package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fajran/subversion-visualization/svnvis"
)

func main() {
	mode := "parse"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	if mode == "parse" {
		p := svnvis.ParseLog(os.Stdin)
		for revlog := range p {
			fmt.Printf("%v\n", revlog)
		}

	} else if mode == "json" {
		prefix := ""
		if len(os.Args) > 2 {
			prefix = os.Args[2]
		}

		logs := svnvis.ParseLogToList(os.Stdin)
		p := svnvis.NewProjectWithPrefix("project name", prefix, logs)
		err := p.WriteToJson(os.Stdout)
		if err != nil {
			log.Fatalf("Error writing to JSON", err)
		}
	}
}
