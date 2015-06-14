package svnvis

import (
	"encoding/json"
	"io"
	"strings"
)

type Project struct {
	Name string
	Logs []*RevisionLog

	Prefix string
}

func NewProject(name string, logs []*RevisionLog) *Project {
	return NewProjectWithPrefix(name, "", logs)
}

func NewProjectWithPrefix(name, prefix string, logs []*RevisionLog) *Project {
	return &Project{
		Name:   name,
		Logs:   logs,
		Prefix: prefix,
	}
}

func GetBranchName(prefix, path string) string {
	t := strings.Split(strings.Trim(path, "/"), "/")
	if len(t) == 0 {
		return ""
	}

	if t[0] == prefix {
		t = t[1:]
	}
	if len(t) == 0 {
		return ""
	}

	if t[0] == "branches" {
		t = t[1:]
	}
	if len(t) == 0 {
		return ""
	}

	return t[0]
}

func (p *Project) GetBranchName(revlog *RevisionLog) string {
	changes := revlog.Changes
	if len(changes) == 0 {
		return ""
	}

	c := changes[0]
	path := c.Path

	return GetBranchName(p.Prefix, path)
}

func (p *Project) WriteToJson(w io.Writer) error {
	type branching struct {
		Branch   string `json:"branch"`
		Revision int    `json:"revision"`
	}
	type revision struct {
		Revision  int        `json:"revision"`
		Author    string     `json:"author"`
		Timestamp string     `json:"timestamp"`
		Branch    string     `json:"branch"`
		Branching *branching `json:"branching,omitempty"`
	}
	var data struct {
		Revisions []revision `json:"revisions"`
	}

	data.Revisions = make([]revision, 0)
	for _, r := range p.Logs {
		rev := revision{
			Revision:  r.Revision,
			Author:    r.Author,
			Timestamp: r.Timestamp,
			Branch:    p.GetBranchName(r),
		}
		b, rr := r.GetBranchingInfo()
		if b != nil {
			rev.Branching = &branching{
				Branch:   GetBranchName(p.Prefix, *b),
				Revision: rr,
			}
		}
		data.Revisions = append(data.Revisions, rev)
	}

	e := json.NewEncoder(w)
	return e.Encode(&data)
}
