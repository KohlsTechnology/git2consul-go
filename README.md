[![Join the chat at https://gitter.im/KohlsTechnology/git2consul-go](https://badges.gitter.im/KohlsTechnology/git2consul-go.svg)](https://gitter.im/KohlsTechnology/git2consul-go?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.com/KohlsTechnology/git2consul-go.svg?branch=master)](https://travis-ci.com/KohlsTechnology/git2consul-go)
[![codecov](https://codecov.io/gh/KohlsTechnology/git2consul-go/branch/master/graph/badge.svg)](https://codecov.io/gh/KohlsTechnology/git2consul-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/KohlsTechnology/git2consul-go)](https://goreportcard.com/report/github.com/KohlsTechnology/git2consul-go)

## Project Deprication and Archival

**As of April 22, 2022 this project is no longer maintained. This repository is being archived(marked as read only).**

# git2consul-go

The git2consul-go tool is used to populate a [Consul](https://www.consul.io) key/value store from a git repo.

The baseline source code was forked from [go-git2consul](https://github.com/Cimpress-MCP/go-git2consul) which was
inspired by the orginal [git2consul](https://github.com/breser/git2consul) tool.

## Improvements Over NodeJS git2consul
* uses the official Consul Go Lang client library
* uses native Go Lang git implementation [go-git](https://github.com/src-d/go-git/)
* removal of nodejs and git runtime dependencies
* configuration is sourced locally instead of it being fetched from the Consul K/V
* transaction atomicity implies the set of keys is stored either entirely or not at all. Along with atomicity the number of the KV API calls is limited. However there is a pending [issue](https://github.com/hashicorp/consul/issues/2921) as Consul transaction endpoint can handle only 64 items in the payload. The transactions are executed in 64 elements chunks.

## Installation

git2consul-go comes in two variants:
* as a single binary file which after downloading can be placed in any working directory - either on the workstation (from which git2consul will be executed) or on the Consul node (depends whether access to the git repository is available from the Consul nodes or not)
* as a source code that can be build on the user workstation ([How to build from src?](#compiling-from-source))

## Documentation

### Example
Simple example usage.

```
$ git2consul -config config.json -basic -user mygituser -password mygitpass -once
```

Simple example config file.
```
{
  "repos": [
    {
      "name": "example",
      "url": "http://github.com/DummyOrg/ExampleRepo.git"
    }
  ]
}
```

### Command Line Options

```
$ git2consul -help
Usage of git2consul:
  -config string
        path to config file
  -debug
        enable debugging mode
  -logfmt string
        specify log format [text | json]  (default "text")
  -once
        run git2consul once and exit
  -version
        show version
```

### Configuration

Configuration is provided with a JSON file and passed in via the `-config` flag. Repository
configuration will take care of cloning the repository into `local_store`, but it will not
be responsible for creating the actual `local_store` directory. Similarly, it is expected
that there is no collision of directory or file that contains the same name as the repository
name under `local_store`, or git2consul will exit with an error. If there is a git repository
under a specified repo name, and the origin URL is different from the one provided in the
configuration, it will be overwritten.

#### Default configuration

git2consul will attempt to use sane defaults for configuration. However, since git2consul needs to know which repository to pull from, minimal configuration is necessary.


| Configuration             | Required | Default Value  | Available Values                           | Description
|---------------------------|----------|----------------|--------------------------------------------| -----------
| local_store               | no       | `os.TempDir()` | `string`                                   | Local cache for git2consul to store its tracked repositories
| webhook:address           | no       |                | `string`                                   | Webhook listener address that git2consul will be using
| webhook:port              | no       | 9000           | `int`                                      | Webhook listener port that git2consul will be using
| repos:name                | yes      |                | `string`                                   | Name of the repository. This will match the webhook path, if any are enabled
| repos:url                 | yes      |                | `string`                                   | The URL of the repository
| repos:branches            | no       | master         | `string`                                   | Tracking branches of the repository
| repos:source_root         | no       |                | `string`                                   | Source root to apply on the repo.
| repos:expand_keys         | no       |                | true, false                                | Enable/disable file content evaluation.
| repos:skip_branch_name    | no       | false          | true, false                                | Enable/disable branch name pruning.
| repos:skip_repo_name      | no       | false          | true, false                                | Enable/disable repository name pruning.
| repos:mount_point         | no       |                | `string`                                   | Sets the prefix which should be used for the path in the Consul KV Store
| repos:credentials:username     | no       |                | `string`                          | Username for the Basic Auth
| repos:credentials:password | no       |                | `string`                          | Password/token for the Basic Auth
| repos:credentials:private_key:pk_key | no       |                | `string`      | Path to the priv key used for the authentication
| repos:credentials:private_key:pk_username     | no       |     git         | `string`              | Username used with the ssh authentication
| repos:credentials:private_key:pk_password | no       |                | `string`                  | Password used with the ssh authentication
| repos:hooks:type          | no       | polling        |  polling, webhook | Type of hook to use to fetch changes on the repository. See [below](#webhooks).
| repos:hooks:interval      | no       | 60             | `int`                                      | Interval, in seconds, to poll if polling is enabled
| repos:hooks:url           | no       | ??             | `string`                                   | ???
| consul:address            | no       | 127.0.0.1:8500 | `string`                                   | Consul address to connect to. It can be either the IP or FQDN with port included
| consul:ssl                | no       | false          | true, false                                | Whether to use HTTPS to communicate with Consul
| consul:ssl_verify         | no       | false          | true, false                                | Whether to verify certificates when connecting via SSL
| consul:token              | no       |                | `string`                                   | Consul API Token

### Webhooks

Webhooks will be served from a single port, and different repositories will be given different endpoints according to their name

Available endpoints:

* `<webhook:address>:<webhook:port>/{repository}/github`
* `<webhook:address>:<webhook:port>/{repository}/stash`
* `<webhook:address>:<webhook:port>/{repository}/bitbucket`
* `<webhook:address>:<webhook:port>/{repository}/gitlab`


### Options

#### source_root (default: undefined)

The "source_root" instructs the app to navigate to the specified directory in the git repo making the value of source_root is trimed from the KV Store key. By default the entire repo is evaluated.

When you configure the source_root with `/top_level/lower_level` the file `/top_level/lower_level/foo/web.json` will be mapped to the KV store as `/foo/web.json`

#### mount_point (default: undefined)

The "mount_point" option sets the prefix for the path in the Consul KV Store under which the keys should be added.

#### expand_keys (default: undefined)

The "expand_keys" instructs the app to evaluate known types of files. The content of the file is evaluated to key-value pair and pushed to the Consul KV store.

##### Supported formats
* Text file - the file content is pushed to the KV store as it is.
* Yaml file - the file is evaluated into key-value pair. i.e `configuration.yml`
```
---
services:
  apache:
    port: 80
  ssh:
    port: 22
```

will be evaluated to the following keys:
* `/configuration/services/apache/port`
* `/configuration/services/ssh/port`

#### skip_branch_name (default: false)

The "skip_branch_name" instructs the app to prune the branch name. If set to true the branch name is pruned from the KV store key.

#### skip_repo_name (default: false)

The "skip_repo_name" instructs the app to prune the repository name. If set to true the repository name is pruned from the KV store key.

#### credentials

The "credentials" option provides the possibility to pass the credentials to authenticate to private git repositories.

Sample config with basic auth (login:password/token)
```
{
  "repos": [
    {
      "name": "example",
      "url": "http://github.com/DummyOrg/ExampleRepo.git",
      "credentials: {
            "username": "foo",
            "password": "bar"
      }
    }
  ]
}
```
Sample config with ssh auth
```
{
  "repos": [
    {
      "name": "example",
      "url": "ssh://github.com/DummyOrg/ExampleRepo.git",
      "credentials: {
            "private_key": {
                  "pk_key": "/path/to/priv_key",
                  "pk_username": "foo",
                  "pk_password": "bar"
            }
      }
    }
  ]
}
```

## Developing

See [CONTRIBUTING.md](.github/CONTRIBUTING.md) for details.

### Dependencies
* Go 1.15+

### Compiling From Source
```
$ make build
```

### End To End Testing

The end to end test can be run by simply running `make test-e2e`.
To run the test, you need to have the [consul binary](https://releases.hashicorp.com/consul/) available in your path.
It simulates a create and update of consul KV pairs and confirms every operation is successful.
The test data is stored within this repo so developers do not have to setup an external repo to test.

The tests can manually be run by starting Consul in dev mode, and then manually running `git2consul` with one of the config files provided.
For example:
```
$ consul agent -dev
$ # In a separate terminal
$ ./git2consul -config pkg/e2e/data/create-config.json -once -debug
```

### Releases
This project is using [goreleaser](https://goreleaser.com). GitHub release creation is automated using Travis
CI. New releases are automatically created when new tags are pushed to the repo.
```
$ TAG=v0.0.2 make tag
```

How to manually create a release without relying on Travis CI.
```
$ TAG=v0.0.2 make tag
$ GITHUB_TOKEN=xxx make clean release
```

## License

See [LICENSE](LICENSE) for details.

## Acknowledgement

See [ACKNOWLEDGEMENT.md](ACKNOWLEDGEMENT.md) for details.

## Code of Conduct

See [CODE_OF_CONDUCT.md](.github/CODE_OF_CONDUCT.md) for details.
