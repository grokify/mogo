SimpleGo
========

[![Used By][used-by-svg]][used-by-url]
[![Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

## Overview

The `simplego` package provides a collection of Go utilities for common tasks.

Go is a small language which is useful from a development and maintenance
perspective but it can result in common tasks being more verbose than other 
languages where there are more productivity methods. The `simplego` package's
goal is to provide wrappers for common tasks in the same spirit of `io/ioutil`
to make programming Go a bit faster.

In terms of scope, `simplego` is designed to cover all common areas.

## Documentation

Documentation is provided using godoc and available on [GoDoc.org](https://godoc.org/github.com/grokify/simplego).

- [crypto](https://pkg.go.dev/github.com/grokify/simplego/crypto)
- [io](https://pkg.go.dev/github.com/grokify/simplego/io)
- [log](https://pkg.go.dev/github.com/grokify/simplego/log)
- [net](https://pkg.go.dev/github.com/grokify/simplego/net)
  - [net/httputilmore](https://pkg.go.dev/github.com/grokify/simplego/net/httputilmore)
  - [net/urlutil](https://pkg.go.dev/github.com/grokify/simplego/net/urlutil)
- [os](https://pkg.go.dev/github.com/grokify/simplego/os)
- mime
  - [multipart](https://pkg.go.dev/github.com/grokify/simplego/mime/multipart)
- [sort](https://pkg.go.dev/github.com/grokify/simplego/sort)
- [strconv](https://pkg.go.dev/github.com/grokify/simplego/strconv)
- [text](https://pkg.go.dev/github.com/grokify/simplego/text)
- [time](https://pkg.go.dev/github.com/grokify/simplego/time)
  - [timeutil](https://pkg.go.dev/github.com/grokify/simplego/time/timeutil)
- [type](https://pkg.go.dev/github.com/grokify/simplego/type)

## Installation

```bash
$ go get github.com/grokify/simplego/...
```

## Contributing

Features, Issues, and Pull Requests are always welcome.

To contribute:

1. Fork it ( http://github.com/grokify/simplego/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

Please report issues and feature requests on [Github](https://github.com/grokify/simplego).

 [used-by-svg]: https://sourcegraph.com/github.com/grokify/simplego/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/grokify/simplego?badge
 [build-status-svg]: https://github.com/grokify/simplego/workflows/go%20build/badge.svg?branch=master
 [build-status-url]: https://github.com/grokify/simplego/actions
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/simplego
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/simplego
 [codeclimate-status-svg]: https://codeclimate.com/github/grokify/simplego/badges/gpa.svg
 [codeclimate-status-url]: https://codeclimate.com/github/grokify/simplego
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/simplego
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/simplego
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/simplego/blob/master/LICENSE
