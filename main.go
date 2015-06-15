package main

import (
	"os"

	"github.com/fajran/subversion-visualization/cli"
)

func main() {
	mode := "parse"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	if mode == "parse" {
		cli.Parse()

	} else if mode == "json" {
		cli.Json()

	} else if mode == "edgelist" {
		cli.EdgeList()

	}
}
