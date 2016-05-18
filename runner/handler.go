package runner

import (
	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/repository"
	"gopkg.in/libgit2/git2go.v24"
)

// func (r *Runner) initHandler(repo *repository.Repository) error {
// 	repo.Lock()
// 	defer repo.Unlock()
//
// 	itr, err := repo.NewReferenceIterator()
// 	if err != nil {
// 		return err
// 	}
// 	defer itr.Free()
//
// 	// Handle all local refs
// 	for {
// 		ref, err := itr.Next()
// 		if err != nil {
// 			break
// 		}
//
// 		b, err := ref.Branch().Name()
// 		if err != nil {
// 			return err
// 		}
//
// 		// Get only local refs
// 		if ref.IsRemote() == false {
// 			log.Debugf("(consul) KV GET ref for %s/%s", repo.Name(), b)
// 			kvRef, err := r.getKVRef(repo, b)
// 			if err != nil {
// 				return err
// 			}
//
// 			localRef := ref.Target().String()
//
// 			// FIXME: Handle case wher ref comes from entirely different repo
// 			if len(kvRef) != 0 && kvRef != localRef {
// 				// Handle modified and deleted files
// 				deltas, err := repo.DiffStatus(kvRef)
// 				if err != nil {
// 					return err
// 				}
// 				for _, d := range deltas {
// 					log.Debugf(" (consul)(trace) Status: %s File: %s", d.Status, d.NewFile.Path)
// 					// switch d.Status {
// 					// case git.DeltaDeleted:
// 					// 	log.Debugf(" === Deleted file: %s", d.NewFile.Path)
// 					// case git.DeltaModified:
// 					// 	log.Debugf(" === Modified file: %s", d.NewFile.Path)
// 					// }
// 				}
// 			}
//
// 			if len(kvRef) == 0 || kvRef != localRef {
// 				log.Debugf("(consul) KV PUT changes for %s/%s", repo.Name(), b)
// 				r.putBranch(repo, ref.Branch())
// 				log.Debugf("(consul) KV PUT ref for %s/%s", repo.Name(), b)
// 				r.putKVRef(repo, b)
// 			}
// 		}
// 	}
//
// 	return nil
// }

func (r *Runner) updateHandler(repo *repository.Repository) error {
	repo.Lock()
	defer repo.Unlock()

	h, err := repo.Head()
	if err != nil {
		return err
	}
	b, err := h.Branch().Name()
	if err != nil {
		return err
	}

	log.Debugf("(consul) KV GET ref for %s/%s", repo.Name(), b)
	kvRef, err := r.getKVRef(repo, b)
	if err != nil {
		return err
	}

	// Local ref
	localRef := h.Target().String()
	log.Debugf("(consul) kvRef: %s | localRef: %s", kvRef, localRef)

	if len(kvRef) != 0 && kvRef != localRef {
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
		for _, d := range deltas {
			switch d.Status {
			case git.DeltaRenamed:
				log.Debugf("(runner)(trace) Renamed file: %s", d.NewFile.Path)
				err := r.deleteKV(repo, d.OldFile.Path)
				if err != nil {
					return err
				}
				err = r.putKV(repo, d.NewFile.Path)
				if err != nil {
					return err
				}
			case git.DeltaAdded, git.DeltaModified:
				log.Debugf("(runner)(trace) Added/Modified file: %s", d.NewFile.Path)
				err := r.putKV(repo, d.NewFile.Path)
				if err != nil {
					return err
				}
			case git.DeltaDeleted:
				log.Debugf("(runner)(trace) Deleted file: %s", d.OldFile.Path)
				err := r.deleteKV(repo, d.OldFile.Path)
				if err != nil {
					return err
				}
			}
		}
	}

	if len(kvRef) == 0 || kvRef != localRef {
		// TODO: Handle modified and deleted files

		// log.Debugf("(consul) KV PUT changes for %s/%s", repo.Name(), b)
		// err := r.putBranch(repo, h.Branch())
		// if err != nil {
		// 	return err
		// }
		log.Debugf("(consul) KV PUT ref for %s/%s", repo.Name(), b)
		r.putKVRef(repo, b)
	}

	return nil
}
