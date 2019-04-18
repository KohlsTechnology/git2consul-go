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
