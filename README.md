# flaw

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][travis-img]][travis-url]
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
err := flaw.Wrap(sql.ErrNoRows).WithCode(404)
```

Collecting multiple errors:

```golang
errs := flaw.ErrorCollector{}
errs = append(errs, fmt.Errorf("oh no"))
```

## Contributing

We are welcome to any contributions. Just fork the
[project](https://github.com/phogolabs/flaw).

[travis-img]: https://travis-ci.org/phogolabs/flaw.svg?branch=master
[travis-url]: https://travis-ci.org/phogolabs/flaw
[report-img]: https://goreportcard.com/badge/github.com/phogolabs/flaw
[report-url]: https://goreportcard.com/report/github.com/phogolabs/flaw
[codecov-url]: https://codecov.io/gh/phogolabs/flaw
[codecov-img]: https://codecov.io/gh/phogolabs/flaw/branch/master/graph/badge.svg
[godoc-url]: https://godoc.org/github.com/phogolabs/flaw
[godoc-img]: https://godoc.org/github.com/phogolabs/flaw?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[software-license-url]: LICENSE
