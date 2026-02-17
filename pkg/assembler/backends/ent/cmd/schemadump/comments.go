package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// SchemaComments maps StructName -> FieldName -> Comment
type SchemaComments map[string]map[string]string

func parseSchemaComments(schemaPath string) (SchemaComments, error) {
	result := make(SchemaComments)

	files, err := os.ReadDir(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema dir: %w", err)
	}

	fset := token.NewFileSet()

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".go") || strings.HasSuffix(f.Name(), "_test.go") {
			continue
		}

		path := filepath.Join(schemaPath, f.Name())
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			continue // Skip unparseable files
		}

		// We look for methods named "Fields" on structs.
		// func (Type) Fields() []ent.Field
		for _, decl := range node.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok || funcDecl.Name.Name != "Fields" || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
				continue
			}

			// Extract receiver type name
			recvType := funcDecl.Recv.List[0].Type
			var structName string
			if ident, ok := recvType.(*ast.Ident); ok {
				structName = ident.Name
			} else {
				continue
			}

			if _, exists := result[structName]; !exists {
				result[structName] = make(map[string]string)
			}

			// Parse the return statement to find field definitions
			// return []ent.Field{ ... }
			if funcDecl.Body == nil {
				continue
			}

			for _, stmt := range funcDecl.Body.List {
				ret, ok := stmt.(*ast.ReturnStmt)
				if !ok || len(ret.Results) == 0 {
					continue
				}

				compLit, ok := ret.Results[0].(*ast.CompositeLit)
				if !ok {
					continue
				}

				// Iterate over elements in the slice: field.String("name").Comment("foo")
				for _, elt := range compLit.Elts {
					parseFieldExpr(elt, result[structName])
				}
			}
		}
	}

	return result, nil
}

func parseFieldExpr(expr ast.Expr, fieldMap map[string]string) {
	// We are looking for chains like:
	// field.String("name").Comment("the comment")
	// This appears as a CallExpr.

	var fieldName, comment string

	// Traverse the call chain
	// The AST for `field.String("x").Comment("y")` is a CallExpr (Comment) wrapped around a CallExpr (String).

	curr := expr
	for {
		call, ok := curr.(*ast.CallExpr)
		if !ok {
			break
		}

		// Check selector: .Comment("...") or .String("...")
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			break // Could be a simple call like field.String("x") which is SelectorExpr, but wrapped in CallExpr
		}

		method := sel.Sel.Name

		if method == "Comment" && len(call.Args) > 0 {
			if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
				comment = strings.Trim(lit.Value, "\"")
			}
		}

		// Identify the field definition call: field.Type("name")
		// Usually the start of the chain has the package name "field" (or a different alias) as the X.
		// field.String("name") -> X is Ident "field", Sel is "String"
		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "field" {
			// This is the base call, e.g. field.String("name")
			if len(call.Args) > 0 {
				if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
					fieldName = strings.Trim(lit.Value, "\"")
				}
			}
		}

		// Move down the chain to the inner expression
		// field.String("x").Comment("y") -> Call(Comment).Fun(Selector).X is Call(String)
		curr = sel.X
	}

	if fieldName != "" && comment != "" {
		fieldMap[fieldName] = comment
	}
}

// FindComments attempts to find the field comments for a given table name.
// It handles simple singular/plural mismatch heuristics.
func (sc SchemaComments) FindComments(tableName string) map[string]string {
	// 1. Try exact match (unlikely as Go structs are PascalCase)
	if fields, ok := sc[tableName]; ok {
		return fields
	}

	// 2. iterate and convert to snake_case
	// naive snake case conversion: BillOfMaterials -> bill_of_materials
	for structName, fields := range sc {
		snake := toSnakeCase(structName)
		// Exact match: bill_of_materials == bill_of_materials
		if snake == tableName {
			return fields
		}
		// Plural match: artifacts == artifact + s
		if tableName == snake+"s" {
			return fields
		}
		// Plural match with 'es': vulnerabilities == vulnerability + es
		// (ent uses a library for this, we use simple heuristics)
		if strings.HasSuffix(tableName, "s") && strings.HasPrefix(tableName, snake) {
			// heuristic: if table name starts with snake version of struct, assume it's the one
			// e.g. certify_vulns starts with certify_vuln
			return fields
		}
		// Special case: "ies" suffix? certify_vulnerabilities vs certify_vulnerability
		// Not needed for current schema likely.
	}

	return nil
}

func toSnakeCase(s string) string {
	var res []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			res = append(res, '_')
		}
		res = append(res, r)
	}
	return strings.ToLower(string(res))
}
