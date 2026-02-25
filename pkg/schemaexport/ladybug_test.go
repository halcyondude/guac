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
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLadybugIntegration(t *testing.T) {
	if os.Getenv("RUN_LADYBUG_TEST") != "true" {
		t.Skip("Skipping Ladybug integration test. Set RUN_LADYBUG_TEST=true to run.")
	}

	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "ladybug-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 1. Generate DDL using our logic
	schema, err := LoadGraphQLSchema("../../pkg/assembler/graphql/schema")
	require.NoError(t, err)
	dump, err := ExportGraphQLToSchemaDump(schema)
	require.NoError(t, err)
	ddl := GenerateKuzuDDL(dump)

	ddlPath := filepath.Join(tmpDir, "schema.ddl")
	err = os.WriteFile(ddlPath, []byte(ddl), 0644)
	require.NoError(t, err)

	dbPath := filepath.Join(tmpDir, "test.db")

	// 2. Prepare Python script
	pythonScript := fmt.Sprintf(`
import real_ladybug as lbug
import os

def main():
    db = lbug.Database("%s")
    conn = lbug.Connection(db)
    
    with open("%s", "r") as f:
        ddl = f.read()
    
    # Kuzu/Ladybug DDL often contains multiple statements.
    # We split by semicolon and execute each.
    
    statements = ddl.split(";")
    for stmt in statements:
        clean_stmt = stmt.strip()
        if not clean_stmt:
            continue
        print(f"Executing statement: {clean_stmt[:50]}...")
        try:
            conn.execute(clean_stmt)
        except Exception as e:
            print(f"Error executing statement: {e}")
            print(f"Full statement: {clean_stmt}")
            raise e

    print("Successfully loaded DDL into LadybugDB")

if __name__ == "__main__":
    main()
`, dbPath, ddlPath)

	scriptPath := filepath.Join(tmpDir, "load_ddl.py")
	err = os.WriteFile(scriptPath, []byte(pythonScript), 0644)
	require.NoError(t, err)

	// 3. Run Python script
	// Ensure real_ladybug is installed: pip install real_ladybug
	cmd := exec.CommandContext(ctx, "python3", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Python output: %s", string(output))
	}
	assert.NoError(t, err, "Python script failed to load DDL into LadybugDB")
	assert.Contains(t, string(output), "Successfully loaded DDL into LadybugDB")
}
