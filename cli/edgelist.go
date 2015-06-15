package cli

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/fajran/subversion-visualization/svnvis"
)

func EdgeList() {
	prefix := ""
	if len(os.Args) > 2 {
		prefix = os.Args[2]
	}

	logs := svnvis.ParseLogToList(os.Stdin)
	p := svnvis.NewProjectWithPrefix("project name", prefix, logs)

	branches := make(map[string][]*svnvis.RevisionLog)

	for _, r := range p.Logs {
		b := p.GetBranchName(r)
		if b == "" {
			continue
		}
		branches[b] = append(branches[b], r)
	}

	for b := range branches {
		revs := make(map[int]*svnvis.RevisionLog)
		nums := make([]int, 0)
		for _, r := range branches[b] {
			revs[r.Revision] = r
			nums = append(nums, r.Revision)
		}

		sort.Ints(nums)
		prev := 0
		for _, n := range nums {
			r := revs[n]
			bb, br := r.GetBranchingInfo()
			if bb != nil {
				fmt.Printf("%d %d %s\n", br, n, p.GetBranchName(r))
			} else {
				fmt.Printf("%d %d %s\n", prev, n, p.GetBranchName(r))
			}
			prev = n
		}
	}

	log.Printf("done!\n")
}
