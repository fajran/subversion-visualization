package main

import (
	"bufio"
	"fmt"
	"os"
	"io"
	"strconv"
	"strings"
)

type RevisionLog struct {
	Revision  int
	Author    string
	Timestamp string
	Info      string
	Changes   []Change
	Log       string
}

type Change struct {
	Type         string
	Path         string
	FromBranch   string
	FromRevision int
}

func isLogSeparator(line string) bool {
	t := strings.Trim(line, "-")
	return len(line) > 0 && len(t) == 0
}

func isEmptyLines(lines []string) bool {
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if len(line) > 0 {
			return false
		}
	}
	return true
}

func splitLogLines(lines <-chan string) <-chan []string {
	c := make(chan []string, 1)

	go func() {
		group := make([]string, 0)
		for line := range lines {
			line = strings.TrimRight(line, " ")

			if isLogSeparator(line) {
				if !isEmptyLines(group) {
					c <- group
				}

				group = make([]string, 0)

			} else {
				group = append(group, line)
			}
		}

		if !isEmptyLines(group) {
			c <- group
		}

		close(c)
	}()

	return c
}

func parseInt(str string) int {
	i, _ := strconv.ParseInt(str, 10, 0)
	return int(i)
}

func parseChange(line string) Change {
	line = strings.Trim(line, " ")
	changeType := string(line[0])

	path := line[2:]
	branch := ""
	rev := 0

	t := " (from"
	pos := strings.Index(path, t)
	if pos != -1 {
		pos2 := strings.LastIndex(path, ":")
		branch = path[pos+len(t)+1 : pos2]
		rev = parseInt(path[pos2+1 : len(path)-1])
		path = path[:pos]
	}

	return Change{
		Type:         changeType,
		Path:         path,
		FromBranch:   branch,
		FromRevision: rev,
	}
}

func parseLog(groups <-chan []string) <-chan *RevisionLog {
	c := make(chan *RevisionLog, 1)

	go func() {
		for lines := range groups {
			header := lines[0]
			p := strings.Split(header, "|")

			revision := parseInt(strings.Trim(p[0], " ")[1:])
			author := strings.Trim(p[1], " ")
			timestamp := strings.Trim(p[2], " ")
			info := strings.Trim(p[3], " ")

			changes := make([]Change, 0)
			logLines := make([]string, 0)

			phase := 0
			for _, line := range lines[2:] {
				if phase == 0 && len(strings.Trim(line, " ")) == 0 {
					phase = 1

				} else if phase == 0 {
					change := parseChange(line)
					changes = append(changes, change)

				} else if phase == 1 {
					logLines = append(logLines, line)
				}
			}

			log := strings.Trim(strings.Join(logLines, "\n"), " \n")

			c <- &RevisionLog{
				Revision:  revision,
				Author:    author,
				Timestamp: timestamp,
				Info:      info,
				Changes:   changes,
				Log:       log,
			}
		}
		close(c)
	}()

	return c
}

func readLines(r io.Reader) <-chan string {
	c := make(chan string)

	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			line := s.Text()
			c <- line
		}
		close(c)
	}()

	return c
}

func ParseLog(r io.Reader) <-chan *RevisionLog {
	lines := readLines(r)
	groups := splitLogLines(lines)
	return parseLog(groups)
}

func main() {
	for revlog := range ParseLog(os.Stdin) {
		fmt.Printf("%v\n", revlog)
	}
}
