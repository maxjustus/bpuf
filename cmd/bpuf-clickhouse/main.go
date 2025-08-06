// Package main provides a ClickHouse User Defined Function (UDF) binary for union-find operations.
// It supports both standard union-find and bipartite union-find modes with JSONEachRow format.
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
)

//go:embed clickhouse_udfs.xml
var combinedUDFXML string

func main() {
	var (
		mode     = flag.String("mode", "unionfind", "UDF mode: 'unionfind' or 'bipartite'")
		printXML = flag.Bool("udf-xml", false, "Print ClickHouse UDF XML configuration for both modes")
	)
	flag.Parse()

	if *printXML {
		printClickHouseXML()
		return
	}

	switch *mode {
	case "unionfind":
		cmd := &UnionFindCmd{}
		cmd.Run()
	case "bipartite":
		cmd := &BipartiteUnionFindCmd{}
		cmd.Run()
	default:
		fmt.Fprintf(os.Stderr, "Unknown mode: %s\n", *mode)
		os.Exit(1)
	}
}

func printClickHouseXML() {
	fmt.Print(combinedUDFXML)
}
