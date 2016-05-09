# go-git2consul

## Defaults

git2consul will attempt to deduce sane defaults for configuration. However, since git2consul needs to know which repository to pull from, minimal configuration is necessary.

### Repository-level configuration

| configuration  | default |
|----------------|---------|
| branches       | master  |
| hooks:type     | polling |
| hooks:interval | 60      |

### Development

go-git2consul is built with libgit2 v23, and therefore uses the matching git2go.v23 library.


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

## Bugs/Issues
* Clone performs checkout on all remote branches, not just the one specified
* consul/watch.go does not GET the ref for initial comparison
