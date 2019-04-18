/*
Copyright 2019 Kohl's Department Stores, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package watch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"

	"github.com/gorilla/mux"
	"gopkg.in/src-d/go-git.v4"
)

// GithubPayload is the response from GitHub
type GithubPayload struct {
	Ref string `json:"ref"`
}

// StashPayload is the response from Stash
type StashPayload struct {
	RefChanges []struct {
		RefID string `json:"refId"`
	} `json:"refChanges"`
}

// BitbucketPayload is the response Bitbucket
type BitbucketPayload struct {
	Push struct {
		Changes []struct {
			New struct {
				Name string `json:"name"`
			} `json:"new"`
		} `json:"changes"`
	} `json:"push"`
}

// GitLabPayload is the response from GitLab
type GitLabPayload struct {
	Ref string `json:"ref"`
}

func (w *Watcher) pollByWebhook(wg *sync.WaitGroup) {
	defer wg.Done()

	if w.once {
		return
	}

	errCh := make(chan error, 1)
	// Passing errCh instead of w.ErrCh to better handle watcher termination
	// since the caller can't determine what type of error it receives from watcher
	go w.ListenAndServe(errCh)

	for {
		select {
		case err := <-errCh:
			w.ErrCh <- err
			close(w.RcvDoneCh) // Stop the watcher if there is a
		case <-w.RcvDoneCh:
			return
		}
	}
}

// ListenAndServe starts the listener server for hooks
func (w *Watcher) ListenAndServe(errCh chan<- error) {
	r := mux.NewRouter()
	r.HandleFunc("/{repository}/github", w.githubHandler)
	r.HandleFunc("/{repository}/stash", w.stashHandler)
	r.HandleFunc("/{repository}/bitbucket", w.bitbucketHandler)
	r.HandleFunc("/{repository}/gitlab", w.gitlabHandler)

	addr := fmt.Sprintf("%s:%d", w.hookSvr.Address, w.hookSvr.Port)
	errCh <- http.ListenAndServe(addr, r)
}

// HTTP handler for github
func (w *Watcher) githubHandler(rw http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	repository := vars["repository"]

	eventType := rq.Header.Get("X-Github-Event")
	if eventType == "" {
		http.Error(rw, "Missing X-Github-Event header", http.StatusBadRequest)
		return
	}
	// Only process pusn events
	if eventType != "push" {
		return
	}

	body, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		http.Error(rw, "Cannot read body", http.StatusInternalServerError)
		return
	}

	payload := &GithubPayload{}
	err = json.Unmarshal(body, payload)
	if err != nil {
		http.Error(rw, "Cannot unmarshal JSON", http.StatusInternalServerError)
		return
	}

	// Check the content
	ref := payload.Ref
	if len(ref) == 0 {
		http.Error(rw, "ref is empty", http.StatusInternalServerError)
		return
	}
	if len(ref) <= 11 || ref[:11] != "refs/heads/" {
		return
	}

	branchName := ref[11:]

	i := sort.Search(len(w.Repositories), func(i int) bool {
		return w.Repositories[i].Name() == repository
	})

	// sort.Search could return last index if not found, so need to check once more
	if i == len(w.Repositories) || w.Repositories[i].Name() != repository {
		return
	}

	repo := w.Repositories[i]
	w.logger.WithField("repository", repo.Name()).Info("Received hook event from GitHub")

	err = repo.Pull(branchName)
	switch {
	case err == git.NoErrAlreadyUpToDate:
		w.logger.Debugf("Up to date: %s/%s", repo.Name(), branchName)
	case err == nil:
		w.logger.Infof("Changed: %s/%s", repo.Name(), branchName)
		w.RepoChangeCh <- repo
	case err != nil:
		w.logger.Errorf("Failed: %s/%s - %s", repo.Name(), branchName, err)
	}
}

// HTTP handler for Stash
func (w *Watcher) stashHandler(rw http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	repository := vars["repository"]

	body, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		http.Error(rw, "Cannot read body", http.StatusInternalServerError)
		return
	}

	payload := &StashPayload{}
	err = json.Unmarshal(body, payload)
	if err != nil {
		http.Error(rw, "Cannot unmarshal JSON", http.StatusInternalServerError)
		return
	}

	ref := payload.RefChanges[0].RefID

	if len(ref) == 0 {
		http.Error(rw, "ref is empty", http.StatusInternalServerError)
		return
	}
	if len(ref) <= 11 || ref[:11] != "refs/heads/" {
		return
	}

	branchName := ref[11:]

	i := sort.Search(len(w.Repositories), func(i int) bool {
		return w.Repositories[i].Name() == repository
	})

	// sort.Search could return last index if not found, so need to check once more
	if i == len(w.Repositories) || w.Repositories[i].Name() != repository {
		return
	}

	repo := w.Repositories[i]
	w.logger.WithField("repository", repo.Name()).Info("Received hook event from Stash")
	err = repo.Pull(branchName)
	switch {
	case err == git.NoErrAlreadyUpToDate:
		w.logger.Debugf("Up to date: %s/%s", repo.Name(), branchName)
	case err == nil:
		w.logger.Infof("Changed: %s/%s", repo.Name(), branchName)
		w.RepoChangeCh <- repo
	case err != nil:
		w.logger.Errorf("Failed: %s/%s - %s", repo.Name(), branchName, err)
	}
}

// HTTP handler for Bitbucket
func (w *Watcher) bitbucketHandler(rw http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	repository := vars["repository"]

	eventType := rq.Header.Get("X-Event-Key")
	if eventType == "" {
		http.Error(rw, "Missing X-Event-key header", http.StatusBadRequest)
		return
	}
	// Only process pusn events
	if eventType != "repo:push" {
		return
	}

	body, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		http.Error(rw, "Cannot read body", http.StatusInternalServerError)
		return
	}

	payload := &BitbucketPayload{}
	err = json.Unmarshal(body, payload)
	if err != nil {
		http.Error(rw, "Cannot unmarshal JSON", http.StatusInternalServerError)
		return
	}

	// Check the content
	ref := payload.Push.Changes[0].New.Name
	if len(ref) == 0 {
		http.Error(rw, "ref is empty", http.StatusInternalServerError)
		return
	}
	if len(ref) <= 11 || ref[:11] != "refs/heads/" {
		return
	}

	branchName := ref[11:]

	i := sort.Search(len(w.Repositories), func(i int) bool {
		return w.Repositories[i].Name() == repository
	})

	// sort.Search could return last index if not found, so need to check once more
	if i == len(w.Repositories) || w.Repositories[i].Name() != repository {
		return
	}

	repo := w.Repositories[i]
	w.logger.WithField("repository", repo.Name()).Info("Received hook event from Bitbucket")
	err = repo.Pull(branchName)
	switch {
	case err == git.NoErrAlreadyUpToDate:
		w.logger.Debugf("Up to date: %s/%s", repo.Name(), branchName)
	case err == nil:
		w.logger.Infof("Changed: %s/%s", repo.Name(), branchName)
		w.RepoChangeCh <- repo
	case err != nil:
		w.logger.Errorf("Failed: %s/%s - %s", repo.Name(), branchName, err)
	}
}

func (w *Watcher) gitlabHandler(rw http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	repository := vars["repository"]

	eventType := rq.Header.Get("X-Gitlab-Event")
	if eventType == "" {
		http.Error(rw, "Missing X-Gitlab-Event header", http.StatusBadRequest)
		return
	}
	// Only process pusn events
	if eventType != "Push Hook" {
		return
	}

	body, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		http.Error(rw, "Cannot read body", http.StatusInternalServerError)
		return
	}

	payload := &GitLabPayload{}
	err = json.Unmarshal(body, payload)
	if err != nil {
		http.Error(rw, "Cannot unmarshal JSON", http.StatusInternalServerError)
		return
	}

	// Check the content
	ref := payload.Ref
	if len(ref) == 0 {
		http.Error(rw, "ref is empty", http.StatusInternalServerError)
		return
	}
	if len(ref) <= 11 || ref[:11] != "refs/heads/" {
		return
	}

	branchName := ref[11:]

	i := sort.Search(len(w.Repositories), func(i int) bool {
		return w.Repositories[i].Name() == repository
	})

	// sort.Search could return last index if not found, so need to check once more
	if i == len(w.Repositories) || w.Repositories[i].Name() != repository {
		return
	}

	repo := w.Repositories[i]
	w.logger.WithField("repository", repo.Name()).Info("Received hook event from GitLab")
	err = repo.Pull(branchName)
	switch {
	case err == git.NoErrAlreadyUpToDate:
		w.logger.Debugf("Up to date: %s/%s", repo.Name(), branchName)
	case err == nil:
		w.logger.Infof("Changed: %s/%s", repo.Name(), branchName)
		w.RepoChangeCh <- repo
	case err != nil:
		w.logger.Errorf("Failed: %s/%s - %s", repo.Name(), branchName, err)
	}
}
