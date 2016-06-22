package repository

import "gopkg.in/libgit2/git2go.v24"

// DiffStatus compares the current workdir with a target ref and return the modified files
func (r *Repository) DiffStatus(ref string) ([]git.DiffDelta, error) {
	deltas := []git.DiffDelta{}

	oid, err := git.NewOid(ref)
	if err != nil {
		return nil, err
	}

	// This can be for a different repo
	obj, err := r.Lookup(oid)
	if err != nil {
		return nil, err
	}

	commit, err := obj.AsCommit()
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	h, err := r.Head()
	if err != nil {
		return nil, err
	}

	obj2, err := r.Lookup(h.Target())
	if err != nil {
		return nil, err
	}

	commit2, err := obj2.AsCommit()
	if err != nil {
		return nil, err
	}

	tree2, err := commit2.Tree()
	if err != nil {
		return nil, err
	}

	do, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}

	diffs, err := r.DiffTreeToTree(tree, tree2, &do)
	if err != nil {
		return nil, err
	}

	n, err := diffs.NumDeltas()
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		diff, err := diffs.GetDelta(i)
		if err != nil {
			return nil, err
		}
		deltas = append(deltas, diff)
	}

	return deltas, nil
}
