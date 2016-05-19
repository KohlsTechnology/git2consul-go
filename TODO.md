# TODO

## Initial version requirements:
* [x] Better error handling of goroutines through errCh
* [x] Possible usage of a runner to abstract git operations from the consul package
* [x] Update on KV should be for modified and deleted files only
* [x] Switch from godep to glide
* [] Better CI pipeline
* [] Webhook polling

## Future additions:
* File format backend
* Improve on mutex locks
* Run -once flag
* Support for source_root and mountpoint
* Support for tags as branches

## Bugs/Issues:
* [x] Need to update diffs on the KV side
  * [x] This includes only updating changed files
  * [x] Delete untracked files
* [x] If repositories slice is empty, stop the program
* [x] Directory check has to check if it's a repository first

## Error handling:
* [x] Better error handling on LoadRepos()
  * [x] Bad configuration should be ignored

## Repo delta KV handling:
* [x] On added, modified: PUT KV
* [x] On delete: DEL KV
* [x] On rename: DEL old KV followed by PUT new KV

## Test coverage
* [] Repository
  * [x] Clone
  * [x] Load
  * [x] Poll
  * [] Pull
  * [] Ref
  * [] Checkout
* [x] Config
  * [x] Load
* [] Runner
  * [] Watch
  * [] KV
  * [] Handler
