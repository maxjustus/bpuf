package e2e

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUDFBinaryStandalone(t *testing.T) {
	// This test doesn't require ClickHouse - just tests our UDF binary directly
	projectRoot, err := findProjectRoot()
	require.NoError(t, err)

	buildCmd := exec.Command("go", "build", "-o", "bin/bpuf-clickhouse", "./cmd/bpuf-clickhouse")
	buildCmd.Dir = projectRoot
	err = buildCmd.Run()
	require.NoError(t, err)

	testData := `{"edges":[["item1","item2"],["item2","item3"],["item4","item5"]]}`

	// Test unionfind mode
	cmd := exec.Command("./bin/bpuf-clickhouse", "--mode=unionfind")
	cmd.Dir = projectRoot
	cmd.Stdin = strings.NewReader(testData)

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "UDF failed: %s", string(output))

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	require.Len(t, lines, 1, "Expected 1 output line with JSON array")

	// Should be valid JSON object with result field
	assert.True(t, strings.HasPrefix(lines[0], "{"))
	assert.True(t, strings.HasSuffix(lines[0], "}"))
	assert.Contains(t, lines[0], `"result"`)
	assert.Contains(t, lines[0], `"value"`)
	assert.Contains(t, lines[0], `"root"`)

	// Test bipartite mode
	bipartiteData := `{"relations":[["x1","c100"],["x1","c101"],["x2","c101"]]}`

	cmd = exec.Command("./bin/bpuf-clickhouse", "--mode=bipartite")
	cmd.Dir = projectRoot
	cmd.Stdin = strings.NewReader(bipartiteData)

	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "Bipartite UDF failed: %s", string(output))

	lines = strings.Split(strings.TrimSpace(string(output)), "\n")
	require.Len(t, lines, 1, "Expected 1 output line with JSON array")

	// Should be valid JSON object with result field
	assert.True(t, strings.HasPrefix(lines[0], "{"))
	assert.True(t, strings.HasSuffix(lines[0], "}"))
	assert.Contains(t, lines[0], `"result"`)
	assert.Contains(t, lines[0], `"u"`)
	assert.Contains(t, lines[0], `"v_root"`)
}

func TestXMLGeneration(t *testing.T) {
	projectRoot, err := findProjectRoot()
	require.NoError(t, err)

	buildCmd := exec.Command("go", "build", "-o", "bin/bpuf-clickhouse", "./cmd/bpuf-clickhouse")
	buildCmd.Dir = projectRoot
	err = buildCmd.Run()
	require.NoError(t, err)

	// Test XML generation
	cmd := exec.Command("./bin/bpuf-clickhouse", "--udf-xml")
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "XML generation failed: %s", string(output))

	xmlStr := string(output)

	// Should contain both UDF definitions
	assert.Contains(t, xmlStr, `<name>unionFind</name>`)
	assert.Contains(t, xmlStr, `<name>bipartiteUnionFind</name>`)
	assert.Contains(t, xmlStr, `--mode=unionfind`)
	assert.Contains(t, xmlStr, `--mode=bipartite`)
	assert.Contains(t, xmlStr, `<format>JSONEachRow</format>`)

	// Should be valid XML structure
	assert.Contains(t, xmlStr, `<?xml version="1.0"?>`)
	assert.Contains(t, xmlStr, `<functions>`)
	assert.Contains(t, xmlStr, `</functions>`)
}
