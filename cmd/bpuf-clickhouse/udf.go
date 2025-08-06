package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/maxjustus/bpuf/unionfind"
)

type UnionFindCmd struct{}

type UnionFindPair struct {
	A string `json:"a"`
	B string `json:"b"`
}

type UnionFindResult struct {
	Value string `json:"value"`
	Root  string `json:"root"`
}

type BipartiteUnionFindCmd struct{}

type BipartiteRelation struct {
	U string `json:"u"`
	V string `json:"v"`
}

type BipartiteResult struct {
	U     string `json:"u"`
	VRoot string `json:"v_root"`
}

// Helper function to read a full line, handling potential partial reads.
func readLine(reader *bufio.Reader) (string, error) {
	var line []byte
	for {
		part, isPrefix, err := reader.ReadLine()
		if err != nil {
			return "", err // Return any errors (including EOF)
		}
		// Append the part read to the line buffer
		line = append(line, part...)
		if !isPrefix {
			break // If isPrefix is false, we have read the entire line
		}
	}
	return string(line), nil
}

// processLines handles the common pattern of reading lines, processing them, and outputting results
func processLines(processor func(line string) (interface{}, error)) {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := readLine(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("ERROR:", err)
			continue
		}

		if line == "" {
			continue
		}

		results, err := processor(line)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}

		// Output result wrapped in object with field name matching return_name in XML
		output := map[string]interface{}{
			"result": results,
		}
		resultJSON, err := json.Marshal(output)
		if err != nil {
			fmt.Println("ERROR:", err)
			continue
		}
		fmt.Println(string(resultJSON))
	}
}

func (c *UnionFindCmd) Run() {
	processLines(func(line string) (interface{}, error) {
		// Parse input from ClickHouse: {"edges":[["a","b"],["c","d"]]}
		var input struct {
			Edges [][2]string `json:"edges"`
		}
		if err := json.Unmarshal([]byte(line), &input); err != nil {
			return nil, fmt.Errorf("could not parse input: %v", err)
		}

		// Convert arrays to UnionFindPair structs
		pairs := make([]UnionFindPair, len(input.Edges))
		for i, edge := range input.Edges {
			pairs[i] = UnionFindPair{A: edge[0], B: edge[1]}
		}

		// Build union-find structure
		uf := unionfind.NewUnionFindWithValues[string](len(pairs) * 2)
		for _, pair := range pairs {
			uf.Union(pair.A, pair.B)
		}

		// Collect all unique values from input pairs
		valueSet := make(map[string]bool)
		for _, pair := range pairs {
			valueSet[pair.A] = true
			valueSet[pair.B] = true
		}

		// Build array of results
		results := make([]UnionFindResult, 0, len(valueSet))
		for value := range valueSet {
			root := uf.FindReturningValue(value)
			results = append(results, UnionFindResult{
				Value: value,
				Root:  root,
			})
		}

		return results, nil
	})
}

func (c *BipartiteUnionFindCmd) Run() {
	processLines(func(line string) (interface{}, error) {
		// Parse input from ClickHouse: {"relations":[["u","v"],["x","y"]]}
		var input struct {
			Relations [][2]string `json:"relations"`
		}
		if err := json.Unmarshal([]byte(line), &input); err != nil {
			return nil, fmt.Errorf("could not parse input: %v", err)
		}

		// Convert arrays to BipartiteRelation structs
		relations := make([]BipartiteRelation, len(input.Relations))
		for i, rel := range input.Relations {
			relations[i] = BipartiteRelation{U: rel[0], V: rel[1]}
		}

		// Build bipartite union-find structure
		buf := unionfind.NewBipartiteUnionFindWithValues[string, string](len(relations) * 2)
		for _, relation := range relations {
			buf.Union(relation.U, relation.V)
		}

		// Collect all unique U values
		uSet := make(map[string]bool)
		for _, relation := range relations {
			uSet[relation.U] = true
		}

		// Build array of results
		results := make([]BipartiteResult, 0, len(uSet))
		for u := range uSet {
			vRoot, exists := buf.FindVRootForU(u)
			if !exists {
				// This shouldn't happen, but handle gracefully
				continue
			}
			results = append(results, BipartiteResult{
				U:     u,
				VRoot: vRoot,
			})
		}

		return results, nil
	})
}
