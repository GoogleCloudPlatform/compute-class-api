/*
* Copyright 2026 Google LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     https://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/
package v1

import (
	"context"
	"path/filepath"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	k8sptr "k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func TestComputeClassValidationSmoke(t *testing.T) {
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "cloud.google.com_computeclasses.yaml")},
	}

	cfg, err := testEnv.Start()
	if err != nil {
		t.Fatalf("Failed to start envtest: %v", err)
	}
	defer testEnv.Stop()

	err = AddToScheme(scheme.Scheme)
	if err != nil {
		t.Fatalf("Failed to add to scheme: %v", err)
	}

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		t.Fatalf("Failed to create k8s client: %v", err)
	}

	tests := []struct {
		name      string
		spec      ComputeClassSpec
		wantValid bool
	}{
		// 1. Comprehensive Positive Smoke Test (Happy Path)
		// Representative of standard valid configurations with multiple features enabled.
		// Proves that basic schema wiring works, valid objects can be persisted,
		// and features don't falsely conflict when used correctly together.
		// Covers: valid ANP, consistent PriorityScores, GpuDirect with ANP,
		// and independent usage of Spot/FlexStart.
		{
			name: "valid-comprehensive-smoke",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				Priorities: []Priority{
					{
						PriorityScore:            k8sptr.To(1),
						AcceleratorNetworkProfile: k8sptr.To("auto-profile"),
						MachineFamily:            k8sptr.To("c3"),
						GpuDirect:                "rdma", // Valid with ANP
					},
					{
						PriorityScore: k8sptr.To(2),
						Spot:          k8sptr.To(true), // Valid without FlexStart
					},
					{
						PriorityScore: k8sptr.To(3),
						FlexStart:     &FlexStart{Enabled: true}, // Valid without Spot
					},
				},
			},
			wantValid: true,
		},
		// 2. Representative of Regex/Pattern Validations
		// Proves that custom CEL regex rules (e.g., self.matches(...)) are correctly
		// compiled and executed by the real Kubernetes API server runtime.
		{
			name: "invalid-anp-starts-with-number-smoke",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				Priorities: []Priority{
					{
						AcceleratorNetworkProfile: k8sptr.To("1profile"),
						MachineFamily:            k8sptr.To("c3"),
					},
				},
			},
			wantValid: false,
		},
		// 3. Representative of Cross-Field Incompatibility Rules
		// Proves that CEL rules enforcing mutual exclusivity between two different fields in the same struct
		// (e.g., FlexStart vs Spot) are correctly interpreted by the control plane.
		{
			name: "invalid-flex-start-with-spot",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				Priorities: []Priority{
					{
						Spot:      k8sptr.To(true),
						FlexStart: &FlexStart{Enabled: true},
					},
				},
			},
			wantValid: false,
		},
		// 4. Representative of Conditional Requirements
		// Proves that rules enforcing one field conditionally based on another field's value
		// (e.g., gpuDirect == 'rdma' requires acceleratorNetworkProfile) are active.
		{
			name: "invalid-gpu-direct-without-anp",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				Priorities: []Priority{
					{
						GpuDirect: "rdma",
						// AcceleratorNetworkProfile is missing
					},
				},
			},
			wantValid: false,
		},
		// 5. Representative of List-Level Aggregations & Consistency
		// Proves that complex CEL logic iterating over arrays (e.g., priorities.all(...))
		// is correctly executed and doesn't hit unexpected runtime limits in Kubernetes.
		{
			name: "invalid-priority-score-mismatch",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				Priorities: []Priority{
					{PriorityScore: k8sptr.To(1)},
					{PriorityScore: nil}, // Mismatch
				},
			},
			wantValid: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cc := &ComputeClass{
				ObjectMeta: metav1.ObjectMeta{Name: tc.name},
				Spec:      tc.spec,
			}
			err := k8sClient.Create(context.Background(), cc)
			if tc.wantValid {
				if err != nil {
					t.Errorf("Expected valid but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			}
			// Clean up
			if err == nil {
				k8sClient.Delete(context.Background(), cc)
			}
		})
	}
}
