# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

> **Version note:** this release contains behavior changes (see _Changed_). They
> stem from fixing latent bugs, but they alter observable output. Decide between
> a minor bump (`v1.5.0`) with the notes below, or a major bump (`v2.0.0`,
> requiring the `/v2` module path) for strict SemVer. The notes below assume
> `v1.5.0`.

### Added
- `Validator.FailFast` — opt-in flag to stop at the first failing field and
  return immediately (similar to Joi's `abortEarly` / Laravel's
  `stopOnFirstFailure`). Defaults to `false` (collect all errors), preserving
  existing behavior.
- Native fuzz targets (`fuzz_test.go`) for the string validators, converters,
  date parser, the phone library, and the full `ValidateStruct` reflection path.
- This `CHANGELOG.md` and a `SECURITY.md`.

### Changed
- **Untagged nested structs and pointers-to-struct are now validated
  recursively.** Previously a nested-struct field needed its own `valid` tag for
  its inner fields to be checked; they were otherwise silently skipped.
  `time.Time` and `decimal.Decimal` are still treated as scalar types (not
  recursed into). **This can surface new validation errors for previously
  unvalidated nested data.**
- **Duplicate errors removed.** A struct that previously reported each field's
  error twice (e.g. 3 fields → 5 errors) now reports one error per field.
- `FieldError.FuncError` is now tagged `json:"-"` (it previously serialized to an
  empty `{}`); use `Message` for client-facing output and `Unwrap()` for the
  underlying error.
- Promoted fields from embedded (anonymous) structs are now validated as if
  declared on the parent.

### Fixed
- **Data race + unbounded memory growth in `requiredIf`** — the rule appended a
  message parameter to a shared cached tag on every failing validation;
  parameters accumulated across calls and concurrent validation raced on the
  slice. The value is now built per call.
- **`O(2^depth)` nested validation** — the field cache duplicated every field
  after the first, so each recursing field was processed twice per level. A
  20-level-deep struct performed ~2M validations (~4s). Now linear (a 500-deep
  chain validates in ~6ms).
- **Cyclic-reference hang (DoS)** — self-referential graphs recursed forever; a
  depth guard (`maxValidationDepth`) now terminates them with an error.
- **Embedded struct fields skipped** — promoted fields were never collected and
  were mis-resolved via the first index element only. Now resolved via the full
  index path (`FieldByIndexErr`); nil embedded pointers are skipped without
  panicking.
- **i18n message dropping** — `Translator.Trans` stopped after the first field
  that used a custom message, leaving the rest untranslated.
- **`sync.Map` data races** in the custom-type unwrapper (the cache was cleared
  by reassigning the map value instead of clearing it in place).
- **Date parsing divergence** — field values and rule parameters used different
  format lists, so the same string could parse to different dates. Unified.
- **`omitempty`/`nullable` substring match** — these were detected with
  `strings.Contains` on the raw tag and could false-positive on a parameter
  value (e.g. `contains=omitempty`). They are now matched as exact options.
- **`map[string]interface{}` entries** holding structs were silently skipped (the
  deref check tested the map's kind, not the entry's).

### Performance
- **Allocation-free success path.** Validation that passes now allocates nothing
  (was 2 allocations/field) and is ~2.5× faster, by reusing the cached field
  name instead of rebuilding it on every field.
- `errorResponsePool` buffer capacity is bounded to avoid pinning oversized
  backing arrays from one-off large error sets.

### Security
- **ReDoS is structurally impossible**: only Go's standard `regexp` (RE2,
  linear-time) engine is used — no backtracking regex library is present.
- Fuzzing (~6.6M executions across all targets) found no panics, hangs, or
  crashers on adversarial input.

### Internal
- Migrated the golangci-lint config to v2 (CI ran `latest` against a v1-format
  config, so the lint job was effectively failing to parse it).
- Split the static MIME table into `mimes.go`, removed dead code, and modernized
  `interface{}` to `any` across the codebase.
- Test coverage raised from ~78% to ~88%; all fixes have regression tests.

## [1.4.0] - earlier

See the Git history for releases prior to this changelog.
