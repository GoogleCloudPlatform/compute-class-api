/*
* Copyright 2026 Google LLC
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
	"testing"
)

// Embedding the test files is needed because tests are also executed inside google3 (this repo is copied by copybara)
// where it's not executed from the package root directory.
//
//go:embed types_test.go
var typesTestGoSource []byte

//go:embed node_system_config_test.go
var nodeSystemConfigTestGoSource []byte

//go:embed integration_test.go
var integrationTestGoSource []byte

//go:embed structure_test.go
var structureTestGoSource []byte

func TestNoRawMapsInTests(t *testing.T) {
	testSources := map[string][]byte{
		"types_test.go":              typesTestGoSource,
		"node_system_config_test.go": nodeSystemConfigTestGoSource,
		"integration_test.go":        integrationTestGoSource,
		"structure_test.go":          structureTestGoSource,
	}

	for filename, src := range testSources {
		t.Run(filename, func(t *testing.T) {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, filename, src, 0)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", filename, err)
			}

			ast.Inspect(node, func(n ast.Node) bool {
				// Look for struct definitions in tests
				structType, ok := n.(*ast.StructType)
				if !ok {
					return true
				}

				for _, field := range structType.Fields.List {
					// Check if field type is any map type
					if _, ok := field.Type.(*ast.MapType); ok {
						pos := fset.Position(field.Pos())
						t.Errorf("Found forbidden map type in struct field at %s. Use concrete Go structs instead.", pos)
					}
				}

				return true
			})
		})
	}
}
