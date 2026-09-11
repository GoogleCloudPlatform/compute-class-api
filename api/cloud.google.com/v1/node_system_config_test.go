/*
* Copyright 2025 Google LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     https://www.apache.org/licenses/LICENSE-2.0
*
*     Unless required by applicable law or agreed to in writing, software
*     distributed under the License is distributed on an "AS IS" BASIS,
*     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*     See the License for the specific language governing permissions and
*     limitations under the License.
 */

package v1

import (
	_ "embed"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
)

//go:embed node_system_config.go
var nodeSystemConfigGoSource []byte

// TestNodeSystemConfigProtobufOrderIsIncreasing automatically checks that for every struct in
// node_system_config.go, the protobuf field numbers are in strictly increasing order.
func TestNodeSystemConfigProtobufOrderIsIncreasing(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "node_system_config.go", nodeSystemConfigGoSource, 0)
	if err != nil {
		t.Fatalf("Failed to parse node_system_config.go: %v", err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, isTypeSpec := n.(*ast.TypeSpec)
		if !isTypeSpec {
			return true
		}
		structType, isStruct := typeSpec.Type.(*ast.StructType)
		if !isStruct {
			return true
		}
		t.Run(typeSpec.Name.Name, func(t *testing.T) {
			lastProtobufNum := 0
			for _, field := range structType.Fields.List {
				if field.Tag == nil {
					continue
				}
				tagString, err := strconv.Unquote(field.Tag.Value)
				if err != nil {
					t.Errorf("could not unquote tag for a field in struct %s: %v", typeSpec.Name.Name, err)
					continue
				}
				protoTag := reflect.StructTag(tagString).Get("protobuf")
				if protoTag == "" {
					continue
				}
				parts := strings.Split(protoTag, ",")
				if len(parts) < 2 {
					t.Logf("skipping field with unparseable protobuf tag: %s", protoTag)
					continue
				}
				num, err := strconv.Atoi(parts[1])
				if err != nil {
					t.Errorf("could not parse protobuf number from tag: %s", protoTag)
					continue
				}
				if num <= lastProtobufNum {
					fieldName := "unknown"
					if len(field.Names) > 0 {
						fieldName = field.Names[0].Name
					}
					t.Errorf("field '%s' has protobuf number %d, which is not greater than the previous number %d", fieldName, num, lastProtobufNum)
				}
				lastProtobufNum = num
			}
		})
		return true
	})
}

func getNodeSystemConfigValidationRules(t *testing.T, structName, ruleSubString string) []string {
	t.Helper()
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "node_system_config.go", nodeSystemConfigGoSource, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse node_system_config.go: %v", err)
	}

	var rules []string
	for _, decl := range node.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != structName {
				continue
			}
			comments := []*ast.CommentGroup{typeSpec.Doc, gd.Doc}
			for _, cg := range comments {
				if cg == nil {
					continue
				}
				for _, comment := range cg.List {
					rule := extractRuleFromComment(t, comment.Text, ruleSubString)
					if rule != nil {
						rules = append(rules, *rule)
					}
				}
			}
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				for _, field := range structType.Fields.List {
					if field.Doc != nil {
						for _, comment := range field.Doc.List {
							rule := extractRuleFromComment(t, comment.Text, ruleSubString)
							if rule != nil {
								rules = append(rules, *rule)
							}
						}
					}
				}
			}
		}
	}
	if len(rules) == 0 {
		t.Fatalf("Could not find validation rules with %q at struct %s in node_system_config.go", ruleSubString, structName)
	}
	return rules
}

func TestSysctlNetIpv4TcpCongestionControlValidationRule(t *testing.T) {
	rules := getTypeValidationRules(t, "ComputeClassSpec", "net__dot__ipv4__dot__tcp_congestion_control")
	var programs []cel.Program
	for _, rule := range rules {
		programs = append(programs, createCELProgram(t, rule))
	}
	tests := []struct {
		name      string
		input     ComputeClassSpec
		wantValid bool
	}{
		{
			name: "autopilot disabled without sysctl",
			input: ComputeClassSpec{
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
					},
				},
			},
			wantValid: true,
		},
		{
			name: "autopilot disabled with sysctl in priorityDefaults",
			input: ComputeClassSpec{
				PriorityDefaults: &PriorityDefaults{
					NodeSystemConfig: &NodeSystemConfig{
						LinuxNodeConfig: &LinuxNodeConfig{
							Sysctls: &SysctlsConfig{
								Net_ipv4_tcp_congestion_control: ptr("bbr"),
							},
						},
					},
				},
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
					},
				},
			},
			wantValid: true,
		},
		{
			name: "autopilot disabled with sysctl in priorities",
			input: ComputeClassSpec{
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
						NodeSystemConfig: &NodeSystemConfig{
							LinuxNodeConfig: &LinuxNodeConfig{
								Sysctls: &SysctlsConfig{
									Net_ipv4_tcp_congestion_control: ptr("bbr"),
								},
							},
						},
					},
				},
			},
			wantValid: true,
		},
		{
			name: "autopilot enabled without sysctl",
			input: ComputeClassSpec{
				Autopilot: &Autopilot{
					Enabled: true,
				},
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
					},
				},
			},
			wantValid: true,
		},
		{
			name: "autopilot enabled with sysctl in priorityDefaults",
			input: ComputeClassSpec{
				Autopilot: &Autopilot{
					Enabled: true,
				},
				PriorityDefaults: &PriorityDefaults{
					NodeSystemConfig: &NodeSystemConfig{
						LinuxNodeConfig: &LinuxNodeConfig{
							Sysctls: &SysctlsConfig{
								Net_ipv4_tcp_congestion_control: ptr("bbr"),
							},
						},
					},
				},
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
					},
				},
			},
			wantValid: false,
		},
		{
			name: "autopilot enabled with sysctl in priorities",
			input: ComputeClassSpec{
				Autopilot: &Autopilot{
					Enabled: true,
				},
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
						NodeSystemConfig: &NodeSystemConfig{
							LinuxNodeConfig: &LinuxNodeConfig{
								Sysctls: &SysctlsConfig{
									Net_ipv4_tcp_congestion_control: ptr("bbr"),
								},
							},
						},
					},
				},
			},
			wantValid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			isValid := true
			for _, program := range programs {
				out, _, err := program.Eval(map[string]interface{}{
					"self": mustConvertToMap(t, tc.input),
				})

				if err != nil {
					t.Fatalf("CEL evaluation failed: %v", err)
				}
				if out.Value() == false {
					isValid = false
					break
				}
			}
			if isValid != tc.wantValid {
				t.Errorf("Validation result = %v, want %v", isValid, tc.wantValid)
			}
		})
	}
}

