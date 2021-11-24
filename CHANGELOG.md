Starting with release v0.1.2 See the [releases page](https://github.com/KohlsTechnology/git2consul-go/releases).

## v0.1.1 (January 6, 2020)
FEATURES:
* None

CHANGES:
* Update to Go 1.13 and switch to Go Modules [GH-16](https://github.com/KohlsTechnology/git2consul-go/pull/16)
* Add git branch, git commit, Go version, and build date to version output [GH-19](https://github.com/KohlsTechnology/git2consul-go/pull/19)
* Enable gitter for this repo [GH-23](https://github.com/KohlsTechnology/git2consul-go/pull/23)

BUG FIXES:
* None

## v0.1.0
#
#### Summary
Release v0.1.0 covers the basic functionality with added features enhancing path management of the keys in the Consul KV Store and authentication for private repositories (f.e. GitHub Enterprise).
#### Features

* Added atomicity to the Consul transactions
* Added mountpoint option which sets the prefix for added keys
* Added skip_branch option which skipps the branch name for the added keys
* Expand YAML file content to k:v
* Added source_root option which allows to point the root of the repository from which we want to process the data
* Added authentication (basic and ssh)

#### Bug fixes

* Fix for transactions - limited to 64 elements chunks
* Fixed reference "not found" issue on branch pull
