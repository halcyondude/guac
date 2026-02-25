# Schema Export Utility

The `schemaexport` package provides high-level utilities to programmatically export the GUAC GraphQL ontology into external analytical and documentation formats.

It serves as a bridge between GUAC's operational graph and high-performance columnar graph engines like [KuzuDB](https://kuzudb.com) and [LadybugDB](https://ladybugdb.com), while also providing automated documentation generation.

## CLI Usage

The core functionality is exposed via the `guacone schemadump` subcommand.

### Examples

```bash
# Generate a JSON dump of the ontology
guacone schemadump --format=json > schema.json

# Generate KuzuDB/LadybugDB DDL (Cypher)
guacone schemadump --format=kuzu > schema.cypher

# Generate rich GitHub Flavored Markdown (GFM) documentation
guacone schemadump --format=markdown > ontology.md
```

### Options
- `--format`: Output format. Supported: `json`, `kuzu`, `markdown` (default `json`).
- `--schema-dir`: Path to the directory containing GUAC `.graphql` schema files (default `pkg/assembler/graphql/schema`).

## Architecture

This utility programmatically parses the GraphQL Abstract Syntax Tree (AST) using `gqlparser`. By deriving the graph ontology directly from the source `.graphql` files, the exporter remains resilient to changes in GUAC's data model without requiring manual code updates.

### Mapping Rules
- **Nodes**: GraphQL `type` definitions containing an `id: ID!` field (or explicitly identified entities like `Artifact`) are mapped to `NODE` tables.
- **Properties**: Scalar fields and Enums are mapped to node properties.
- **Edges**: Fields pointing to other object types, Unions, or Interfaces are mapped to `REL` tables.
- **Polymorphism**: Relationships involving Unions or Interfaces are accurately represented as `REL` tables with multiple `FROM`/`TO` pairs.

## Development & Testing

### Unit Tests
Verifies internal parsing, DDL generation, and Markdown formatting logic.
```bash
go test ./pkg/schemaexport/...
```

### LadybugDB Integration Test
Verifies that the generated DDL is natively compatible with LadybugDB by attempting to load the full GUAC schema into an embedded instance.

**Requirements:**
- Python 3.x
- `real_ladybug` Python library (`pip install real_ladybug`)

**Running the test:**
```bash
RUN_LADYBUG_TEST=true go test -v ./pkg/schemaexport/ladybug_test.go ./pkg/schemaexport/schemaexport.go ./pkg/schemaexport/kuzu.go
```
