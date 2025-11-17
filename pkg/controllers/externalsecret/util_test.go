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
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
)

func TestNewExternalSecretCondition(t *testing.T) {
	before := time.Now()
	condition := NewExternalSecretCondition(
		esv1beta1.ExternalSecretReady,
		v1.ConditionTrue,
		"TestReason",
		"Test message",
	)
	after := time.Now()

	if condition.Type != esv1beta1.ExternalSecretReady {
		t.Errorf("expected Type %v, got %v", esv1beta1.ExternalSecretReady, condition.Type)
	}
	if condition.Status != v1.ConditionTrue {
		t.Errorf("expected Status %v, got %v", v1.ConditionTrue, condition.Status)
	}
	if condition.Reason != "TestReason" {
		t.Errorf("expected Reason %q, got %q", "TestReason", condition.Reason)
	}
	if condition.Message != "Test message" {
		t.Errorf("expected Message %q, got %q", "Test message", condition.Message)
	}

	// Verify LastTransitionTime is set to now
	transitionTime := condition.LastTransitionTime.Time
	if transitionTime.Before(before) || transitionTime.After(after) {
		t.Errorf("LastTransitionTime %v not between %v and %v", transitionTime, before, after)
	}
}

func TestGetExternalSecretCondition(t *testing.T) {
	status := esv1beta1.ExternalSecretStatus{
		Conditions: []esv1beta1.ExternalSecretStatusCondition{
			{
				Type:    esv1beta1.ExternalSecretReady,
				Status:  v1.ConditionTrue,
				Reason:  "Ready",
				Message: "Secret is ready",
			},
			{
				Type:    esv1beta1.ExternalSecretDeleted,
				Status:  v1.ConditionFalse,
				Reason:  "NotDeleted",
				Message: "Secret is not deleted",
			},
		},
	}

	tests := []struct {
		name     string
		condType esv1beta1.ExternalSecretConditionType
		want     *esv1beta1.ExternalSecretStatusCondition
	}{
		{
			name:     "find Ready condition",
			condType: esv1beta1.ExternalSecretReady,
			want: &esv1beta1.ExternalSecretStatusCondition{
				Type:    esv1beta1.ExternalSecretReady,
				Status:  v1.ConditionTrue,
				Reason:  "Ready",
				Message: "Secret is ready",
			},
		},
		{
			name:     "find Deleted condition",
			condType: esv1beta1.ExternalSecretDeleted,
			want: &esv1beta1.ExternalSecretStatusCondition{
				Type:    esv1beta1.ExternalSecretDeleted,
				Status:  v1.ConditionFalse,
				Reason:  "NotDeleted",
				Message: "Secret is not deleted",
			},
		},
		{
			name:     "condition not found",
			condType: "NonExistentCondition",
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetExternalSecretCondition(status, tt.condType)
			if tt.want == nil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("expected condition, got nil")
			}
			if got.Type != tt.want.Type {
				t.Errorf("Type: expected %v, got %v", tt.want.Type, got.Type)
			}
			if got.Status != tt.want.Status {
				t.Errorf("Status: expected %v, got %v", tt.want.Status, got.Status)
			}
			if got.Reason != tt.want.Reason {
				t.Errorf("Reason: expected %v, got %v", tt.want.Reason, got.Reason)
			}
			if got.Message != tt.want.Message {
				t.Errorf("Message: expected %v, got %v", tt.want.Message, got.Message)
			}
		})
	}
}

