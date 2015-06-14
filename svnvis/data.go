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
	FromBranch   *string
	FromRevision int
}

func (r *RevisionLog) IsBranching() bool {
	b, _ := r.GetBranchingInfo()
	return b != nil
}

func (r *RevisionLog) GetBranchingInfo() (*string, int) {
	if len(r.Changes) == 0 {
		return nil, 0
	}

	c := r.Changes[0]
	return c.FromBranch, c.FromRevision
}
