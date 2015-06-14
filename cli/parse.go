package cli

import (
	"fmt"
	"os"

	"github.com/fajran/subversion-visualization/svnvis"
)

func Parse() {
	p := svnvis.ParseLog(os.Stdin)
	for revlog := range p {
		fmt.Printf("%v\n", revlog)
	}
}
