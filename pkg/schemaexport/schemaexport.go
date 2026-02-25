//
// Copyright 2024 The GUAC Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schemaexport

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

type Property struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Comment string `json:"comment,omitempty"`
}

type Node struct {
	Table      string     `json:"table"`
	Properties []Property `json:"properties"`
	Edges      []Edge     `json:"edges,omitempty"`
}

type Edge struct {
	ToTable    string     `json:"to_table"`
	Column     string     `json:"column,omitempty"` // For GraphQL, this might be the field name
	Symbol     string     `json:"symbol"`
	Properties []Property `json:"properties,omitempty"`
}

type SchemaDump struct {
	Nodes []Node `json:"nodes"`
}

func LoadGraphQLSchema(schemaDir string) (*ast.Schema, error) {
	files, err := os.ReadDir(schemaDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema directory: %w", err)
	}

	var sources []*ast.Source
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".graphql") {
			content, err := os.ReadFile(filepath.Join(schemaDir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to read schema file %s: %w", file.Name(), err)
			}
			sources = append(sources, &ast.Source{
				Name:  file.Name(),
				Input: string(content),
			})
		}
	}

	schema, gqlErr := gqlparser.LoadSchema(sources...)
	if gqlErr != nil {
		return nil, fmt.Errorf("failed to load schema: %w", gqlErr)
	}

	return schema, nil
}

func ExportGraphQLToSchemaDump(schema *ast.Schema) (*SchemaDump, error) {
	dump := &SchemaDump{}

	// Identify all types that are not internal/scalars
	for _, def := range schema.Types {
		if def.BuiltIn {
			continue
		}
		// Skip Input, Enum, Scalar, Union, Interface for first pass (we'll use them as properties or relationships)
		if def.Kind != ast.Object {
			continue
		}
		// Skip Query, Mutation, Subscription
		if def.Name == "Query" || def.Name == "Mutation" || def.Name == "Subscription" {
			continue
		}
		// Skip Connection and Edge types (standard GraphQL pagination)
		if strings.HasSuffix(def.Name, "Connection") || strings.HasSuffix(def.Name, "Edge") {
			continue
		}

		// Check if it has an 'id' field to be considered a Node
		// Or if it is a known GUAC entity that should be a node
		hasID := false
		for _, field := range def.Fields {
			if field.Name == "id" {
				hasID = true
				break
			}
		}

		if !hasID && def.Name != "Artifact" {
			continue
		}

		node := Node{
			Table: def.Name,
		}

		for _, field := range def.Fields {
			// Check if it's a scalar or points to another type
			fieldType := field.Type.Name()
			fieldDef, ok := schema.Types[fieldType]

			if !ok {
				// Handle built-in scalars or custom scalars
				node.Properties = append(node.Properties, Property{
					Name:    field.Name,
					Type:    mapGraphQLTypeToKuzu(fieldType),
					Comment: field.Description,
				})
				continue
			}

			switch fieldDef.Kind {
			case ast.Scalar, ast.Enum:
				node.Properties = append(node.Properties, Property{
					Name:    field.Name,
					Type:    mapGraphQLTypeToKuzu(fieldType),
					Comment: field.Description,
				})
			case ast.Object, ast.Union, ast.Interface:
				// If it's an object, check if it has an ID
				isEntity := false
				if fieldDef.Kind == ast.Object {
					for _, f := range fieldDef.Fields {
						if f.Name == "id" {
							isEntity = true
							break
						}
					}
				} else {
					// Unions and Interfaces are generally collections of entities in GUAC
					isEntity = true
				}

				if isEntity {
					// It's a relationship
					toTables := getToTables(schema, fieldDef)
					for _, toTable := range toTables {
						node.Edges = append(node.Edges, Edge{
							ToTable: toTable,
							Symbol:  fmt.Sprintf("%s_%s", def.Name, field.Name),
							Column:  field.Name,
						})
					}
				} else {
					// It's a nested object without ID, treat as STRING (JSON)
					node.Properties = append(node.Properties, Property{
						Name:    field.Name,
						Type:    "STRING",
						Comment: field.Description,
					})
				}
			}
		}

		dump.Nodes = append(dump.Nodes, node)
	}

	// Sort nodes for deterministic output
	sort.Slice(dump.Nodes, func(i, j int) bool {
		return dump.Nodes[i].Table < dump.Nodes[j].Table
	})

	return dump, nil
}

func getToTables(schema *ast.Schema, def *ast.Definition) []string {
	switch def.Kind {
	case ast.Object:
		return []string{def.Name}
	case ast.Union:
		return def.Types
	case ast.Interface:
		var types []string
		for _, t := range schema.GetPossibleTypes(def) {
			types = append(types, t.Name)
		}
		return types
	default:
		return nil
	}
}

func mapGraphQLTypeToKuzu(gqlType string) string {
	switch gqlType {
	case "ID":
		return "STRING" // Or UUID if appropriate
	case "String":
		return "STRING"
	case "Int":
		return "INT64"
	case "Float":
		return "DOUBLE"
	case "Boolean":
		return "BOOL"
	case "Time":
		return "TIMESTAMP"
	default:
		return "STRING"
	}
}
