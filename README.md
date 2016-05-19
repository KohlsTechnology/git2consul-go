# go-git2consul

## Default configuration

git2consul will attempt to use sane defaults for configuration. However, since git2consul needs to know which repository to pull from, minimal configuration is necessary.

| configuration        | required | default        |
|----------------------|----------|----------------|
| local_store          | no       | `os.TempDir()` |
| repos:name           | yes      |                |
| repos:url            | yes      |                |
| repos:branches       | no       | master         |
| repos:hooks:type     | no       | polling        |
| repos:hooks:interval | no       | 60             |

## Development dependencies:
* Go 1.6
* libgit2 v0.24.0
* [glide](https://github.com/Masterminds/glide)
