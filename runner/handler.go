package runner

import (
	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/repository"
)

func (r *Runner) initHandler(repo *repository.Repository) error {
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
			kvRef, err := r.getKVRef(repo, b)
			if err != nil {
				return err
			}

			localRef := ref.Target().String()

			if len(kvRef) == 0 || kvRef != localRef {
				log.Debugf("(consul) KV PUT changes for %s/%s", repo.Name(), b)
				r.putBranch(repo, ref.Branch())
				log.Debugf("(consul) KV PUT ref for %s/%s", repo.Name(), b)
				r.putKVRef(repo, b)
			}
		}
	}

	return nil
}

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

	if len(kvRef) == 0 || kvRef != localRef {
		log.Debugf("(consul) KV PUT changes for %s/%s", repo.Name(), b)
		err := r.putBranch(repo, h.Branch())
		if err != nil {
			return err
		}
		log.Debugf("(consul) KV PUT ref for %s/%s", repo.Name(), b)
		r.putKVRef(repo, b)
	}

	return nil
}
