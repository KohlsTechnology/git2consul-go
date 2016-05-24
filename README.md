# go-git2consul

## Default configuration

git2consul will attempt to use sane defaults for configuration. However, since git2consul needs to know which repository to pull from, minimal configuration is necessary.

| Configuration        | Required | Default Value  | Avail. Values                              | Description
|----------------------|----------|----------------|--------------------------------------------| -----------
| local_store          | no       | `os.TempDir()` | `string`                                   | Local cache for git2consul to store its tracked repositories
| webhook_port         | no       | 9000           | `int`                                      |  Webhook port that that git2consul will be using
| repos:name           | yes      |                | `string`                                   | Name of the repository, this will match against the path of the webhook, if any is present
| repos:url            | yes      |                | `string`                                   | The URL of the repository
| repos:branches       | no       | master         | `string`                                   | Tracking branches of the repository
| repos:hooks:type     | no       | polling        |  polling, github, stash, bitbucket, gitlab | Type of hook to use to fetch changes on the repository
| repos:hooks:interval | no       | 60             | `int`                                      | Interval, in seconds, to poll if polling is enabled
| consul:address       | no       | 127.0.0.1:8500 | `string`                                   | Consul address to connect to. It can be either the IP or FQDN with port included
| consul:ssl           | no       | false          | true, false                                | Whether to use HTTPS to communicate with Consul
| consul:ssl_verify    | no       | false          | true, false                                | Whether to verify certificates when connecting via SSL
| consul:token         | no       |                | `string`                                   | Consul API Token

## Development dependencies:
* Go 1.6
* libgit2 v0.24.0
* [glide](https://github.com/Masterminds/glide)