func TestSetExternalSecretCondition(t *testing.T) {
	tests := []struct {
		name              string
		initialConditions []esv1beta1.ExternalSecretStatusCondition
		newCondition      esv1beta1.ExternalSecretStatusCondition
		expectedCount     int
		expectedStatus    v1.ConditionStatus
	}{
		{
			name:              "add new condition to empty list",
			initialConditions: []esv1beta1.ExternalSecretStatusCondition{},
			newCondition: esv1beta1.ExternalSecretStatusCondition{
				Type:    esv1beta1.ExternalSecretReady,
				Status:  v1.ConditionTrue,
				Reason:  "Ready",
				Message: "Secret is ready",
			},
			expectedCount:  1,
			expectedStatus: v1.ConditionTrue,
		},
		{
			name: "update existing condition with same status",
			initialConditions: []esv1beta1.ExternalSecretStatusCondition{
				{
					Type:               esv1beta1.ExternalSecretReady,
					Status:             v1.ConditionTrue,
					Reason:             "OldReason",
					Message:            "Old message",
					LastTransitionTime: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
				},
			},
			newCondition: esv1beta1.ExternalSecretStatusCondition{
				Type:               esv1beta1.ExternalSecretReady,
				Status:             v1.ConditionTrue,
				Reason:             "NewReason",
				Message:            "New message",
				LastTransitionTime: metav1.Now(),
			},
			expectedCount:  1,
			expectedStatus: v1.ConditionTrue,
		},
		{
			name: "update existing condition with different status",
			initialConditions: []esv1beta1.ExternalSecretStatusCondition{
				{
					Type:               esv1beta1.ExternalSecretReady,
					Status:             v1.ConditionTrue,
					Reason:             "Ready",
					Message:            "Secret is ready",
					LastTransitionTime: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
				},
			},
			newCondition: esv1beta1.ExternalSecretStatusCondition{
				Type:               esv1beta1.ExternalSecretReady,
				Status:             v1.ConditionFalse,
				Reason:             "NotReady",
				Message:            "Secret is not ready",
				LastTransitionTime: metav1.Now(),
			},
			expectedCount:  1,
			expectedStatus: v1.ConditionFalse,
		},
		{
			name: "add new condition to existing list",
			initialConditions: []esv1beta1.ExternalSecretStatusCondition{
				{
					Type:    esv1beta1.ExternalSecretReady,
					Status:  v1.ConditionTrue,
					Reason:  "Ready",
					Message: "Secret is ready",
				},
			},
			newCondition: esv1beta1.ExternalSecretStatusCondition{
				Type:    esv1beta1.ExternalSecretDeleted,
				Status:  v1.ConditionFalse,
				Reason:  "NotDeleted",
				Message: "Secret is not deleted",
			},
			expectedCount:  2,
			expectedStatus: v1.ConditionFalse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &esv1beta1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-es",
					Namespace: "default",
				},
				Status: esv1beta1.ExternalSecretStatus{
					Conditions: tt.initialConditions,
				},
			}

			SetExternalSecretCondition(es, tt.newCondition)

			if len(es.Status.Conditions) != tt.expectedCount {
				t.Errorf("expected %d conditions, got %d", tt.expectedCount, len(es.Status.Conditions))
			}

			// Find the new condition
			var found *esv1beta1.ExternalSecretStatusCondition
			for i := range es.Status.Conditions {
				if es.Status.Conditions[i].Type == tt.newCondition.Type {
					found = &es.Status.Conditions[i]
					break
				}
			}

			if found == nil {
				t.Fatal("condition not found in status")
			}

			if found.Status != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, found.Status)
			}
		})
	}
}

func TestSetExternalSecretCondition_SameCondition(t *testing.T) {
	// Test that setting the exact same condition doesn't change anything
	transitionTime := metav1.NewTime(time.Now().Add(-1 * time.Hour))
	es := &esv1beta1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-es",
			Namespace: "default",
		},
		Status: esv1beta1.ExternalSecretStatus{
			Conditions: []esv1beta1.ExternalSecretStatusCondition{
				{
					Type:               esv1beta1.ExternalSecretReady,
					Status:             v1.ConditionTrue,
					Reason:             "Ready",
					Message:            "Secret is ready",
					LastTransitionTime: transitionTime,
				},
			},
		},
	}

	// Set the exact same condition
	sameCondition := esv1beta1.ExternalSecretStatusCondition{
		Type:    esv1beta1.ExternalSecretReady,
		Status:  v1.ConditionTrue,
		Reason:  "Ready",
		Message: "Secret is ready",
	}

	SetExternalSecretCondition(es, sameCondition)

	// Verify conditions list still has 1 item
	if len(es.Status.Conditions) != 1 {
		t.Errorf("expected 1 condition, got %d", len(es.Status.Conditions))
	}
}

func TestFilterOutCondition(t *testing.T) {
	conditions := []esv1beta1.ExternalSecretStatusCondition{
		{
			Type:    esv1beta1.ExternalSecretReady,
			Status:  v1.ConditionTrue,
			Reason:  "Ready",
			Message: "Secret is ready",
		},
		{
			Type:    esv1beta1.ExternalSecretDeleted,
			Status:  v1.ConditionFalse,
			Reason:  "NotDeleted",
			Message: "Secret is not deleted",
		},
	}

	tests := []struct {
		name          string
		condType      esv1beta1.ExternalSecretConditionType
		expectedCount int
	}{
		{
			name:          "filter out Ready",
			condType:      esv1beta1.ExternalSecretReady,
			expectedCount: 1,
		},
		{
			name:          "filter out Deleted",
			condType:      esv1beta1.ExternalSecretDeleted,
			expectedCount: 1,
		},
		{
			name:          "filter out non-existent",
			condType:      "NonExistent",
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterOutCondition(conditions, tt.condType)
			if len(result) != tt.expectedCount {
				t.Errorf("expected %d conditions, got %d", tt.expectedCount, len(result))
			}

			// Verify the filtered condition is not in the result
			for _, cond := range result {
				if cond.Type == tt.condType {
					t.Errorf("condition %v should have been filtered out", tt.condType)
				}
			}
		})
	}
}
