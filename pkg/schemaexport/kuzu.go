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
	"sort"
	"strings"
)

func GenerateKuzuDDL(dump *SchemaDump) string {
	var sb strings.Builder

	var nodeDefinitions []string
	var edgeDefinitions []string

	// 1. Process Node Tables
	for _, node := range dump.Nodes {
		var cols []string

		// Ensure 'id' is present and is the primary key
		hasID := false
		for _, prop := range node.Properties {
			if prop.Name == "id" {
				hasID = true
				break
			}
		}

		for _, prop := range node.Properties {
			cols = append(cols, fmt.Sprintf("    %s %s", prop.Name, prop.Type))
		}

		if hasID {
			cols = append(cols, "    PRIMARY KEY (id)")
		}

		nodeDefinitions = append(nodeDefinitions, fmt.Sprintf("CREATE NODE TABLE %s (\n%s\n);", node.Table, strings.Join(cols, ",\n")))
	}

	// 2. Process Relationship Tables
	relTables := make(map[string][]Edge)
	relFromTo := make(map[string][]string) // symbol -> []"FROM T1 TO T2"

	for _, node := range dump.Nodes {
		for _, edge := range node.Edges {
			relTables[edge.Symbol] = append(relTables[edge.Symbol], edge)
			relFromTo[edge.Symbol] = append(relFromTo[edge.Symbol], fmt.Sprintf("FROM %s TO %s", node.Table, edge.ToTable))
		}
	}

	// Sort symbols for deterministic output
	var symbols []string
	for s := range relTables {
		symbols = append(symbols, s)
	}
	sort.Strings(symbols)

	for _, symbol := range symbols {
		edges := relTables[symbol]
		fromTos := relFromTo[symbol]
		
		var cols []string

		for _, ft := range fromTos {
			cols = append(cols, "    "+ft)
		}

		for _, prop := range edges[0].Properties {
			cols = append(cols, fmt.Sprintf("    %s %s", prop.Name, prop.Type))
		}

		edgeDefinitions = append(edgeDefinitions, fmt.Sprintf("CREATE REL TABLE %s (\n%s\n);", symbol, strings.Join(cols, ",\n")))
	}

	for _, d := range nodeDefinitions {
		sb.WriteString(d + "\n\n")
	}

	for _, d := range edgeDefinitions {
		sb.WriteString(d + "\n\n")
	}

	return sb.String()
}
