# Gitlab CI Component Linter

![GitHub Release](https://img.shields.io/github/v/release/sapher/gitlab-ci-component-linter)

Linter for Gitlab CI components. Validate that project follow the [guidelines](https://docs.gitlab.com/ee/ci/components/#directory-structure) defined by Giltab for CI components.

> This loot name is not that great, if you have a better idea, please open an issue.

## Install

You can either install it using `go install` :

```bash
go install github.com/sapher/gitlab-ci-component-linter
```

Or download the binary from the [release page](https://github.com/sapher/gitlab-ci-component-linter/releases).

## Usage

### CLI

```
Validate Gitlab CI Component files against set of rules

Usage:
  gitlab-ci-component-linter [flags]

Flags:
  -h, --help             help for gitlab-ci-component-linter
      --output string    Output format, one of: json, yaml, junit (default "json")
      --soft-fail        Wether to fail or not on error
      --workdir string   Working directory (default ".")
```

### Git hook

You can use this linter as git hook with `pre-commit` to validate your changes before commiting them.

```yaml
repos:
  - repo: https://github.com/sapher/gitlab-ci-component-linter
    rev: v1.1.0
    hooks:
      - id: gitlab-ci-component-linter
```

### Gitlab CI

You can use this linter as a job in your Gitlab CI pipeline.

```yaml
stages:
  - lint

lint:
  stage: lint
  image: alpine:latest
  script:
    # Download binary
    - apk add --no-cache wget
    - wget -q -O gitlab-ci-component-linter.zip  https://github.com/sapher/gitlab-ci-component-linter/releases/download/v1.0.0/gitlab-ci-component-linter_1.0.0_linux_amd64.zip
    - unzip gitlab-ci-component-linter.zip
    - mv gitlab-ci-component-linter_v1.0.0 /usr/local/bin/gitlab-ci-component-linter
    # Run linter
    - gitlab-ci-component-linter $CI_PROJECT_DIR --output junit | tee junit.xml
  artifacts:
    reports:
      junit: junit.xml
```
