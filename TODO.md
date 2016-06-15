# TODO

## Initial version requirements:
* [x] Better error handling of goroutines through errCh
* [x] Possible usage of a runner to abstract git operations from the consul package
* [x] Update on KV should be for modified and deleted files only
* [x] Switch from godep to glide
* [x] Switch to apex/log
* [x] Webhook polling
  * [x] GitHub
  * [x] Stash
  * [x] Bitbucket
  * [x] Gitlab
* [x] Accept consul configuration for the client
* [x] Add -once flag to run git2consul once
* [ ] Better CD/CI pipeline
  * [ ] Cross-platform builds
  * [ ] Travis/appveyor

## Bugs/Issues:
* [x] Need to update diffs on the KV side
  * [x] This includes only updating changed files
  * [x] Delete untracked files
* [x] If repositories slice is empty, stop the program
* [x] Directory check has to check if it's a repository first
* [x] Runner, and watchers need a Stop() to handle cleanup better
* [x] Handle DoneCh better on both the watcher and runner
* [x] Handle initial load state better
  * [x] Watcher should handle initial changes from load state

## Error handling:
* [x] Better error handling on LoadRepos()
  * [x] Bad configuration should be ignored
* [ ] Handle repository error with git reset or re-clone

## Repo delta KV handling:
* [x] On added, modified: PUT KV
* [x] On delete: DEL KV
* [x] On rename: DEL old KV followed by PUT new KV

## Test coverage
* [ ] Repository
  * [x] New
  * [x] Clone
  * [x] Load
  * [ ] Pull
  * [ ] Ref
  * [x] Checkout
* [x] Config
  * [x] Load
* [ ] Runner
* [ ] Watch
  * [ ] Watcher
  * [ ] Interval
  * [ ] Webhook
* [ ] KV
  * [ ] Handler
  * [ ] Branch
  * [ ] KV
  * [ ] InitHandler
  * [ ] UpdateHandler

Test suite enhancement:
* [ ] git-init on repo should be done on init()
* [ ] Setup and teardown for each test during
  * [ ] Setup resets "remote" repo to initial commit
  * [ ] Teardown cleans local store

* Instead of testutil, we can use mocks to set up a mock repository.Repository object
