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

// Package v1 contains integration tests for ComputeClass API validation.
//
// IMPORTANT — READ BEFORE MODIFYING THIS FILE:
// integration_test.go implements Layer 2 (Integration/envtest) validation for ComputeClass API.
// It is strictly a minimal smoke test for post-compilation Kubernetes API server checks:
//   - Verifying that the generated CRD YAML installs cleanly into a real K8s API server.
//   - Validating K8s CEL cost/budget estimation limits across the entire CRD schema (verified
//     automatically when TestComputeClassValidationSmoke starts envtest and applies the CRD).
//   - Smoke-testing one representative case per structural class of validation mechanism
//     (e.g., regex, cross-field, list aggregation, root-level metadata vs spec).
//
// Strict Rules for Contributors (Human & AI):
//   - DO NOT add unit test permutations, boundary values, or feature-specific edge cases here.
//     All permutation and edge-case testing belongs in Layer 1 (types_test.go).
//   - DO NOT add new test cases for a new field/rule if its validation structure (e.g., regex,
//     cross-field check, list comprehension) is already represented in TestComputeClassValidationSmoke.
//     Existing smoke startup already validates that your new rule compiles and fits within K8s CEL budget limits.
//   - DO NOT spin up additional envtest.Environment instances, add separate Test* functions,
//     patch CRDs at runtime, or add test-only dependencies to go.mod.
//   - If Layer 1 (types_test.go) lacks support for a CEL feature needed by your rule (e.g., oldSelf
//     transition rules, cel.OptionalTypes(), or markers at a new AST location), extend the Layer 1
//     test harness (createCELProgram / getTypeValidationRules in types_test.go) instead of falling
//     back to envtest.
//   - Only add a single minimal case to TestComputeClassValidationSmoke if introducing a genuinely
//     new structural validation mechanism not yet covered below.
//   - SUBMIT REQUIREMENT: Any commit modifying this file must include
//     `E2E_ACK=<justification confirming unique structure and no permutations/edge cases>`
//     in the commit message, or Gerrit will block submission.
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
				NetworkConfig: &NetworkConfig{
					SubnetPriorities: []SubnetPriority{
						{
							Name:     "subnet-a",
							PodRange: "pod-range-a",
						},
					},
				},
				NodePoolConfig: &NodePoolConfig{
					ResourceLabels: map[string]ResourceLabelValue{
						"billing":    "marketing",
						"env-prod_1": "valid-val_1",
					},
				},
				Priorities: []Priority{
					{
						PriorityScore:             k8sptr.To(1),
						AcceleratorNetworkProfile: k8sptr.To("auto-profile"),
						MachineFamily:             k8sptr.To("c3"),
						GpuDirect:                 "rdma", // Valid with ANP
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
						MachineFamily:             k8sptr.To("c3"),
					},
				},
			},
			wantValid: false,
		},
		{
			name: "invalid-subnet-priority-bad-pod-range-smoke",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				NetworkConfig: &NetworkConfig{
					SubnetPriorities: []SubnetPriority{
						{
							Name:     "subnet-a",
							PodRange: "123-bad-pod-range",
						},
					},
				},
				Priorities: []Priority{
					{MachineFamily: k8sptr.To("c3")},
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
		// 6. Representative of ResourceLabels Validations
		// Proves that resourceLabels key regex and value pattern validation work in Kubernetes runtime.
		{
			name: "invalid-resource-labels-key-format",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				NodePoolConfig: &NodePoolConfig{
					ResourceLabels: map[string]ResourceLabelValue{
						"1billing": "marketing",
					},
				},
			},
			wantValid: false,
		},
		{
			name: "invalid-resource-labels-value-format",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				NodePoolConfig: &NodePoolConfig{
					ResourceLabels: map[string]ResourceLabelValue{
						"billing": "Marketing",
					},
				},
			},
			wantValid: false,
		},
		// 7. Multi-NIC (additionalNetworks) & SubnetPriority Validations
		{
			name: "valid-additional-networks-with-pod-range",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				NetworkConfig: &NetworkConfig{
					AdditionalNetworks: []AdditionalNetwork{
						{
							Name: "multinet-vpc-1",
							SubnetPriorities: []SubnetPriority{
								{
									Name:     "multinet-subnet-1",
									PodRange: "multinet-pods-1",
								},
							},
						},
					},
				},
				Priorities: []Priority{
					{MachineFamily: k8sptr.To("c3")},
				},
			},
			wantValid: true,
		},
		{
			name: "valid-additional-networks-without-pod-range",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				NetworkConfig: &NetworkConfig{
					AdditionalNetworks: []AdditionalNetwork{
						{
							Name: "rdma-vpc-2",
							SubnetPriorities: []SubnetPriority{
								{
									Name: "rdma-subnet-2",
								},
							},
						},
					},
				},
				Priorities: []Priority{
					{MachineFamily: k8sptr.To("c3")},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid-primary-subnet-priorities-missing-pod-range",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				NetworkConfig: &NetworkConfig{
					SubnetPriorities: []SubnetPriority{
						{
							Name: "subnet-a",
						},
					},
				},
				Priorities: []Priority{
					{MachineFamily: k8sptr.To("c3")},
				},
			},
			wantValid: false,
		},
		{
			name: "invalid-additional-networks-too-many-subnet-priorities",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				NetworkConfig: &NetworkConfig{
					AdditionalNetworks: []AdditionalNetwork{
						{
							Name: "multinet-vpc-1",
							SubnetPriorities: []SubnetPriority{
								{Name: "subnet-a"},
								{Name: "subnet-b"},
							},
						},
					},
				},
				Priorities: []Priority{
					{MachineFamily: k8sptr.To("c3")},
				},
			},
			wantValid: false,
		},
		{
			name: "invalid-additional-networks-with-anp",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				NetworkConfig: &NetworkConfig{
					AdditionalNetworks: []AdditionalNetwork{
						{
							Name: "multinet-vpc-1",
							SubnetPriorities: []SubnetPriority{
								{Name: "multinet-subnet-1"},
							},
						},
					},
				},
				Priorities: []Priority{
					{
						AcceleratorNetworkProfile: k8sptr.To("auto-profile"),
						MachineFamily:             k8sptr.To("c3"),
					},
				},
			},
			wantValid: false,
		},
		// 8. Representative of Root-Level Rules on metadata.name
		// Proves that a rule on the whole object can read metadata.name and is
		// run by the API server on create. The case name is also the object
		// name, so it has 64 characters, one more than the limit.
		{
			name: "invalid-name-longer-than-63-characters-xxxxxxxxxxxxxxxxxxxxxxxxx",
			spec: ComputeClassSpec{
				WhenUnsatisfiable: "DoNotScaleUp",
				Priorities: []Priority{
					{MachineFamily: k8sptr.To("c3")},
				},
			},
			wantValid: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cc := &ComputeClass{
				ObjectMeta: metav1.ObjectMeta{Name: tc.name},
				Spec:       tc.spec,
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
