package main

import (
	"fmt"
	"strings"
)

func generateMarkdown(dump SchemaDump) {
	fmt.Println("# GUAC Schema Documentation")
	fmt.Println()

	fmt.Println("## Entities")
	for _, node := range dump.Nodes {
		fmt.Printf("### %s\n", node.Table)

		if len(node.NaturalKeys) > 0 {
			fmt.Printf("#### %s: Natural Keys\n", node.Table)
			fmt.Println("| Keys |")
			fmt.Println("| --- |")
			for _, key := range node.NaturalKeys {
				fmt.Printf("| [%s] |\n", strings.Join(key, ", "))
			}
			fmt.Println()
		}

		fmt.Printf("#### %s: Fields\n", node.Table)
		fmt.Println("| Field | Type | Comment |")
		fmt.Println("| --- | --- | --- |")
		for _, prop := range node.Properties {
			fmt.Printf("| %s | %s | %s |\n", prop.Name, prop.Type, prop.Comment)
		}
		fmt.Println()

		if len(node.Edges) > 0 {
			fmt.Printf("#### %s: Edges\n", node.Table)
			fmt.Println("| To | Edge Field | On Delete |")
			fmt.Println("| --- | --- | --- |")
			for _, edge := range node.Edges {
				fmt.Printf("| %s | %s | %s |\n", edge.ToTable, edge.Symbol, edge.Column)
			}
			fmt.Println()
		}
	}
}
