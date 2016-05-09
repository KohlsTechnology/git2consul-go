# go-git2consul

## Defaults

git2consul will attempt to deduce sane defaults for configuration. However, since git2consul needs to know which repository to pull from, minimal configuration is necessary.

### Repository-level configuration

| configuration  | default |
|----------------|---------|
| branches       | master  |
| hooks:type     | polling |
| hooks:interval | 60      |

## Development

### Build dependencies:
* Go 1.6
* libgit2 v0.23.x
* git2go.v23
