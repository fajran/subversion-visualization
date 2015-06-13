package svnvis

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
