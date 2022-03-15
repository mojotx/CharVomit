
# CharVomit
![Go Coverage](https://raw.githubusercontent.com/mojotx/CharVomit/master/coverage_badge.png)

Generate random passwords using Go's
[crypto/rand](https://golang.org/pkg/crypto/rand/) functions.

The passwords should look like the cat puked up random characters,
hence the name.

## Usage

```shell
Usage: ./CharVomit [ length ]

If a password length is not specified, 32 is used.

Other optional flags are:
  -d	use numeric digits
  -h	show help and exit
  -l	use lower-case letters
  -s	use symbols: !#%+:=?@
  -u	use upper-case letters
  -w	use weak characters (2-9, A-Z, a-z)

Note that optional flags must precede the password length.

For example, a 8-character password of all capital letters:
./CharVomit -u 8
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
