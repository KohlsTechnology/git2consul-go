package repository

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/libgit2/git2go.v24"
)

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

	h, err := r.Head()
	if err != nil {
		return nil, err
	}
	obj2, err := r.Lookup(h.Target())
	if err != nil {
		return nil, err
	}

	tree2 := &git.Tree{
		Object: *obj2,
	}

	diffs, err := r.DiffTreeToTree(tree, tree2, &git.DiffOptions{})
	if err != nil {
		return nil, err
	}

	log.Debugf("(git)(trace) Diffs from func: %+v | Repo ref: %s | Diff ref: %s\n", diffs, h.Target().String(), ref)

	stats, err := diffs.Stats()
	if err != nil {
		return nil, err
	}
	log.Debugf("(git)(trace) Diffs files changed from func: %d\n", stats.FilesChanged())

	n, err := diffs.NumDeltas()
	log.Debugf("(git)(trace) Diffs num from func: %d\n", n)

	for i := 0; i < n; i++ {
		diff, err := diffs.GetDelta(i)
		if err != nil {
			return nil, err
		}
		deltas = append(deltas, diff)
	}

	return deltas, nil
}
