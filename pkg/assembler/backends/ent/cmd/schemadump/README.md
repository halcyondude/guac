# GUAC Schema Dump Utility

This utility programmatically interrogates the GUAC Ent backend's metadata registry to produce a machine-readable, fully-typed definition of the entire graph schema. It facilitates automated migrations to graph databases like [RyuGraph](https://github.com/predictable-labs/ryugraph) (a fork of [kuzu](https://github.com/kuzudb/kuzu)), by providing definitive typing and relationship mapping.

## Overview

GUAC's Ent backend maintains a detailed schema in `pkg/assembler/backends/ent/migrate/schema.go`. This utility extracts:
- **Nodes**: All entity tables with their property names and static types.
- **Edges**: 
    - **One-to-Many**: Discovered via foreign key constraints.
    - **Many-to-Many**: Discovered via join tables in the Ent schema.
- **Static Typing**: Maps Ent Go types to standardized database types (e.g., `UUID`, `STRING`, `INT64`, `TIMESTAMP`, `DOUBLE`).

## Usage

To run the utility and generate a JSON schema dump:

```bash
go run pkg/assembler/backends/ent/cmd/schemadump/main.go > schema_dump.json
```

## How It Works

The utility (`main.go`) leverages the `github.com/guacsec/guac/pkg/assembler/backends/ent/migrate` package, which contains the runtime schema definition used by Ent.

1.  **Registry Iteration**: It iterates through the global `migrate.Tables` slice.
2.  **Node Identification**: Every table is treated as a potential Node. Columns are mapped to a `properties` list, with types translated from Go `field.Type` to static string representations (e.g., `field.TypeUUID` -> `"UUID"`).
3.  **Edge Discovery**:
    *   **Foreign Keys**: For every foreign key in a table, an Edge is created connecting the table to the referenced table.
    *   **Join Tables**: The utility intelligently identifies "pure" join tables (tables with exactly 2 foreign keys and no primary data columns). These are flattened into direct Edges between the two referenced entities, simplifying the graph model.
4.  **Property Extraction**: Edge properties (like "justification" or "timestamp" on a relationship) are preserved and added to the Edge definition.

## Applications

The `schema_dump.json` output is the foundational blueprint for a modern graph data pipeline.

### 1. DDL Generation for RyuGraph
[RyuGraph](https://github.com/predictable-labs/ryugraph) requires static schemas for its high-performance vector execution. You can use the JSON output to generate `CREATE NODE TABLE` and `CREATE REL TABLE` statements:

**Node Creation:**
```sql
-- Generated from "artifacts" node definition
CREATE NODE TABLE artifacts (
    id UUID, 
    algorithm STRING, 
    digest STRING, 
    PRIMARY KEY (id)
);
```

**Edge Creation:**
```sql
-- Generated from "dependencies" edge definition
CREATE REL TABLE dependencies (
    FROM package_versions TO package_names, 
    dependency_type STRING, 
    justification STRING
);
```

### 2. Ontological Documentation
The JSON dump serves as a machine-readable ontology. It can be fed into documentation generators to produce an interactive "Knowledge Graph" of the GUAC data model, allowing users to explore how `Packages`, `Artifacts`, and `Vulnerabilities` connect without reading source code.

### 3. Python ETL Pipeline (DuckDB to RyuGraph)
A robust ingestion pipeline can be built using Python and DuckDB to bridge the gap between GUAC's Postgres data and RyuGraph.

**Architecture:**
1.  **Extract**: Use DuckDB's high-performance Postgres scanner to read data directly from the running GUAC database.
2.  **Transform**: "Shred" complex structures into the node and edge formats defined by `schema_dump.json`.
3.  **Store**: Save the intermediate state as a single `.duckdb` file (preferred over multiple Parquet files for atomicity and simplified management).
4.  **Load**: Use RyuGraph's `COPY FROM` command to ingest the prepared data from DuckDB.

```python
# Pseudo-code for Python implementation
import duckdb
import json

# Load schema mapping
schema = json.load(open('schema_dump.json'))

# Connect to DuckDB and attach Postgres
con = duckdb.connect('guac_graph.duckdb')
con.execute("INSTALL postgres; LOAD postgres;")
con.execute("ATTACH 'dbname=guac user=postgres' AS pg (TYPE POSTGRES);")

# Dynamically generate INSERT SELECT statements based on schema
for node in schema['nodes']:
    # Generate SQL to move data from pg.public.<table> to DuckDB table
    pass
```

### 4. Advanced Graph Analysis (Cypher)
Once loaded into RyuGraph/KuzuDB, you can run powerful Cypher queries for supply chain insights.

**Scenario: Transitive Dependency Search**
Find all projects that depend on a specific library (e.g., `log4j`), potentially identifying impact blast radius.
```cypher
MATCH (root:Project)-[:DEPENDS_ON*]->(dep:Package)
WHERE dep.name = 'log4j'
RETURN root.name, dep.version
```

**Scenario: License Drift Detection**
Search for cases where the declared license of a package differs from the discovered license in its newer versions.
```cypher
MATCH (p:Package)-[:HAS_VERSION]->(v1:PackageVersion)-[:HAS_LICENSE]->(l1:License)
MATCH (p)-[:HAS_VERSION]->(v2:PackageVersion)-[:HAS_LICENSE]->(l2:License)
WHERE v2.timestamp > v1.timestamp AND l1.name <> l2.name
RETURN p.name, v1.version, l1.name, v2.version, l2.name
```

**Scenario: Visualizations**
The graph structure natively supports 2D/3D force-directed layout algorithms. Tools like Gephi or graph-app-kit can connect to the database to visualize the "constellation" of open source software, clustering projects by ecosystem (Go, Rust, JS) or connectivity.

**Scenario: CNCF Insights (Issue #1709)**
To address the need for ecosystem-wide insights:
*   **Vulnerability Prevalence**: `MATCH (v:Vulnerability)<-[:AFFECTS]-(p:Package) RETURN v.id, count(p) as impact_count ORDER BY impact_count DESC`
*   **SBOM Availability**: `MATCH (p:Package) OPTIONAL MATCH (p)-[:HAS_SBOM]->(s:SBOM) RETURN p.name, s IS NOT NULL as has_sbom`
