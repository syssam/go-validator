# Contributing

Thanks for your interest in improving `go-validator`. This guide covers the
local workflow and the conventions CI enforces.

## Prerequisites

- Go **1.23** or newer (see `go.mod`).
- [`golangci-lint`](https://golangci-lint.run/) **v2** for local linting.

## Getting started

```sh
git clone https://github.com/syssam/go-validator
cd go-validator
go mod download
go test ./...
```

## Before you open a PR

CI runs the checks below; run them locally first so the PR goes green on the
first try.

```sh
gofmt -l .            # must print nothing
go vet ./...
golangci-lint run ./...
go test -race ./...    # the suite must pass under the race detector
go test -cover ./...
```

A one-liner:

```sh
gofmt -l . && go vet ./... && golangci-lint run ./... && go test -race ./...
```

## Conventions

- **Formatting:** `gofmt` is mandatory; CI fails on unformatted files.
- **No panics:** validators must return errors, never panic — even on
  adversarial input (see `fuzz_test.go`).
- **Thread-safety:** validation must be safe to run concurrently. Do not mutate
  cached metadata (`field`/`ValidTag`) during validation. Package-level config
  maps (`MessageMap`, `RuleMap`, …) are configured once at startup, not at
  runtime — use `CustomTypeRuleMap` for runtime-safe registration.
- **Tests are required.** Every bug fix needs a regression test that fails
  against the old behavior; every new rule needs unit tests. Keep the
  allocation-free success path allocation-free (`go test -bench=. -benchmem`).

## Adding a validation rule

1. **No-parameter rule** (e.g. `distinct`): add to `RuleMap` in `types.go`.
2. **Parameter rule** (e.g. `min=value`): add to `ParamRuleMap` in `types.go`.
3. **String pattern rule** (e.g. `email`): add the pattern to `patterns.go`, the
   function to `validator_string.go`, and register it in `StringRulesMap`.
4. Add the default message to `MessageMap` in `message.go` (and the `lang/*`
   files for translations).
5. Add tests in the matching `*_test.go` file; add a benchmark if the rule does
   non-trivial work.

## Fuzzing

Input-facing changes should be fuzzed:

```sh
go test -fuzz=FuzzValidateStruct -fuzztime=60s
```

## Reporting bugs & security issues

- **Bugs:** open a GitHub issue with a minimal reproduction.
- **Security vulnerabilities:** do **not** open a public issue — follow
  [`SECURITY.md`](SECURITY.md).
