package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnionFindCmd(t *testing.T) {
	input := `{"edges":[["user1","user2"],["user2","user3"],["user4","user5"]]}`

	// Create pipes for stdin and stdout
	stdinR, stdinW, _ := os.Pipe()
	stdoutR, stdoutW, _ := os.Pipe()

	// Backup original stdin/stdout
	oldStdin := os.Stdin
	oldStdout := os.Stdout

	// Replace stdin/stdout
	os.Stdin = stdinR
	os.Stdout = stdoutW

	// Write input to stdin pipe
	go func() {
		defer func() { _ = stdinW.Close() }()
		_, _ = stdinW.WriteString(input)
	}()

	// Run the command
	cmd := &UnionFindCmd{}
	cmd.Run()

	// Close stdout writer and restore
	_ = stdoutW.Close()
	os.Stdin = oldStdin
	os.Stdout = oldStdout

	// Read output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(stdoutR)
	output := buf.String()
	t.Logf("Raw output: %q", output)

	// Parse output - should be a single line with JSON object containing array
	lines := strings.Split(strings.TrimSpace(output), "\n")
	require.Len(t, lines, 1) // Should have 1 line

	// Parse the wrapped result
	var wrappedResult struct {
		Result []UnionFindResult `json:"result"`
	}
	err := json.Unmarshal([]byte(lines[0]), &wrappedResult)
	require.NoError(t, err)
	resultArray := wrappedResult.Result
	require.Len(t, resultArray, 5) // Should have 5 unique values

	// Convert to map for easier testing
	results := make(map[string]string) // value -> root
	for _, result := range resultArray {
		results[result.Value] = result.Root
	}

	// Check that user1, user2, user3 are in same group
	assert.Equal(t, results["user1"], results["user2"])
	assert.Equal(t, results["user2"], results["user3"])

	// Check that user4, user5 are in same group (different from above)
	assert.Equal(t, results["user4"], results["user5"])
	assert.NotEqual(t, results["user1"], results["user4"])
}

func TestBipartiteUnionFindCmd(t *testing.T) {
	input := `{"relations":[["entity1","group100"],["entity1","group101"],["entity2","group101"],["entity3","group200"]]}`

	// Create pipes for stdin and stdout
	stdinR, stdinW, _ := os.Pipe()
	stdoutR, stdoutW, _ := os.Pipe()

	// Backup original stdin/stdout
	oldStdin := os.Stdin
	oldStdout := os.Stdout

	// Replace stdin/stdout
	os.Stdin = stdinR
	os.Stdout = stdoutW

	// Write input to stdin pipe
	go func() {
		defer func() { _ = stdinW.Close() }()
		_, _ = stdinW.WriteString(input)
	}()

	// Run the command
	cmd := &BipartiteUnionFindCmd{}
	cmd.Run()

	// Close stdout writer and restore
	_ = stdoutW.Close()
	os.Stdin = oldStdin
	os.Stdout = oldStdout

	// Read output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(stdoutR)
	output := buf.String()

	// Parse output - should be a single line with JSON object containing array
	lines := strings.Split(strings.TrimSpace(output), "\n")
	require.Len(t, lines, 1) // Should have 1 line

	// Parse the wrapped result
	var wrappedResult struct {
		Result []BipartiteResult `json:"result"`
	}
	err := json.Unmarshal([]byte(lines[0]), &wrappedResult)
	require.NoError(t, err)
	resultArray := wrappedResult.Result
	require.Len(t, resultArray, 3) // Should have 3 unique U values

	// Convert to map for easier testing
	results := make(map[string]string) // u -> v_root
	for _, result := range resultArray {
		results[result.U] = result.VRoot
	}

	// entity1 and entity2 should map to same V root
	// (because entity1 -> group100, group101 and entity2 -> group101)
	assert.Equal(t, results["entity1"], results["entity2"])

	// entity3 should map to different V root
	assert.NotEqual(t, results["entity1"], results["entity3"])
	assert.Equal(t, "group200", results["entity3"])
}

func TestEmptyInput(t *testing.T) {
	// Create pipes for stdin and stdout
	stdinR, stdinW, _ := os.Pipe()
	stdoutR, stdoutW, _ := os.Pipe()

	// Backup original stdin/stdout
	oldStdin := os.Stdin
	oldStdout := os.Stdout

	// Replace stdin/stdout
	os.Stdin = stdinR
	os.Stdout = stdoutW

	// Close stdin immediately (no input)
	_ = stdinW.Close()

	cmd := &UnionFindCmd{}
	cmd.Run()

	// Close stdout writer and restore
	_ = stdoutW.Close()
	os.Stdin = oldStdin
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(stdoutR)
	output := strings.TrimSpace(buf.String())

	// Should have no output for empty input
	assert.Empty(t, output)
}
