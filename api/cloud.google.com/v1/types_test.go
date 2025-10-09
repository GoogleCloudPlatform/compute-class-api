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
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// TestProtobufOrderIsIncreasing automatically checks that for every struct in
// types.go, the protobuf field numbers are in strictly increasing order.
// This test works by parsing the source file and inspecting the AST, so it
// doesn't need to be manually updated when new structs are added.
func TestProtobufOrderIsIncreasing(t *testing.T) {
	fset := token.NewFileSet()
	// The path is relative to the package directory.
	node, err := parser.ParseFile(fset, "types.go", nil, 0)
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
