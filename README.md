# flaw

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][action-img]][action-url]
[![Coverage][codecov-img]][codecov-url]
[![Go Report Card][report-img]][report-url]

A flaw is Golang is a replacement of errors package

## Installation

Make sure you have a working Go environment. Go version 1.13.x is supported.

[See the install instructions for Go](http://golang.org/doc/install.html).

To install flaw, simply run:

```
$ go get github.com/phogolabs/flaw
```

## Getting Started

Wrapping an error:

```golang
err := flaw.Wrap(sql.ErrNoRows).
	WithMessage("user does not exist").
	WithCode(404)
```

If you print the error with `fmt.Println(err)` you will receive the following
output:

```bash
code: 404 message: user does not exist cause: sql: no rows in result set
```

If you want more detailed information, you can use `fmt.Printf("%+v", err)`
formatting:

```bash
    code: 404
 message: user does not exist
   cause: sql: no rows in result set
   stack:
 --- /Users/ralch/go/src/github.com/phogolabs/flaw/cmd/main.go:19 (main)
 --- /usr/local/Cellar/go/1.13.1/libexec/src/runtime/proc.go:203 (main)
 --- /usr/local/Cellar/go/1.13.1/libexec/src/runtime/asm_amd64.s:1357 (goexit)
 ```

Collecting multiple errors:

```golang
errs := flaw.ErrorCollector{}
errs = append(errs, fmt.Errorf("insufficient funds"))
errs = append(errs, fmt.Errorf("maximum allowance reached"))
```

Then you can print the error `fmt.Println(errs)`:

```bash
[insufficient funds, maximum allowance reached]
```

You can print the result with better formatting `fmt.Printf("%+v", errs)`:

```bash
 --- insufficient funds
 --- maximum allowance reached
```

## Contributing

We are open for any contributions. Just fork the
[project](https://github.com/phogolabs/stride).

[report-img]: https://goreportcard.com/badge/github.com/phogolabs/flaw
[report-url]: https://goreportcard.com/report/github.com/phogolabs/flaw
[codecov-url]: https://codecov.io/gh/phogolabs/flaw
[codecov-img]: https://codecov.io/gh/phogolabs/flaw/branch/master/graph/badge.svg
[action-img]: https://github.com/phogolabs/flaw/workflows/pipeline/badge.svg
[action-url]: https://github.com/phogolabs/flaw/actions
[godoc-url]: https://godoc.org/github.com/phogolabs/flaw
[godoc-img]: https://godoc.org/github.com/phogolabs/flaw?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[software-license-url]: LICENSE
