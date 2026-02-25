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

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/table"
)

func GenerateMarkdown(dump *SchemaDump) string {
	var sb strings.Builder

	sb.WriteString("# GUAC Ontology Documentation\n\n")
	sb.WriteString("This document provides a detailed overview of the GUAC graph ontology, including entities (nodes) and their relationships (edges).\n\n")

	// Table of Contents
	sb.WriteString("## Table of Contents\n\n")
	l := list.NewWriter()
	l.SetStyle(list.StyleBulletCircle)

	l.AppendItem("[Entities](#entities)")
	l.Indent()
	for _, node := range dump.Nodes {
		l.AppendItem(fmt.Sprintf("[%s](#%s)", node.Table, strings.ToLower(node.Table)))
	}
	l.UnIndent()

	l.AppendItem("[Relationships](#relationships)")
	l.Indent()
	
	relMap := make(map[string][]string)
	for _, node := range dump.Nodes {
		for _, edge := range node.Edges {
			relMap[edge.Symbol] = append(relMap[edge.Symbol], fmt.Sprintf("%s to %s", node.Table, edge.ToTable))
		}
	}
	var relSymbols []string
	for s := range relMap {
		relSymbols = append(relSymbols, s)
	}
	sort.Strings(relSymbols)

	for _, s := range relSymbols {
		l.AppendItem(fmt.Sprintf("[%s](#%s)", s, strings.ToLower(s)))
	}
	l.UnIndent()

	sb.WriteString(l.RenderMarkdown() + "\n\n")

	// Entities Section
	sb.WriteString("## Entities\n\n")
	for _, node := range dump.Nodes {
		sb.WriteString(fmt.Sprintf("### %s\n\n", node.Table))
		if node.Description != "" {
			sb.WriteString(node.Description + "\n\n")
		}

		sb.WriteString("#### Properties\n\n")
		tw := table.NewWriter()
		tw.AppendHeader(table.Row{"Property", "Type", "Description"})
		for _, prop := range node.Properties {
			desc := prop.Description
			if desc == "" {
				desc = "-"
			}
			tw.AppendRow(table.Row{prop.Name, prop.Type, desc})
		}
		sb.WriteString(tw.RenderMarkdown() + "\n\n")

		if len(node.Edges) > 0 {
			sb.WriteString("#### Outgoing Relationships\n\n")
			el := list.NewWriter()
			el.SetStyle(list.StyleBulletCircle)
			for _, edge := range node.Edges {
				el.AppendItem(fmt.Sprintf("[%s](#%s) to [%s](#%s)", 
					edge.Symbol, strings.ToLower(edge.Symbol), 
					edge.ToTable, strings.ToLower(edge.ToTable)))
			}
			sb.WriteString(el.RenderMarkdown() + "\n\n")
		}
	}

	// Relationships Section
	sb.WriteString("## Relationships\n\n")
	for _, symbol := range relSymbols {
		sb.WriteString(fmt.Sprintf("### %s\n\n", symbol))
		
		// Find description from one of the edges
		var desc string
		for _, node := range dump.Nodes {
			for _, edge := range node.Edges {
				if edge.Symbol == symbol && edge.Description != "" {
					desc = edge.Description
					break
				}
			}
			if desc != "" {
				break
			}
		}
		if desc != "" {
			sb.WriteString(desc + "\n\n")
		}

		sb.WriteString("#### Definitions\n\n")
		dl := list.NewWriter()
		dl.SetStyle(list.StyleBulletCircle)
		for _, def := range relMap[symbol] {
			dl.AppendItem(def)
		}
		sb.WriteString(dl.RenderMarkdown() + "\n\n")
	}

	return sb.String()
}
