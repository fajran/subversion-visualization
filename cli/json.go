package cli

import (
	"log"
	"os"

	"github.com/fajran/subversion-visualization/svnvis"
)

func Json() {
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
