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

package cmd

import (
	stdjson "encoding/json"
	"fmt"
	"os"

	"github.com/guacsec/guac/pkg/schemaexport"
	"github.com/spf13/cobra"
)

var schemadumpOptions struct {
	format    string
	schemaDir string
}

var schemadumpCmd = &cobra.Command{
	Use:   "schemadump",
	Short: "dump the GUAC schema in various formats (JSON, Kuzu DDL, Markdown)",
	Long: `schemadump is a utility to export the GUAC GraphQL ontology into external formats.
It programmatically parses the GraphQL schema definitions to ensure compatibility with 
external analytical tools and documentation formats.

Supported formats:
  - json: A machine-readable representation of the graph ontology.
  - kuzu: Data Definition Language (DDL) for KuzuDB and LadybugDB.
  - markdown: Rich, human-readable documentation of entities and relationships.`,
	Example: `  # Export schema as Kuzu DDL
  guacone schemadump --format=kuzu > schema.cypher

  # Generate Markdown documentation
  guacone schemadump --format=markdown > ontology.md

  # Export as JSON from a custom schema directory
  guacone schemadump --format=json --schema-dir=./my-schemas`,
	Run: func(cmd *cobra.Command, args []string) {
		schema, err := schemaexport.LoadGraphQLSchema(schemadumpOptions.schemaDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load schema: %v\n", err)
			os.Exit(1)
		}

		dump, err := schemaexport.ExportGraphQLToSchemaDump(schema)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to export schema: %v\n", err)
			os.Exit(1)
		}

		switch schemadumpOptions.format {
		case "json":
			encoder := stdjson.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(dump); err != nil {
				fmt.Fprintf(os.Stderr, "failed to encode JSON: %v\n", err)
				os.Exit(1)
			}
		case "kuzu":
			ddl := schemaexport.GenerateKuzuDDL(dump)
			fmt.Print(ddl)
		case "markdown":
			md := schemaexport.GenerateMarkdown(dump)
			fmt.Print(md)
		default:
			fmt.Fprintf(os.Stderr, "unknown format: %s\n", schemadumpOptions.format)
			os.Exit(1)
		}
	},
}

func init() {
	schemadumpCmd.Flags().StringVar(&schemadumpOptions.format, "format", "json", "output format: json, kuzu, markdown")
	schemadumpCmd.Flags().StringVar(&schemadumpOptions.schemaDir, "schema-dir", "pkg/assembler/graphql/schema", "directory containing GraphQL schema files")
	rootCmd.AddCommand(schemadumpCmd)
}
