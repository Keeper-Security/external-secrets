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
package v1beta1

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestGenericStoreValidator_ValidateCreate(t *testing.T) {
	validator := &GenericStoreValidator{}
	ctx := context.Background()

	tests := []struct {
		name    string
		obj     runtime.Object
		wantErr bool
	}{
		{
			name: "invalid object type",
			obj: &ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "not-a-store",
					Namespace: "default",
				},
			},
			wantErr: true,
		},
		{
			name:    "nil provider",
			obj: &SecretStore{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-store",
					Namespace: "default",
				},
				Spec: SecretStoreSpec{
					Provider: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "empty provider",
			obj: &SecretStore{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-store",
					Namespace: "default",
				},
				Spec: SecretStoreSpec{
					Provider: &SecretStoreProvider{},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreate(ctx, tt.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenericStoreValidator_ValidateUpdate(t *testing.T) {
	validator := &GenericStoreValidator{}
	ctx := context.Background()

	validStore := &SecretStore{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-store",
			Namespace: "default",
		},
		Spec: SecretStoreSpec{
			Provider: &SecretStoreProvider{},
		},
	}

	tests := []struct {
		name    string
		oldObj  runtime.Object
		newObj  runtime.Object
		wantErr bool
	}{
		{
			name:   "update to empty provider",
			oldObj: validStore,
			newObj: &SecretStore{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-store",
					Namespace: "default",
				},
				Spec: SecretStoreSpec{
					Provider: &SecretStoreProvider{},
				},
			},
			wantErr: true,
		},
		{
			name:   "update to nil provider",
			oldObj: validStore,
			newObj: &SecretStore{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-store",
					Namespace: "default",
				},
				Spec: SecretStoreSpec{
					Provider: nil,
				},
			},
			wantErr: true,
		},
		{
			name:   "invalid new object type",
			oldObj: validStore,
			newObj: &ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "not-a-store",
					Namespace: "default",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateUpdate(ctx, tt.oldObj, tt.newObj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenericStoreValidator_ValidateDelete(t *testing.T) {
	validator := &GenericStoreValidator{}
	ctx := context.Background()

	tests := []struct {
		name    string
		obj     runtime.Object
		wantErr bool
	}{
		{
			name: "valid delete",
			obj: &SecretStore{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-store",
					Namespace: "default",
				},
			},
			wantErr: false,
		},
		{
			name: "delete cluster store",
			obj: &ClusterSecretStore{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster-store",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateDelete(ctx, tt.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
