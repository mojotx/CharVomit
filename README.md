
# CharVomit

[![CI](https://github.com/mojotx/CharVomit/actions/workflows/ci.yml/badge.svg)](https://github.com/mojotx/CharVomit/actions/workflows/ci.yml)
[![CodeQL](https://github.com/mojotx/CharVomit/actions/workflows/codeql.yml/badge.svg)](https://github.com/mojotx/CharVomit/actions/workflows/codeql.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/mojotx/CharVomit/v2.svg)](https://pkg.go.dev/github.com/mojotx/CharVomit/v2)

Generate random passwords using Go's
[crypto/rand](https://golang.org/pkg/crypto/rand/) functions.

The passwords should look like the cat puked up random characters,
hence the name.

## Usage

```shell
Usage: CharVomit [ length ] [flags]

If a password length is not specified, 32 is used. With no character flags,
the default is equivalent to `CharVomit 32 -w -s`: weak, non-ambiguous
characters plus symbols. The `-d` flag is redundant with `-w` because the weak
pool already contains the non-ambiguous digits.

Other optional flags are:
    -d, --digits           use numeric digits
    -h, --help             show help and exit
    -l, --lowercase        use lower-case letters
    -o, --output-file      write the generated password to a file instead of stdout
    -s, --symbols          use symbols: !#%+:=?@
    -u, --uppercase        use upper-case letters
    -v, --version          show version
    -w, --weak             use weak characters (2-9, A-N, P-Z, a-k, m-z)
    -x, --exclude string   excluded characters (will be removed)

By default, CharVomit writes the generated password to standard output. To
avoid exposing the password in terminal scrollback or shell history, prefer
--output-file / -o and write it to a secure local file.

For example, a 8-character password of all capital letters:
CharVomit -u 8

Write to a file instead of the terminal:
CharVomit --output-file /path/to/password.txt 20

Also note that certain characters that are confusing are ignored by default,
such as '0', 'O', '1', and 'l'. You can still get those characters, if you wish,
by using the -u, -l, and -d flags. The default is equivalent to `CharVomit 32 -w -s`.

```

## Breaking changes in v2.0.0

The module path is now `github.com/mojotx/CharVomit/v2`. If you imported
`pkg/arg` or `pkg/CharVomit` as a library, update your import paths
accordingly. The `pkg/arg.ParseArgs`, `ParseConfig`, `Parse`, `Usage`,
`RegisterFlags` functions and package-global `Config` variable have been
removed, as has `RemoveIndex`; none were reachable from the CLI, which has
used `github.com/spf13/cobra` for flag parsing since v2.0.0. Build an
`arg.ConfigType` directly instead — see "Library usage" below.

## Library usage

If you want to use the package from Go code instead of the CLI, prefer building a local `arg.ConfigType` value rather than relying on package-global state.

```go
config := arg.ConfigType{WeakChars: true, Symbols: true}

cv := CharVomit.NewCharVomit("")
if err := cv.SetAcceptableChars(config); err != nil {
    panic(err)
}
```

## Examples

### 32-character password, without ambiguous characters

```shell
$ CharVomit
Va9nBzgtW:Xt@28pcXW+6zpjb@DuyqJ3
```

### 20-character password, with weak (non-ambiguous) characters, no symbols

```shell
$ CharVomit -w 20
qm995CZrA7pRC4SgfDrJ
```

### 20-character password, with all upper- and lower-case letters and digits, as well as the symbols '!' and '@'

```shell
$ CharVomit -l -u -s -x '#%+:=?' 20
Xl!bXDnZC@srbxBDNzdj
```

## To-Do

* Implement functionality to specify number of duplicate characters
* Improve documentation

## Installation

If you have a Go compiler installed, you can use this command:

```shell
go install -v github.com/mojotx/CharVomit/v2/cmd/CharVomit@latest
```

Alternatively, you can download the latest [release](https://github.com/mojotx/CharVomit/releases).
