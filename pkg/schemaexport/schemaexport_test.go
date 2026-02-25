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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportGraphQLToSchemaDump(t *testing.T) {
	// We can use the actual schema in the repo for testing
	schema, err := LoadGraphQLSchema("../../pkg/assembler/graphql/schema")
	assert.NoError(t, err)
	assert.NotNil(t, schema)

	dump, err := ExportGraphQLToSchemaDump(schema)
	assert.NoError(t, err)
	assert.NotNil(t, dump)

	// Check for some known nodes
	foundArtifact := false
	for _, node := range dump.Nodes {
		if node.Table == "Artifact" {
			foundArtifact = true
			break
		}
	}
	assert.True(t, foundArtifact, "Artifact node not found in dump")
}

func TestGenerateKuzuDDL(t *testing.T) {
	dump := &SchemaDump{
		Nodes: []Node{
			{
				Table: "Person",
				Properties: []Property{
					{Name: "id", Type: "STRING"},
					{Name: "name", Type: "STRING"},
				},
				Edges: []Edge{
					{ToTable: "Company", Symbol: "Person_worksAt", Column: "worksAt"},
				},
			},
			{
				Table: "Company",
				Properties: []Property{
					{Name: "id", Type: "STRING"},
					{Name: "name", Type: "STRING"},
				},
			},
		},
	}

	ddl := GenerateKuzuDDL(dump)
	assert.Contains(t, ddl, "CREATE NODE TABLE Person")
	assert.Contains(t, ddl, "CREATE NODE TABLE Company")
	assert.Contains(t, ddl, "CREATE REL TABLE Person_worksAt (FROM Person TO Company)")
}
