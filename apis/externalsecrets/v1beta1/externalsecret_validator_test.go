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
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
)

func TestValidateExternalSecret(t *testing.T) {
	tests := []struct {
		name    string
		obj     runtime.Object
		wantErr bool
	}{
		{
			name:    "nil",
			obj:     nil,
			wantErr: true,
		},
		{
			name: "deletion policy delete with merge creation policy",
			obj: &ExternalSecret{
				Spec: ExternalSecretSpec{
					Target: ExternalSecretTarget{
						DeletionPolicy: DeletionPolicyDelete,
						CreationPolicy: CreatePolicyMerge,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "deletion policy delete with none creation policy",
			obj: &ExternalSecret{
				Spec: ExternalSecretSpec{
					Target: ExternalSecretTarget{
						DeletionPolicy: DeletionPolicyDelete,
						CreationPolicy: CreatePolicyNone,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "deletion policy merge with none creation policy",
			obj: &ExternalSecret{
				Spec: ExternalSecretSpec{
					Target: ExternalSecretTarget{
						DeletionPolicy: DeletionPolicyMerge,
						CreationPolicy: CreatePolicyNone,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid deletion policy delete with owner creation policy",
			obj: &ExternalSecret{
				Spec: ExternalSecretSpec{
					Target: ExternalSecretTarget{
						DeletionPolicy: DeletionPolicyDelete,
						CreationPolicy: CreatePolicyOwner,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid deletion policy merge with owner creation policy",
			obj: &ExternalSecret{
				Spec: ExternalSecretSpec{
					Target: ExternalSecretTarget{
						DeletionPolicy: DeletionPolicyMerge,
						CreationPolicy: CreatePolicyOwner,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid deletion policy merge with merge creation policy",
			obj: &ExternalSecret{
				Spec: ExternalSecretSpec{
					Target: ExternalSecretTarget{
						DeletionPolicy: DeletionPolicyMerge,
						CreationPolicy: CreatePolicyMerge,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "generator with find",
			obj: &ExternalSecret{
				Spec: ExternalSecretSpec{
					DataFrom: []ExternalSecretDataFromRemoteRef{
						{
							Find: &ExternalSecretFind{},
							SourceRef: &SourceRef{
								GeneratorRef: &GeneratorRef{},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "generator with extract",
			obj: &ExternalSecret{
				Spec: ExternalSecretSpec{
					DataFrom: []ExternalSecretDataFromRemoteRef{
						{
							Extract: &ExternalSecretDataRemoteRef{},
							SourceRef: &SourceRef{
								GeneratorRef: &GeneratorRef{},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "generator without find or extract is valid",
			obj: &ExternalSecret{
				Spec: ExternalSecretSpec{
					DataFrom: []ExternalSecretDataFromRemoteRef{
						{
							SourceRef: &SourceRef{
								GeneratorRef: &GeneratorRef{
									APIVersion: "generators.external-secrets.io/v1alpha1",
									Kind:       "Password",
									Name:       "test-generator",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple dataFrom with mixed validity",
			obj: &ExternalSecret{
				Spec: ExternalSecretSpec{
					DataFrom: []ExternalSecretDataFromRemoteRef{
						{
							Extract: &ExternalSecretDataRemoteRef{
								Key: "valid-key",
							},
						},
						{
							Find: &ExternalSecretFind{},
							SourceRef: &SourceRef{
								GeneratorRef: &GeneratorRef{},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid with empty dataFrom",
			obj: &ExternalSecret{
				Spec: ExternalSecretSpec{
					DataFrom: []ExternalSecretDataFromRemoteRef{},
				},
			},
			wantErr: false,
		},
		{
			name: "valid with nil dataFrom",
			obj: &ExternalSecret{
				Spec: ExternalSecretSpec{
					DataFrom: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "wrong type object",
			obj: &SecretStore{
				Spec: SecretStoreSpec{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateExternalSecret(tt.obj); (err != nil) != tt.wantErr {
				t.Errorf("validateExternalSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
