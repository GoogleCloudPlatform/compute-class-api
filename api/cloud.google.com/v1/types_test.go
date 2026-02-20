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

// Embedding the file is needed because the test is also executed inside google3 (this repo is copied by copybara)
// where it's not executed from the package root directory.
//
//go:embed types.go
var typesGoSource []byte

// TestProtobufOrderIsIncreasing automatically checks that for every struct in
// types.go, the protobuf field numbers are in strictly increasing order.
// This test works by parsing the source file and inspecting the AST, so it
// doesn't need to be manually updated when new structs are added.
func TestProtobufOrderIsIncreasing(t *testing.T) {
	fset := token.NewFileSet()
	// The path is relative to the package directory.
	node, err := parser.ParseFile(fset, "types.go", typesGoSource, 0)
	if err != nil {
		t.Fatalf("Failed to parse types.go: %v", err)
	}

	// ast.Inspect traverses the AST of the parsed file.
	ast.Inspect(node, func(n ast.Node) bool {
		// We are looking for type declarations.
		typeSpec, isTypeSpec := n.(*ast.TypeSpec)
		if !isTypeSpec {
			return true // Continue traversal.
		}

		// We are only interested in struct types.
		structType, isStruct := typeSpec.Type.(*ast.StructType)
		if !isStruct {
			return true // Continue traversal.
		}

		// Run a subtest for each struct found.
		t.Run(typeSpec.Name.Name, func(t *testing.T) {
			lastProtobufNum := 0
			// Iterate over all fields in the struct.
			for _, field := range structType.Fields.List {
				if field.Tag == nil {
					continue
				}

				// field.Tag.Value is a raw string like `json:"..." protobuf:"..."`
				// We need to unquote it to handle escape characters.
				tagString, err := strconv.Unquote(field.Tag.Value)
				if err != nil {
					t.Errorf("could not unquote tag for a field in struct %s: %v", typeSpec.Name.Name, err)
					continue
				}

				// Use reflect.StructTag to easily parse the tags.
				protoTag := reflect.StructTag(tagString).Get("protobuf")
				if protoTag == "" {
					continue
				}

				// The protobuf tag is comma-separated, e.g., "bytes,1,opt,name=metadata"
				parts := strings.Split(protoTag, ",")
				if len(parts) < 2 {
					t.Logf("skipping field with unparseable protobuf tag: %s", protoTag)
					continue
				}

				// The second part should be the field number.
				num, err := strconv.Atoi(parts[1])
				if err != nil {
					t.Errorf("could not parse protobuf number from tag: %s", protoTag)
					continue
				}

				// Check if the number is strictly greater than the previous one.
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

		// We've processed this struct, no need to inspect its children nodes.
		return false
	})
}

func TestGpuTopologyValidationRule(t *testing.T) {
	rule := getTypeValidationRule(t, "ComputeClassSpec", "gpu.topology")
	program := createCELProgram(t, rule)

	tests := []struct {
		name      string
		input     map[string]interface{}
		wantValid bool
	}{
		{
			name: "valid:_no_gpu_topology",
			input: map[string]interface{}{
				"priorities": []map[string]interface{}{
					{
						"machineFamily": "c3",
					},
				},
			},
			wantValid: true,
		},
		{
			name: "valid:_gpu_topology_with_a4x_with_placement",
			input: map[string]interface{}{
				"priorities": []map[string]interface{}{
					{
						"machineFamily": "a4x",
						"gpu": map[string]interface{}{
							"topology": "1x72",
						},
						"placement": map[string]interface{}{
							"policyName": "workloadPolicy",
						},
					},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid:_gpu_topology_with_a4x_without_placement",
			input: map[string]interface{}{
				"priorities": []map[string]interface{}{
					{
						"machineFamily": "a4x",
						"gpu": map[string]interface{}{
							"topology": "1x72",
						},
					},
				},
			},
			wantValid: false,
		},
		{
			name: "valid:_gpu_topology_with_nvidia-gb200_with_placement",
			input: map[string]interface{}{
				"priorities": []map[string]interface{}{
					{
						"gpu": map[string]interface{}{
							"type":     "nvidia-gb200",
							"topology": "1x72",
						},
						"placement": map[string]interface{}{
							"policyName": "workloadPolicy",
						},
					},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid:_gpu_topology_with_nvidia-gb200_without_policy",
			input: map[string]interface{}{
				"priorities": []map[string]interface{}{
					{
						"gpu": map[string]interface{}{
							"type":     "nvidia-gb200",
							"topology": "1x72",
						},
					},
				},
			},
			wantValid: false,
		},
		{
			name: "valid:_nvidia-gb200_with_policy_without_topology",
			input: map[string]interface{}{
				"priorities": []map[string]interface{}{
					{
						"machineFamily": "a4x",
						"gpu": map[string]interface{}{
							"type": "nvidia-gb200",
						},
						"placement": map[string]interface{}{
							"policyName": "workloadPolicy",
						},
					},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid:_gpu_topology_with_c3",
			input: map[string]interface{}{
				"priorities": []map[string]interface{}{
					{
						"machineFamily": "c3",
						"gpu": map[string]interface{}{
							"topology": "1x72",
						},
					},
				},
			},
			wantValid: false,
		},
		{
			name: "invalid:_gpu_topology_with_other_gpu_type",
			input: map[string]interface{}{
				"priorities": []map[string]interface{}{
					{
						"gpu": map[string]interface{}{
							"type":     "nvidia-h100-80gb",
							"topology": "1x72",
						},
					},
				},
			},
			wantValid: false,
		},
		{
			name: "valid:_gpu_without_topology",
			input: map[string]interface{}{
				"priorities": []map[string]interface{}{
					{
						"machineFamily": "c3",
						"gpu": map[string]interface{}{
							"type": "nvidia-h100-80gb",
						},
					},
				},
			},
			wantValid: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, _, err := program.Eval(map[string]interface{}{
				"self": tc.input,
			})
			if err != nil {
				t.Fatalf("CEL evaluation failed: %v", err)
			}

			if out.Value() != tc.wantValid {
				t.Errorf("Validation result = %v, want %v", out.Value(), tc.wantValid)
			}
		})
	}
}

func TestLocationZoneTypes_WithLocationZones(t *testing.T) {
	rule := getTypeValidationRule(t, "Location", "!(has(self.zones) && has(self.zoneTypes))")
	program := createCELProgram(t, rule)

	tests := []struct {
		name      string
		input     map[string]interface{}
		wantValid bool
	}{
		{
			name: "only_zones",
			input: map[string]interface{}{
				"zones": []string{"us-central1-a"},
			},
			wantValid: true,
		},
		{
			name: "only_zoneTypes",
			input: map[string]interface{}{
				"zoneTypes": []string{"STANDARD"},
			},
			wantValid: true,
		},
		{
			name:      "neither",
			input:     map[string]interface{}{},
			wantValid: true,
		},
		{
			name: "both_set",
			input: map[string]interface{}{
				"zones":     []string{"us-central1-a"},
				"zoneTypes": []string{"STANDARD"},
			},
			wantValid: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, _, err := program.Eval(map[string]interface{}{
				"self": tc.input,
			})
			if err != nil {
				t.Fatalf("CEL evaluation failed: %v", err)
			}

			if out.Value() != tc.wantValid {
				t.Errorf("Validation result = %v, want %v", out.Value(), tc.wantValid)
			}
		})
	}
}

func getTypeValidationRule(t *testing.T, structName, ruleSubString string) string {
	t.Helper()
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "types.go", typesGoSource, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse types.go: %v", err)
	}

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
			if typeSpec.Doc == nil && gd.Doc == nil {
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
						return *rule
					}
				}
			}
		}
	}
	t.Fatalf("Could not find validation rule with %q at struct %s in types.go", ruleSubString, structName)
	return ""
}

func getFieldValidationRule(t *testing.T, structName, fieldName, ruleSubString string) string {
	t.Helper()
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "types.go", typesGoSource, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse types.go: %v", err)
	}

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

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue // skip embedded fields
				}
				if field.Names[0].Name != fieldName {
					continue
				}
				if field.Doc == nil {
					continue
				}
				for _, comment := range field.Doc.List {
					rule := extractRuleFromComment(t, comment.Text, ruleSubString)
					if rule != nil {
						return *rule
					}
				}

			}
		}
	}
	t.Fatalf("Could not find field validation rule with %q on field %s.%s", ruleSubString, structName, fieldName)
	return ""
}

