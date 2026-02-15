# AGENTS.md

Guidance for coding agents working in `github.com/dector/kdly`.

This repository is a Go implementation of a KDL v2 parser and serializer.

## Scope and Priority

1. Follow this file for repo-specific workflow and style.
2. Follow standard Go conventions when this file is silent.

## Project Layout

- `main.go`: public package `kdly` API; wrapper/re-export layer.
- `pkg/parser/`: parser implementation (state machine, KDL model types, parser tests).
- `pkg/serializer/`: serializer implementation and tests.
- `ror.kdl`: `ror` task definitions.
- `samples/parse/`: separate sample module.
- `tasks/`: planning/notes artifacts; not runtime code.

## Toolchain

- Go version in module: `go 1.25.4`.
- Main test assertion library: `github.com/stretchr/testify/assert`.
- Task runner available: `ror`.

## Build, Test, and Lint Commands

### Canonical Commands

- Run all configured tasks help: `ror --help`
- Run repository test task: `ror test`
- Direct full test run: `go test ./...`
- Build all packages: `go build ./...`

### Running a Single Test (Important)

- Single package: `go test ./pkg/parser`
- Single exact test: `go test ./pkg/parser -run '^TestSingleSimpleNode$'`
- Single exact test in serializer: `go test ./pkg/serializer -run '^TestSerializeServerConfig$'`
- Single subtest: `go test ./pkg/serializer -run '^TestSerializerEdgeCases$/^empty document$'`
- Force rerun (disable test cache): `go test ./pkg/parser -run '^TestSingleSimpleNode$' -count=1`

### Useful Narrowing Patterns

- Run all tests matching prefix: `go test ./pkg/parser -run '^TestTypeAnnotations'`
- Run one package with verbose names: `go test -v ./pkg/parser -run '^TestNodeTypeAnnotation_Basic$'`
- Run tests across all packages by name: `go test ./... -run '^TestSerializeServerConfig$'`

### Lint/Format (No Dedicated Linter Config Checked In)

- Format all code: `gofmt -w .`
- Vet all packages: `go vet ./...`
- Optional import normalization (if installed): `goimports -w .`

## Test Suite Reality (Important for Agents)

- `go test ./...` currently fails in `pkg/parser` due many spec/TODO tests in `parser_todo_test.go`.
- `pkg/serializer` tests currently pass.
- Do not assume red `go test ./...` means your change is wrong; inspect failing test names.
- For incremental work, run targeted tests with `-run` for affected behavior.
- When reporting results, clearly separate pre-existing failures from new failures.

## Coding Style Guidelines

### Formatting and File Hygiene

- Always produce `gofmt`-formatted Go code.
- Keep files ASCII unless a non-ASCII literal is required by test/spec coverage.
- Prefer small, focused diffs over broad refactors.
- Preserve existing public API names unless explicitly asked to change them.

### Imports

- Use standard Go import grouping:
  - standard library
  - third-party modules
  - local module imports
- Keep import blocks sorted as `gofmt`/`goimports` expects.
- Avoid unused imports; keep package-level imports minimal.

### Types and Data Structures

- Keep model types in parser domain (`Document`, `Node`, `Value`, `Property`) consistent.
- Use `ValueType` constants (`ValueTypeString`, `ValueTypeNumber`, etc.) instead of raw strings.
- Prefer explicit struct literals in tests for readability of expected ASTs.
- For mutable parser state, keep pointer receivers (`func (p *Parser) ...`).

### Naming

- Exported identifiers: `PascalCase` with Go doc comments when exported.
- Internal helpers: `camelCase` and concise, behavior-oriented names.
- Test names:
  - Keep descriptive `TestXxx` names for hand-written tests.
  - Preserve underscore-heavy naming in spec-derived tests (`Test_...`) for compatibility.

### Control Flow and Parser Internals

- Parser implementation uses a hand-written state machine; keep transitions explicit.
- Keep parser position tracking (`pos`, `line`, `col`) correct when adding branches.
- When parsing alternatives, prefer guard checks before consuming input.
- Maintain existing semantics for comments, slashdash, type annotations, and raw strings.

### Error Handling

- In parser internals, follow current pattern:
  - call `p.panicAt(...)` for parse errors
  - recover in `Parse` and return `error`
- Include actionable parse context (`line`, `col`, specific expectation) in new errors.
- Outside parser core, return normal Go `error` values; do not introduce panic-based flow.
- Keep error text stable when tests assert specific substrings.

### Comments and Documentation

- Keep comments for non-obvious behavior, invariants, and spec edge cases.
- Avoid redundant comments that restate trivial code.
- Add/update doc comments when changing exported API behavior.

## Testing Conventions

- Use `assert.NoError`/`assert.Equal` patterns already used in tests.
- For expected errors, use both:
  - `assert.Error(t, err)`
  - `assert.Contains(t, err.Error(), "...")` when message semantics matter.
- Keep test structure clear: `input`, parse/act, `want`, assertions.
- Prefer exact `-run '^TestName$'` matching when validating one test.

## Package-Specific Notes

### `pkg/parser`

- Core behavior lives in parser state and helper functions (`parse*`, `skip*`, `looksLikeNumber`).
- Be careful when changing identifier or number logic; many tests depend on edge cases.
- `WithAllowDuplicateProperties(bool)` and `WithDisableTypeAnnotations(bool)` are parser-level toggles.

### `main.go` (public package `kdly`)

- This is a thin wrapper over `pkg/parser` and `pkg/serializer`.
- Wrapper methods intentionally expose simpler toggles:
  - `WithAllowDuplicateProperties()` (enables duplicates)
  - `WithNoTypeAnnotations()` (disables type annotations)
- Keep wrapper API ergonomic and backward compatible.

### `pkg/serializer`

- Keep output deterministic and round-trip friendly with parser.
- Preserve current indentation style (2 spaces per nesting level).
- Maintain value formatting decisions (`formatIdentifier`, raw string preference, escaping).

## Suggested Agent Workflow

1. Inspect touched package and nearest tests first.
2. Implement minimal change that matches existing architecture.
3. Run narrow tests (`go test ./pkg/<pkg> -run '^TestName$' -count=1`).
4. Run broader package tests for impacted package.
5. If feasible, run `ror test` and report known baseline failures separately.

## When Unsure

- Prefer consistency with adjacent code over introducing a new pattern.
- If behavior is spec-sensitive, encode it in tests before or along with changes.
- Document any intentional deviation from KDL v2 expectations in tests/comments.