func TestKubeletInsecureKubeletReadonlyPortEnabledValidationRule(t *testing.T) {
	rules := getTypeValidationRules(t, "ComputeClassSpec", "insecureKubeletReadonlyPortEnabled")
	var programs []cel.Program
	for _, rule := range rules {
		programs = append(programs, createCELProgram(t, rule))
	}
	tests := []struct {
		name      string
		input     ComputeClassSpec
		wantValid bool
	}{
		{
			name: "autopilot disabled without insecureKubeletReadonlyPortEnabled",
			input: ComputeClassSpec{
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
					},
				},
			},
			wantValid: true,
		},
		{
			name: "autopilot disabled with insecureKubeletReadonlyPortEnabled true in priorityDefaults",
			input: ComputeClassSpec{
				PriorityDefaults: &PriorityDefaults{
					NodeSystemConfig: &NodeSystemConfig{
						KubeletConfig: &KubeletConfig{
							InsecureKubeletReadonlyPortEnabled: ptr(true),
						},
					},
				},
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
					},
				},
			},
			wantValid: true,
		},
		{
			name: "autopilot disabled with insecureKubeletReadonlyPortEnabled true in priorities",
			input: ComputeClassSpec{
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
						NodeSystemConfig: &NodeSystemConfig{
							KubeletConfig: &KubeletConfig{
								InsecureKubeletReadonlyPortEnabled: ptr(true),
							},
						},
					},
				},
			},
			wantValid: true,
		},
		{
			name: "autopilot enabled without insecureKubeletReadonlyPortEnabled",
			input: ComputeClassSpec{
				Autopilot: &Autopilot{
					Enabled: true,
				},
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
					},
				},
			},
			wantValid: true,
		},
		{
			name: "autopilot enabled with insecureKubeletReadonlyPortEnabled false in priorityDefaults",
			input: ComputeClassSpec{
				Autopilot: &Autopilot{
					Enabled: true,
				},
				PriorityDefaults: &PriorityDefaults{
					NodeSystemConfig: &NodeSystemConfig{
						KubeletConfig: &KubeletConfig{
							InsecureKubeletReadonlyPortEnabled: ptr(false),
						},
					},
				},
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
					},
				},
			},
			wantValid: true,
		},
		{
			name: "autopilot enabled with insecureKubeletReadonlyPortEnabled false in priorities",
			input: ComputeClassSpec{
				Autopilot: &Autopilot{
					Enabled: true,
				},
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
						NodeSystemConfig: &NodeSystemConfig{
							KubeletConfig: &KubeletConfig{
								InsecureKubeletReadonlyPortEnabled: ptr(false),
							},
						},
					},
				},
			},
			wantValid: true,
		},
		{
			name: "autopilot enabled with insecureKubeletReadonlyPortEnabled true in priorityDefaults",
			input: ComputeClassSpec{
				Autopilot: &Autopilot{
					Enabled: true,
				},
				PriorityDefaults: &PriorityDefaults{
					NodeSystemConfig: &NodeSystemConfig{
						KubeletConfig: &KubeletConfig{
							InsecureKubeletReadonlyPortEnabled: ptr(true),
						},
					},
				},
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
					},
				},
			},
			wantValid: false,
		},
		{
			name: "autopilot enabled with insecureKubeletReadonlyPortEnabled true in priorities",
			input: ComputeClassSpec{
				Autopilot: &Autopilot{
					Enabled: true,
				},
				Priorities: []Priority{
					{
						MachineFamily: ptr("c3"),
						NodeSystemConfig: &NodeSystemConfig{
							KubeletConfig: &KubeletConfig{
								InsecureKubeletReadonlyPortEnabled: ptr(true),
							},
						},
					},
				},
			},
			wantValid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			isValid := true
			for _, program := range programs {
				out, _, err := program.Eval(map[string]interface{}{
					"self": mustConvertToMap(t, tc.input),
				})

				if err != nil {
					t.Fatalf("CEL evaluation failed: %v", err)
				}
				if out.Value() == false {
					isValid = false
					break
				}
			}
			if isValid != tc.wantValid {
				t.Errorf("Validation result = %v, want %v", isValid, tc.wantValid)
			}
		})
	}
}
