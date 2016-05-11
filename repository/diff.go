package repository

import "gopkg.in/libgit2/git2go.v24"

// Compares the current workdir with a target ref and return the modified files
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

	tree := &git.Tree{
		Object: *obj,
	}

	diffs, err := r.DiffTreeToWorkdir(tree, &git.DiffOptions{})
	if err != nil {
		return nil, err
	}

	n, err := diffs.NumDeltas()

	for i := 0; i < n; i++ {
		diff, err := diffs.GetDelta(i)
		if err != nil {
			return nil, err
		}
		deltas = append(deltas, diff)
	}

	return deltas, nil
}
