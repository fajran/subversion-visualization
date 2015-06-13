package main

import (
	"fmt"
	"os"

	"svnvis"
)

func main() {
	p := svnvis.ParseLog(os.Stdin)
	for revlog := range p {
		fmt.Printf("%v\n", p)
	}
}
