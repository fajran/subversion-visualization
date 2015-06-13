package main

import (
	"fmt"
	"os"

	"github.com/fajran/subversion-visualization/svnvis"
)

func main() {
	p := svnvis.ParseLog(os.Stdin)
	for revlog := range p {
		fmt.Printf("%v\n", revlog)
	}
}
