package main

import (
	"encoding/json"
	"fmt"
	"os"

	"entgo.io/ent/schema/field"
	"github.com/guacsec/guac/pkg/assembler/backends/ent/migrate"
)

type Property struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Node struct {
	Table      string     `json:"table"`
	Properties []Property `json:"properties"`
}

type Edge struct {
	FromTable  string     `json:"from_table"`
	ToTable    string     `json:"to_table"`
	Column     string     `json:"column,omitempty"`
	Symbol     string     `json:"symbol"`
	Properties []Property `json:"properties,omitempty"`
}

type SchemaDump struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

func mapType(t field.Type) string {
	switch t {
	case field.TypeBool:
		return "BOOL"
	case field.TypeTime:
		return "TIMESTAMP"
	case field.TypeInt, field.TypeInt8, field.TypeInt16, field.TypeInt32, field.TypeInt64,
		field.TypeUint, field.TypeUint8, field.TypeUint16, field.TypeUint32, field.TypeUint64:
		return "INT64"
	case field.TypeFloat32, field.TypeFloat64:
		return "DOUBLE"
	case field.TypeString, field.TypeEnum:
		return "STRING"
	case field.TypeUUID:
		return "UUID"
	case field.TypeJSON:
		return "STRING" // Usually processed as JSON string or shredded
	default:
		return "STRING"
	}
}

func main() {
	dump := SchemaDump{}

	for _, table := range migrate.Tables {
		// Identify if it's a join table (Many-to-Many edge)
		isJoinTable := len(table.ForeignKeys) == 2 && len(table.Columns) == 2

		if isJoinTable {
			edge := Edge{
				FromTable: table.ForeignKeys[0].RefTable.Name,
				ToTable:   table.ForeignKeys[1].RefTable.Name,
				Symbol:    table.Name,
			}
			dump.Edges = append(dump.Edges, edge)
			continue
		}

		// Handle as a Node
		node := Node{
			Table: table.Name,
		}
		for _, col := range table.Columns {
			node.Properties = append(node.Properties, Property{
				Name: col.Name,
				Type: mapType(col.Type),
			})
		}
		dump.Nodes = append(dump.Nodes, node)

		// Extract edges from this table's foreign keys
		for _, fk := range table.ForeignKeys {
			edge := Edge{
				FromTable: table.Name,
				ToTable:   fk.RefTable.Name,
				Column:    fk.Columns[0].Name,
				Symbol:    fk.Symbol,
			}
			// If the table is an "observation" or relationship with properties,
			// include its non-FK properties in the edge.
			for _, col := range table.Columns {
				isFK := false
				for _, fkCol := range fk.Columns {
					if col.Name == fkCol.Name {
						isFK = true
						break
					}
				}
				if col.Name != "id" && !isFK {
					edge.Properties = append(edge.Properties, Property{
						Name: col.Name,
						Type: mapType(col.Type),
					})
				}
			}

			dump.Edges = append(dump.Edges, edge)
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(dump); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode schema dump: %v\n", err)
		os.Exit(1)
	}
}
