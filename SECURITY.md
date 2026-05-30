# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.4.x and later | :white_check_mark: |
| < 1.4   | :x:                |

Only the latest minor release receives security fixes. Please upgrade before
reporting issues against older versions.

## Reporting a Vulnerability

Please report security issues **privately** — do not open a public issue for an
undisclosed vulnerability.

- Preferred: open a private advisory via GitHub →
  **Security** → **Report a vulnerability**
  (<https://github.com/syssam/go-validator/security/advisories/new>).
- Include a minimal reproduction, the affected version, and the impact.

You can expect an initial acknowledgement within a few days. Once a fix is
available, a patched release and an advisory will be published.

## Security Posture

This library validates untrusted input, so the following properties are
maintained deliberately:

- **No catastrophic-backtracking (ReDoS).** All pattern matching uses Go's
  standard `regexp` package, which is backed by RE2 and guarantees linear-time
  matching. No backtracking regex engine (e.g. `regexp2`, PCRE) is used.
- **Bounded recursion.** Struct validation is depth-guarded
  (`maxValidationDepth`), so cyclic or pathologically deep object graphs return
  an error instead of exhausting the stack or hanging.
- **No panics on adversarial input.** The input-facing validators, converters,
  date parser, and the full reflection path are exercised by fuzz targets
  (`fuzz_test.go`). Re-run with, e.g.:

  ```sh
  go test -fuzz=FuzzValidateStruct -fuzztime=60s
  ```

### Notes for integrators

- Validation cost scales with the number of fields and collection elements. For
  endpoints that accept untrusted, deeply nested or large payloads, enforce
  request size/shape limits at the transport layer in addition to validation.
- The package-level configuration maps (`MessageMap`, `RuleMap`, etc.) are not
  synchronized; configure them once at startup, before concurrent validation.
