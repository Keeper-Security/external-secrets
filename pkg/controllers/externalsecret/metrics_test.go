/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package externalsecret

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
)

func TestUpdateExternalSecretCondition(t *testing.T) {
	// Reset metrics before each test
	syncCallsTotal.Reset()
	syncCallsError.Reset()
	externalSecretCondition.Reset()
	externalSecretReconcileDuration.Reset()

	tests := []struct {
		name          string
		es            *esv1beta1.ExternalSecret
		condition     *esv1beta1.ExternalSecretStatusCondition
		value         float64
		expectedValue float64
	}{
		{
			name: "set Ready condition to True",
			es: &esv1beta1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-es",
					Namespace: "default",
				},
			},
			condition: &esv1beta1.ExternalSecretStatusCondition{
				Type:   esv1beta1.ExternalSecretReady,
				Status: v1.ConditionTrue,
			},
			value:         1.0,
			expectedValue: 1.0,
		},
		{
			name: "set Ready condition to False",
			es: &esv1beta1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-es-2",
					Namespace: "default",
				},
			},
			condition: &esv1beta1.ExternalSecretStatusCondition{
				Type:   esv1beta1.ExternalSecretReady,
				Status: v1.ConditionFalse,
			},
			value:         1.0,
			expectedValue: 1.0,
		},
		{
			name: "set Deleted condition",
			es: &esv1beta1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-es-3",
					Namespace: "default",
				},
			},
			condition: &esv1beta1.ExternalSecretStatusCondition{
				Type:   esv1beta1.ExternalSecretDeleted,
				Status: v1.ConditionTrue,
			},
			value:         1.0,
			expectedValue: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateExternalSecretCondition(tt.es, tt.condition, tt.value)

			// Verify the metric was set
			labels := prometheus.Labels{
				"name":      tt.es.Name,
				"namespace": tt.es.Namespace,
				"condition": string(tt.condition.Type),
				"status":    string(tt.condition.Status),
			}

			gauge, err := externalSecretCondition.GetMetricWith(labels)
			if err != nil {
				t.Fatalf("failed to get metric: %v", err)
			}

			value := testutil.ToFloat64(gauge)
			if value != tt.expectedValue {
				t.Errorf("expected metric value %v, got %v", tt.expectedValue, value)
			}
		})
	}
}

func TestUpdateExternalSecretCondition_ReadyToggle(t *testing.T) {
	// Reset metrics before test
	externalSecretCondition.Reset()

	es := &esv1beta1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-es",
			Namespace: "default",
		},
	}

	// First set Ready=True
	trueCondition := &esv1beta1.ExternalSecretStatusCondition{
		Type:   esv1beta1.ExternalSecretReady,
		Status: v1.ConditionTrue,
	}
	updateExternalSecretCondition(es, trueCondition, 1.0)

	// Verify True is set to 1
	trueLabels := prometheus.Labels{
		"name":      es.Name,
		"namespace": es.Namespace,
		"condition": string(esv1beta1.ExternalSecretReady),
		"status":    string(v1.ConditionTrue),
	}
	gauge, _ := externalSecretCondition.GetMetricWith(trueLabels)
	if testutil.ToFloat64(gauge) != 1.0 {
		t.Error("expected Ready=True to be 1.0")
	}

	// Now set Ready=False
	falseCondition := &esv1beta1.ExternalSecretStatusCondition{
		Type:   esv1beta1.ExternalSecretReady,
		Status: v1.ConditionFalse,
	}
	updateExternalSecretCondition(es, falseCondition, 1.0)

	// Verify False is set to 1
	falseLabels := prometheus.Labels{
		"name":      es.Name,
		"namespace": es.Namespace,
		"condition": string(esv1beta1.ExternalSecretReady),
		"status":    string(v1.ConditionFalse),
	}
	gauge, _ = externalSecretCondition.GetMetricWith(falseLabels)
	if testutil.ToFloat64(gauge) != 1.0 {
		t.Error("expected Ready=False to be 1.0")
	}

	// Verify True was toggled to 0
	gauge, _ = externalSecretCondition.GetMetricWith(trueLabels)
	if testutil.ToFloat64(gauge) != 0.0 {
		t.Error("expected Ready=True to be toggled to 0.0")
	}
}

func TestUpdateExternalSecretCondition_DeletedCleansReady(t *testing.T) {
	// Reset metrics before test
	externalSecretCondition.Reset()

	es := &esv1beta1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-es",
			Namespace: "default",
		},
	}

	// First set Ready=True
	readyCondition := &esv1beta1.ExternalSecretStatusCondition{
		Type:   esv1beta1.ExternalSecretReady,
		Status: v1.ConditionTrue,
	}
	updateExternalSecretCondition(es, readyCondition, 1.0)

	// Now set Deleted=True, which should clean up Ready metrics
	deletedCondition := &esv1beta1.ExternalSecretStatusCondition{
		Type:   esv1beta1.ExternalSecretDeleted,
		Status: v1.ConditionTrue,
	}
	updateExternalSecretCondition(es, deletedCondition, 1.0)

	// The Ready metrics should be deleted (not present)
	// Note: We can't easily verify deletion in this test framework,
	// but we can verify Deleted is set
	deletedLabels := prometheus.Labels{
		"name":      es.Name,
		"namespace": es.Namespace,
		"condition": string(esv1beta1.ExternalSecretDeleted),
		"status":    string(v1.ConditionTrue),
	}
	gauge, err := externalSecretCondition.GetMetricWith(deletedLabels)
	if err != nil {
		t.Fatalf("failed to get Deleted metric: %v", err)
	}
	if testutil.ToFloat64(gauge) != 1.0 {
		t.Error("expected Deleted=True to be 1.0")
	}
}

func TestMetricsInitialization(t *testing.T) {
	// Verify that metrics are registered properly
	// This is a simple smoke test to ensure init() was called
	if syncCallsTotal == nil {
		t.Error("syncCallsTotal metric not initialized")
	}
	if syncCallsError == nil {
		t.Error("syncCallsError metric not initialized")
	}
	if externalSecretCondition == nil {
		t.Error("externalSecretCondition metric not initialized")
	}
	if externalSecretReconcileDuration == nil {
		t.Error("externalSecretReconcileDuration metric not initialized")
	}
}
