# deslop

A standalone Go linter that enforces opinionated code style and guides LLMs toward preferred libraries and patterns.

## Install

```sh
go install github.com/matthewmueller/deslop@latest
```

## Usage

```sh
deslop ./...
```

## Analyzers

### Template-based

These analyzers flag discouraged patterns and include usage templates in the diagnostic message, helping LLMs produce correct code on the first try.

| Analyzer | Flags | Suggests |
|----------|-------|----------|
| `envcheck` | `os.Getenv`, `os.LookupEnv` | `internal/env` with `caarlos0/env/v11` |
| `clicheck` | `os.Args`, `flag`, cobra, urfave | `internal/cli` with `livebud/cli` |
| `muxcheck` | `http.NewServeMux`, gorilla, chi, gin, echo | `livebud/mux` |
| `nolog` | stdlib `log` package | `log/slog` + `matthewmueller/logs` |

### Test style

| Analyzer | Flags |
|----------|-------|
| `notabletest` | `[]struct{}` table-driven test patterns |
| `ischeck` | `t.Fatal`/`t.Error` for assertions, testify imports |
| `testnaming` | `Test_Something` underscored test names |

### Error handling

| Analyzer | Flags |
|----------|-------|
| `errpath` | `if err == nil` (escape hatch: `nil == err`), unnecessary else after return |
| `errwrap` | `fmt.Errorf` without `%w` |
| `lowercaseerror` | `fmt.Errorf("Uppercase...")` — error strings should start lowercase |

### Code style

| Analyzer | Flags |
|----------|-------|
| `noelse` | Any `else` block — use early returns instead |
| `noglobals` | Exported package-level `var` (allows `Err*`, `_`, unexported) |
| `nogetter` | Methods named `GetX` — just use the noun |
| `noinit` | `init()` functions |
| `nonakedreturn` | Bare `return` in named-return functions |
| `nostutter` | Type names that repeat the package name (allows exact match like `user.User`) |
| `contextfirst` | `context.Context` not as first parameter |
| `nocommentedcode` | Commented-out code (function calls, `:=`, keywords) |
| `nofmt` | `fmt.Println` |

### Complexity

| Analyzer | Default | Flag |
|----------|---------|------|
| `nesting` | Max depth 3 | `-nesting.max` |
| `params` | Max 4 parameters | `-params.max` |
| `funcnaming` | Max 30 character function name | `-funcnaming.max` |
| `linelength` | Max 120 characters per line (skips strings) | `-linelength.max` |
| `funclength` | Max 500 lines per function | `-funclength.max` |
