# TODO

## Initial version requirements:
* Better error handling of goroutines through errCh
* Possible usage of a runner to abstract git operations from the consul package
* Test coverage
* Better CI pipeline
* Webhook polling

## Future additions:
* File format backend
* Update on KV should be for modified and deleted files only
* Improve on mutex locks
* Run -once flag
* Support for source_root and mountpoint
* Support for tags as branches

## Bugs/Issues:
* Need to update diffs on the KV side
  * This includes only updating changed files
  * Delete untracked files
* If repositories slice is empty, stop the program
* Directory check has to check if it's a repository first

## Error handling:
* Should error on client connection be fatal/exit or should it just log an error?
  * Leaning towards logging and retry
* Better error handling on Loadrepos()
  * Bad configuration should be ignored

## Repo delta KV handling:
* On added, modified: PUT KV
* On delete: DEL KV
* On rename: DEL old KV followed by PUT new KV
