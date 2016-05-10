## TODO

### Initial version requirements:
* Better error handling of goroutines through errCh
* Possible usage of a runner to abstract git operations from the consul package
* Test coverage
* Better CI pipeline
* Webhook polling

### Future additions:
* File format backend
* Update on KV should be for modified and deleted files only
* Improve on mutex locks
* Run -once flag
* Support for source_root and mountpoint
* Support for tags as branches

### Bugs/Issues:
* Clone performs checkout on all remote branches, not just the one specified
* OpenRepository() needs to check on repo.Url as well
* If repositories slice is empty, stop the program

### Error handling:
* Should error on client connection be fatal/exit or should it just log an error?
  * Leaning towards logging and retry
* Better error handling on Loadrepos()
  * Bad configuration should be ignored
