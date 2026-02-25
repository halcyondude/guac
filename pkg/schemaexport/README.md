# Schema Export Utility

This package provides utilities to programmatically export the GUAC GraphQL ontology into external formats, such as KuzuDB DDL (Cypher) and JSON. 

It is designed to bridge the gap between GUAC's operational graph (Graph-based metadata) and high-performance analytical graph databases (columnar graph engines like [KuzuDB](https://kuzudb.com) and [LadybugDB](https://ladybugdb.com)).

## The `guacone schemadump` Command

The logic in this package powers the `guacone schemadump` command. This utility interrogates the GUAC GraphQL schema files to produce a machine-readable definition of the entire graph schema.

### Usage

```bash
# Generate a JSON dump of the ontology
guacone schemadump --format=json > schema.json

# Generate KuzuDB/LadybugDB DDL (Cypher)
guacone schemadump --format=kuzu > schema.cypher
```

### Options
- `--format`: Output format. Supported: `json`, `kuzu` (default `json`).
- `--schema-dir`: Directory containing GUAC `.graphql` schema files (default `pkg/assembler/graphql/schema`).

## Architecture

This utility parses the GraphQL Abstract Syntax Tree (AST) directly. This ensures that the exporter remains resilient and automatically updates whenever the GUAC maintainers update the core ontology in the `.graphql` files.

### Mapping Rules
- **Nodes**: GraphQL `type` definitions that contain an `id: ID!` field (or are explicitly marked as entities like `Artifact`) are mapped to Kuzu `NODE` tables.
- **Properties**: Scalar fields (String, Int, Boolean, etc.) and Enums are mapped to node properties.
- **Edges**: Fields that point to other object types, Unions, or Interfaces are mapped to Kuzu `REL` tables. 
- **Unions/Interfaces**: Relationships involving Unions or Interfaces result in `REL` tables with multiple `FROM`/`TO` pairs, accurately representing the polymorphic nature of GUAC links (e.g., `IsOccurrence` linking to either a `Package` or a `Source`).

## Testing

### Unit Tests
Standard Go unit tests verify the parsing and DDL generation logic.
```bash
go test ./pkg/schemaexport/...
```

### LadybugDB Integration Test
There is an integration test that verifies the generated DDL is natively compatible with LadybugDB. It uses a Python bridge to load the schema into a real embedded database instance.

**Requirements:**
- Python 3
- `real_ladybug` library: `pip install real_ladybug`

**Run Integration Test:**
```bash
RUN_LADYBUG_TEST=true go test -v ./pkg/schemaexport/ladybug_test.go ./pkg/schemaexport/schemaexport.go ./pkg/schemaexport/kuzu.go
```