func extractRuleFromComment(t *testing.T, commentText string, ruleSubString string) *string {
	if !strings.Contains(commentText, "+kubebuilder:validation:XValidation:rule") || !strings.Contains(commentText, ruleSubString) {
		return nil
	}
	idx := strings.Index(commentText, "rule=")
	if idx == -1 {
		return nil
	}
	rest := commentText[idx+len("rule="):]
	quotedRule, err := strconv.QuotedPrefix(rest)
	if err != nil {
		t.Logf("Failed to parse quoted rule from comment: %v", err)
		return nil
	}
	var rule string
	rule, err = strconv.Unquote(quotedRule)
	if err != nil {
		t.Logf("Failed to unquote rule: %v", err)
		return nil
	}
	return &rule
}

func createCELProgram(t *testing.T, rule string) cel.Program {
	t.Helper()
	env, err := cel.NewEnv(
		cel.Variable("self", cel.MapType(cel.StringType, cel.DynType)),
	)
	if err != nil {
		t.Fatalf("Failed to create CEL environment: %v", err)
	}

	ast, issues := env.Compile(rule)
	if issues != nil && issues.Err() != nil {
		t.Fatalf("Failed to compile CEL rule: %v", issues.Err())
	}

	program, err := env.Program(ast)
	if err != nil {
		t.Fatalf("Failed to create CEL program: %v", err)
	}

	return program
}
