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

		// Fields Table
		fmt.Println("#### Fields")
		fmt.Println("| Field | Type | Comment |")
		fmt.Println("| --- | --- | --- |")
		for _, prop := range node.Properties {
			fmt.Printf("| %s | %s | %s |\n", prop.Name, prop.Type, prop.Comment)
		}
		fmt.Println()

		// Edges Table
		if len(node.Edges) > 0 {
			fmt.Println("#### Edges")
			fmt.Println("| To | Edge Field | On Delete |")
			fmt.Println("| --- | --- | --- |")
			for _, edge := range node.Edges {
				// We don't have exact OnDelete or Field info from `Edge` struct directly as it was simplified.
				// But we can show Symbol and Column.
				// User example: | To | Symbol | Column |
				fmt.Printf("| %s | %s | %s |\n", edge.ToTable, edge.Symbol, edge.Column)
			}
			fmt.Println()
		}

		// Natural Keys Table
		if len(node.NaturalKeys) > 0 {
			fmt.Println("#### Natural Keys (Composite Unique Constraints)")
			fmt.Println("| Keys |")
			fmt.Println("| --- |")
			for _, key := range node.NaturalKeys {
				fmt.Printf("| `[%s]` |\n", strings.Join(key, ", "))
			}
			fmt.Println()
		}
	}
}
