package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"entgo.io/ent/schema/field"
	"github.com/guacsec/guac/pkg/assembler/backends/ent/migrate"
)

type Property struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Comment string `json:"comment,omitempty"`
}

type Node struct {
	Table       string     `json:"table"`
	Properties  []Property `json:"properties"`
	NaturalKeys [][]string `json:"natural_keys,omitempty"`
	Edges       []Edge     `json:"edges,omitempty"`
}

type Edge struct {
	ToTable    string     `json:"to_table"`
	Column     string     `json:"column,omitempty"`
	Symbol     string     `json:"symbol"`
	Properties []Property `json:"properties,omitempty"`
}

type SchemaDump struct {
	Nodes []Node `json:"nodes"`
}

func mapType(t field.Type) string {
	switch t {
	case field.TypeBool:
		return "BOOL"
	case field.TypeTime:
		return "TIMESTAMP"
	case field.TypeInt64, field.TypeInt: // Go int is usually 64-bit on modern systems
		return "INT64"
	case field.TypeInt32:
		return "INT32"
	case field.TypeInt16:
		return "INT16"
	case field.TypeInt8:
		return "INT8"
	case field.TypeUint, field.TypeUint64:
		return "UINT64"
	case field.TypeUint32:
		return "UINT32"
	case field.TypeUint16:
		return "UINT16"
	case field.TypeUint8:
		return "UINT8"
	case field.TypeFloat32:
		return "FLOAT"
	case field.TypeFloat64:
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
	outputFormat := flag.String("o", "json", "output format: json, markdown, or kuzu")
	schemaPath := flag.String("schema-path", "pkg/assembler/backends/ent/schema", "path to the ent schema directory")
	flag.Parse()

	dump := SchemaDump{}

	// Attempt to parse comments from schema files
	comments, err := parseSchemaComments(*schemaPath)
	if err != nil {
		// Log warning but continue without comments
		fmt.Fprintf(os.Stderr, "warning: failed to parse schema comments from %s: %v\n", *schemaPath, err)
	}

	// Pass 1: Create Nodes (skipping join tables)
	// We map table name to *Node index in dump.Nodes to easily append edges later.
	nodeMap := make(map[string]int)

	for _, table := range migrate.Tables {
		// Identify if it's a join table (Many-to-Many edge)
		if len(table.ForeignKeys) == 2 && len(table.Columns) == 2 {
			continue // Skip join tables in pass 1
		}

		node := Node{
			Table: table.Name,
		}

		fieldComments := comments.FindComments(table.Name)

		// Identify FK columns and map them to their semantic name (simplified for natural keys)
		fkMap := make(map[string]string)
		for _, fk := range table.ForeignKeys {
			for _, col := range fk.Columns {
				// Heuristic: strip "_id" suffix for better readability in natural keys
				// e.g. "package_id" -> "package", "name_id" -> "name"
				semanticName := strings.TrimSuffix(col.Name, "_id")
				fkMap[col.Name] = semanticName
			}
		}

		for _, col := range table.Columns {
			comment := ""
			if fieldComments != nil {
				comment = fieldComments[col.Name]
			}
			node.Properties = append(node.Properties, Property{
				Name:    col.Name,
				Type:    mapType(col.Type),
				Comment: comment,
			})
		}

		// Extract natural keys from unique indexes
		for _, idx := range table.Indexes {
			if idx.Unique {
				var key []string
				for _, col := range idx.Columns {
					if semanticName, isFK := fkMap[col.Name]; isFK {
						key = append(key, semanticName)
					} else {
						key = append(key, col.Name)
					}
				}
				node.NaturalKeys = append(node.NaturalKeys, key)
			}
		}

		dump.Nodes = append(dump.Nodes, node)
		nodeMap[table.Name] = len(dump.Nodes) - 1
	}

	// Pass 2: Extract Edges (from FKs and Join Tables)
	for _, table := range migrate.Tables {
		// Handle Join Tables (Many-to-Many)
		if len(table.ForeignKeys) == 2 && len(table.Columns) == 2 {
			fromTable := table.ForeignKeys[0].RefTable.Name
			toTable := table.ForeignKeys[1].RefTable.Name

			edge := Edge{
				ToTable: toTable,
				Symbol:  table.Name,
			}

			if idx, ok := nodeMap[fromTable]; ok {
				dump.Nodes[idx].Edges = append(dump.Nodes[idx].Edges, edge)
			}
			continue
		}

		// Handle Foreign Keys (One-to-Many / Many-to-One)
		if idx, ok := nodeMap[table.Name]; ok {
			for _, fk := range table.ForeignKeys {
				edge := Edge{
					ToTable: fk.RefTable.Name,
					Column:  fk.Columns[0].Name,
					Symbol:  fk.Symbol,
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
						comment := ""
						fieldComments := comments.FindComments(table.Name)
						if fieldComments != nil {
							comment = fieldComments[col.Name]
						}
						edge.Properties = append(edge.Properties, Property{
							Name:    col.Name,
							Type:    mapType(col.Type),
							Comment: comment,
						})
					}
				}
				dump.Nodes[idx].Edges = append(dump.Nodes[idx].Edges, edge)
			}
		}
	}

	// No need for re-parsing flags here, they are parsed at the beginning
	switch *outputFormat {
	case "markdown":
		generateMarkdown(dump)
	case "kuzu":
		generateKuzuDDL(dump)
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(dump); err != nil {
			fmt.Fprintf(os.Stderr, "failed to encode schema dump: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown output format: %s\n", *outputFormat)
		os.Exit(1)
	}
}
