# Gitlab CI Component Linter

Linter for Gitlab CI components. Validate that project follow the [guidelines](https://docs.gitlab.com/ee/ci/components/#directory-structure) defined by Giltab for CI components.

## Install

```bash
go install github.com/sapher/gitlab-ci-component-linter
```

## Usage

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
