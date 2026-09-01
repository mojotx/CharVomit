# gfi-finder

## Environment

- Never create git commits or push on the user's behalf. When work is ready, suggest a commit message instead and let the user commit/push themselves.

## Build and Test

After ANY Go code change, run the same checks CI does (`.github/workflows/ci.yml`):

```sh
go mod tidy && git diff --exit-code -- go.mod go.sum
go build ./...
go test -race ./...
golangci-lint run --timeout=5m --allow-parallel-runners --max-same-issues 0 --max-issues-per-linter 0 ./...
```

Use `./...` exactly as written — don't let it get mangled into a path fragment (e.g. `go test -race mojotx.`).

CI also runs `go vet ./...` as its own step, but that's redundant to run locally — `golangci-lint` already runs the `govet` linter (plus staticcheck and many others), so a clean `golangci-lint` run implies `go vet` is clean too.

## Conventions

- Tests use `github.com/stretchr/testify` (`assert`/`require`), table-driven where the cases share shape. Use named struct fields in table-test case literals, not positional, once a case has more than a couple of fields.
- Mock at the HTTP layer, not by hand-rolling fake clients: use an `http.RoundTripper` stub against the real `go-gh` `api.RESTClient`/`api.GraphQLClient`, matching `cli/cli`'s own testing conventions.
- No inline multi-line query strings in Go source (e.g. GraphQL). Put them in their own file (e.g. `timeline.graphql`) and pull them in via `//go:embed`.
- Package doc comments go immediately above the `package` clause, no blank line in between.
- Don't add an `if err != nil` check around a write that can never fail in practice (e.g. writes into a buffered `tabwriter.Writer`) — ignore it explicitly with `_, _ = fmt.Fprintf(...)` and only check the error at the real I/O boundary (e.g. `Flush()`).
- Always check errors, or explicitly ignore them with `_, _ = ...` when safe. Never silently swallow errors without a good reason.
- When wrapping errors with `fmt.Errorf`, put `%w` at the end (`"doing X: %w"`), not the start — keeps error chains reading newest-to-oldest. Don't add an annotation that only restates "it failed" with no new information (e.g. `fmt.Errorf("failed: %w", err)` — just return `err`).
- Keep interfaces small and defined by the consumer (see `search.Searcher`, `linkcheck.Checker`, `search.RESTGetter`): only the methods actually called, not a mirror of the whole client. Constructors return concrete types (`*Client`), not interfaces.
- If a function's parameter list grows past a handful of related flags/options, prefer an options struct (see `cmd.Options`) over more positional parameters.
- Use `github.com/spf13/cobra` for Go CLI applications
- prefer `github.com/rs/zerolog` for logging
- Suggested git commit messages should follow best practices: 50 characters for the first line, a blank line, and then 72 characters max for following lines
- Use the `gh` GitHub CLI tool if you need to view PR status, comments, issues, etc.
