package kv

import (
	"github.com/apex/log"
	"github.com/cleung2010/go-git2consul/repository"
	"gopkg.in/libgit2/git2go.v24"
)

// Function that handles initial fetching of the KV on start
func (h *KVHandler) HandleInit(repos []*repository.Repository) error {
	for _, repo := range repos {
		err := h.handleRepoInit(repo)
		if err != nil {
			return err
		}
	}

	return nil
}

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
			log.Debugf("(consul) KV GET ref for %s/%s", repo.Name(), b)
			kvRef, err := h.getKVRef(repo, b)
			if err != nil {
				return err
			}

			localRef := ref.Target().String()

			if len(kvRef) == 0 {
				// There is no ref in the KV, push the entire branch
				log.Debugf("(consul)(trace) KV PUT changes for %s/%s", repo.Name(), b)
				h.putBranch(repo, ref.Branch())

				log.Debugf("(consul) KV PUT ref for %s/%s", repo.Name(), b)
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
				log.Debugf("(consul) KV PUT ref for %s/%s", repo.Name(), b)
			}
		}
	}

	return nil
}

func (h *KVHandler) handleDeltas(repo *repository.Repository, deltas []git.DiffDelta) error {
	// Handle modified and deleted files
	for _, d := range deltas {
		switch d.Status {
		case git.DeltaRenamed:
			log.Debugf("(runner)(trace) Renamed file: %s", d.NewFile.Path)
			err := h.deleteKV(repo, d.OldFile.Path)
			if err != nil {
				return err
			}
			err = h.putKV(repo, d.NewFile.Path)
			if err != nil {
				return err
			}
		case git.DeltaAdded, git.DeltaModified:
			log.Debugf("(runner)(trace) Added/Modified file: %s", d.NewFile.Path)
			err := h.putKV(repo, d.NewFile.Path)
			if err != nil {
				return err
			}
		case git.DeltaDeleted:
			log.Debugf("(runner)(trace) Deleted file: %s", d.OldFile.Path)
			err := h.deleteKV(repo, d.OldFile.Path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
