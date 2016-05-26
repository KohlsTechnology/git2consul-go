package watch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"gopkg.in/libgit2/git2go.v24"
)

type GithubPayload struct {
	Ref string `json:"ref"`
}

type StashPayload struct {
	RefChanges []struct {
		RefId string `json:"refId"`
	} `json:"refChanges"`
}

type BitbucketPayload struct {
	Push struct {
		Changes []struct {
			New struct {
				Name string `json:"name"`
			} `json:"new"`
		} `json:"changes"`
	} `json:"push"`
}

type GitLabPayload struct {
	Ref string `json:"ref"`
}

func (w *Watcher) pollByWebhook() {
	errCh := make(chan error, 1)
	// Passing errCh instead of w.ErrCh to better handle watcher termination
	// since the caller can't determine what type of error it receives from watcher
	go w.ListenAndServe(errCh)

	for {
		select {
		case err := <-errCh:
			w.ErrCh <- err
			w.Stop() // Stop the watcher if it's unable to serve
		case <-w.DoneCh:
			return
		}
	}
}

func (w *Watcher) ListenAndServe(errCh chan<- error) {
	r := mux.NewRouter()
	r.HandleFunc("/{repository}/github", w.githubHandler)
	r.HandleFunc("/{repository}/stash", w.stashHandler)
	r.HandleFunc("/{repository}/bitbucket", w.bitbucketHandler)
	r.HandleFunc("/{repository}/gitlab", w.gitlabHandler)

	addr := fmt.Sprintf(":%d", w.webhookPort)
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

	w.logger.Info("Received hook event from GitHub")
	repo := w.Repositories[i]
	analysis, err := repo.Pull(branchName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// If there is a change, send the repo RepoChangeCh
	switch {
	case analysis&git.MergeAnalysisUpToDate != 0:
		w.logger.Debugf("Up to date: %s/%s", repo.Name(), branchName)
	case analysis&git.MergeAnalysisNormal != 0, analysis&git.MergeAnalysisFastForward != 0:
		w.logger.Infof("Changed: %s/%s", repo.Name(), branchName)
		w.RepoChangeCh <- repo
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

	ref := payload.RefChanges[0].RefId

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

	w.logger.Info("Received hook event from Stash")
	repo := w.Repositories[i]
	analysis, err := repo.Pull(branchName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// If there is a change, send the repo RepoChangeCh
	switch {
	case analysis&git.MergeAnalysisUpToDate != 0:
		w.logger.Debugf("Up to date: %s/%s", repo.Name(), branchName)
	case analysis&git.MergeAnalysisNormal != 0, analysis&git.MergeAnalysisFastForward != 0:
		w.logger.Infof("Changed: %s/%s", repo.Name(), branchName)
		w.RepoChangeCh <- repo
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

	w.logger.Info("Received hook event from Bitbucket")
	repo := w.Repositories[i]
	analysis, err := repo.Pull(branchName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// If there is a change, send the repo RepoChangeCh
	switch {
	case analysis&git.MergeAnalysisUpToDate != 0:
		w.logger.Debugf("Up to date: %s/%s", repo.Name(), branchName)
	case analysis&git.MergeAnalysisNormal != 0, analysis&git.MergeAnalysisFastForward != 0:
		w.logger.Infof("Changed: %s/%s", repo.Name(), branchName)
		w.RepoChangeCh <- repo
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

	w.logger.Info("Received hook event from GitLab")
	repo := w.Repositories[i]
	analysis, err := repo.Pull(branchName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// If there is a change, send the repo RepoChangeCh
	switch {
	case analysis&git.MergeAnalysisUpToDate != 0:
		w.logger.Debugf("Up to date: %s/%s", repo.Name(), branchName)
	case analysis&git.MergeAnalysisNormal != 0, analysis&git.MergeAnalysisFastForward != 0:
		w.logger.Infof("Changed: %s/%s", repo.Name(), branchName)
		w.RepoChangeCh <- repo
	}
}
