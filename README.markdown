
# CharVomit
![Go Coverage](https://raw.githubusercontent.com/mojotx/CharVomit/master/coverage_badge.png)

Generate random passwords using Go's
[crypto/rand](https://golang.org/pkg/crypto/rand/) functions.

The passwords should look like the cat puked up random characters,
hence the name.

## Usage

```shell
Usage: CharVomit [ length ]

If a password length is not specified, 32 is used.

Other optional flags are:
  -d	use numeric digits
  -h	show help and exit
  -l	use lower-case letters
  -s	use symbols: !#%+:=?@
  -u	use upper-case letters
  -v	show version
  -w	use weak characters (2-9, A-N, P-Z, a-k, m-z)
  -x string
    	excluded characters (will be removed)

Note that optional flags must precede the password length.

For example, a 8-character password of all capital letters:
CharVomit -u 8

Also note that certain characters that are confusing are ignored by default,
such as '0', 'O', '1', and 'l'. You can still get those characters, if you wish,
by using the -u, -l, and -d flags. The default is equivalent to -w -s.

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

### 16-character password, with all upper- and lower-case letters and digits, as well as the symbols '!' and '@'

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
go install -v github.com/mojotx/CharVomit/cmd/CharVomit@latest
```

Alternatively, you can download the latest [release](https://github.com/mojotx/CharVomit/releases).
