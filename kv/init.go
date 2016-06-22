package kv

import (
	"github.com/Cimpress-MCP/go-git2consul/repository"
	"gopkg.in/libgit2/git2go.v24"
)

// HandleInit handles initial fetching of the KV on start
func (h *KVHandler) HandleInit(repos []*repository.Repository) error {
	for _, repo := range repos {
		err := h.handleRepoInit(repo)
		if err != nil {
			return err
		}
	}

	return nil
}

// Handles differences on all branches of a repository, comparing the ref
// of the branch against the one in the KV
func (h *KVHandler) handleRepoInit(repo *repository.Repository) error {
	repo.Lock()
	defer repo.Unlock()

	itr, err := repo.NewReferenceIterator()
	if err != nil {
		return err
	}
	defer itr.Free()

	// Handle all local refs
	for {
		ref, err := itr.Next()
		if err != nil {
			break
		}

		b, err := ref.Branch().Name()
		if err != nil {
			return err
		}

		// Get only local refs
		if ref.IsRemote() == false {
			h.logger.Infof("KV GET ref: %s/%s", repo.Name(), b)
			kvRef, err := h.getKVRef(repo, b)
			if err != nil {
				return err
			}

			localRef := ref.Target().String()

			if len(kvRef) == 0 {
				// There is no ref in the KV, push the entire branch
				h.logger.Infof("KV PUT changes: %s/%s", repo.Name(), b)
				h.putBranch(repo, ref.Branch())

				h.logger.Infof("KV PUT ref: %s/%s", repo.Name(), b)
				h.putKVRef(repo, b)
			} else if kvRef != localRef {
				// Check if the ref belongs to that repo
				err := repo.CheckRef(kvRef)
				if err != nil {
					return err
				}

				// Handle modified and deleted files
				deltas, err := repo.DiffStatus(kvRef)
				if err != nil {
					return err
				}
				h.handleDeltas(repo, deltas)

				err = h.putKVRef(repo, b)
				if err != nil {
					return err
				}
				h.logger.Debugf("KV PUT ref: %s/%s", repo.Name(), b)
			}
		}
	}

	return nil
}

// Helper function that handles deltas
func (h *KVHandler) handleDeltas(repo *repository.Repository, deltas []git.DiffDelta) error {
	// Handle modified and deleted files
	for _, d := range deltas {
		switch d.Status {
		case git.DeltaRenamed:
			h.logger.Debugf("Detected renamed file: %s", d.NewFile.Path)
			h.logger.Infof("KV DEL %s/%s/%s", repo.Name(), repo.Branch(), d.OldFile.Path)
			err := h.deleteKV(repo, d.OldFile.Path)
			if err != nil {
				return err
			}
			h.logger.Infof("KV PUT %s/%s/%s", repo.Name(), repo.Branch(), d.NewFile.Path)
			err = h.putKV(repo, d.NewFile.Path)
			if err != nil {
				return err
			}
		case git.DeltaAdded, git.DeltaModified:
			h.logger.Debugf("Detected added/modified file: %s", d.NewFile.Path)
			h.logger.Infof("KV PUT %s/%s/%s", repo.Name(), repo.Branch(), d.NewFile.Path)
			err := h.putKV(repo, d.NewFile.Path)
			if err != nil {
				return err
			}
		case git.DeltaDeleted:
			h.logger.Debugf("Detected deleted file: %s", d.OldFile.Path)
			h.logger.Infof("KV DEL %s/%s/%s", repo.Name(), repo.Branch(), d.OldFile.Path)
			err := h.deleteKV(repo, d.OldFile.Path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
