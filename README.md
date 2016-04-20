# go-git2consul

# Defaults

git2consul will attempt to deduce sane defaults for configuration. However, as git2consul needs to know which repository to pull from, minimal configuration is necessary.

## Repository-level configuration

| configuration  | default |
|----------------|---------|
| branches       | master  |
| hooks:type     | polling |
| hooks:interval | 60      |
