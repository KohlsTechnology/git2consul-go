# go-git2consul

go-git2consul is a port of [git2consul](https://github.com/Cimpress-MCP/git2consul), which had great success and adoption. go-git2consul takes on the same basic principles as its predecessor, and attempts to improve upon some of its feature sets as well as add new ones. There are a few advantages to go-git2consul, including the use of the official Consul API, which is written in Go, and removing runtime dependencies such as node and git.

## Default configuration

git2consul will attempt to use sane defaults for configuration. However, since git2consul needs to know which repository to pull from, minimal configuration is necessary.

| Configuration        | Required | Default Value  | Available Values                           | Description
|----------------------|----------|----------------|--------------------------------------------| -----------
| local_store          | no       | `os.TempDir()` | `string`                                   | Local cache for git2consul to store its tracked repositories
| webhook:address      | no       |                | `string`                                   | Webhook listener address that that git2consul will be using
| webhook:port         | no       | 9000           | `int`                                      | Webhook listener port that that git2consul will be using
| repos:name           | yes      |                | `string`                                   | Name of the repository. This will match the webhook path, if any are enabled
| repos:url            | yes      |                | `string`                                   | The URL of the repository
| repos:branches       | no       | master         | `string`                                   | Tracking branches of the repository
| repos:hooks:type     | no       | polling        |  polling, github, stash, bitbucket, gitlab | Type of hook to use to fetch changes on the repository
| repos:hooks:interval | no       | 60             | `int`                                      | Interval, in seconds, to poll if polling is enabled
| consul:address       | no       | 127.0.0.1:8500 | `string`                                   | Consul address to connect to. It can be either the IP or FQDN with port included
| consul:ssl           | no       | false          | true, false                                | Whether to use HTTPS to communicate with Consul
| consul:ssl_verify    | no       | false          | true, false                                | Whether to verify certificates when connecting via SSL
| consul:token         | no       |                | `string`                                   | Consul API Token

## Webhooks

Webhooks will be served from a single port, and different repositories will be given different endpoints according to their name

Available endpoints:

* <webhook:address>:<webhook:port>/{repository}/github
* <webhook:address>:<webhook:port>/{repository}/stash
* <webhook:address>:<webhook:port>/{repository}/bitbucket
* <webhook:address>:<webhook:port>/{repository}/gitlab

## Future feature additions
* File format backend
* Run -once flag
* Support for source_root and mountpoint
* Support for tags as branches
* Support for Consul HTTP Basic Auth
* Logger support for other handlers other than text
* Auth support for webhooks banckends


## Development dependencies:
* Go 1.6
* libgit2 v0.24.0
* [glide](https://github.com/Masterminds/glide)
