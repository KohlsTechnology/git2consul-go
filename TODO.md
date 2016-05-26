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
  * [ ] Travis tests

## Bugs/Issues:
* [x] Need to update diffs on the KV side
  * [x] This includes only updating changed files
  * [x] Delete untracked files
* [x] If repositories slice is empty, stop the program
* [x] Directory check has to check if it's a repository first
* [x] Runner, and watchers need a Stop() to handle cleanup better

## Error handling:
* [x] Better error handling on LoadRepos()
  * [x] Bad configuration should be ignored

## Repo delta KV handling:
* [x] On added, modified: PUT KV
* [x] On delete: DEL KV
* [x] On rename: DEL old KV followed by PUT new KV

## Test coverage
* [ ] Repository
  * [x] Clone
  * [x] Load
  * [x] Poll
  * [ ] Pull
  * [ ] Ref
  * [ ] Checkout
* [x] Config
  * [x] Load
* [ ] Runner
  * [ ] Watch
  * [ ] KV
  * [ ] Handler

## Webhook polling
Will be served from a single port, and different repos will be given different endpoints
E.g. test-example will have its optional webhook endpoint at:
* localhost:<port>/test-example/github
* localhost:<port>/test-example/stash
* localhost:<port>/test-example/bitbucket
* localhost:<port>/test-example/gitlab
