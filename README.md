MoGo
====

[![Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![CodeClimate][codeclimate-status-svg]][codeclimate-status-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![SLOC][loc-svg]][repo-url]
[![License][license-svg]][license-url]

## Overview

The `mogo` (pronounced "Mo Go" for "More Go") package provides a collection of Go utilities for common tasks.

Go is a small language which is useful from a development and maintenance
perspective but it can result in common tasks being more verbose than other 
languages where there are more productivity methods. The `mogo` package's
goal is to provide wrappers for common tasks in the same spirit of `io/ioutil`
to make programming Go a bit faster.

In terms of scope, `mogo` is designed to cover all common areas, with a move to fewer dependencies.

See the following modules for extended functionality (with more dependencies):

* GoAuth - [github.com/grokify/goauth](https://github.com/grokify/goauth)
* GoCharts - [github.com/grokify/gocharts](https://github.com/grokify/gocharts)
* GoPhoneNumbers - [github.com/grokify/gophonenumbers](https://github.com/grokify/gophonenumbers)

## Documentation

Documentation is provided using godoc and available on [GoDoc.org](https://godoc.org/github.com/grokify/mogo).

## Installation

```bash
$ go get github.com/grokify/mogo/...
```

## SBOM

The SBOM is in the file [`SBOM.spdx.json`]. It is created using GitHub's built-in SBOM feature with the following additions:

1. It is converted to indented JSON for easier reading.
1. The `packages[0].downloadLocation` value as the commit hash added per [https://spdx.github.io/spdx-spec/v2.3/package-information/](downloadLocation).

## Credits

See [`CREDITS`](CREDITS.md)

## Contributing

Features, Issues, and Pull Requests are always welcome.

To contribute:

1. Fork it ( http://github.com/grokify/mogo/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

Please report issues and feature requests on [Github](https://github.com/grokify/mogo).

 [used-by-svg]: https://sourcegraph.com/github.com/grokify/mogo/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/grokify/mogo?badge
 [build-status-svg]: https://github.com/grokify/mogo/workflows/test/badge.svg?branch=master
 [build-status-url]: https://github.com/grokify/mogo/actions/workflows/test.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/mogo
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/mogo
 [codeclimate-status-svg]: https://codeclimate.com/github/grokify/mogo/badges/gpa.svg
 [codeclimate-status-url]: https://codeclimate.com/github/grokify/mogo
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/mogo
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/mogo
 [license-svg]: https://img.shields.io/badge/license-MIT-mogo.svg
 [license-url]: https://github.com/grokify/mogo/blob/master/LICENSE
 [loc-svg]: https://tokei.rs/b1/github/grokify/mogo
 [repo-url]: https://github.com/grokify/mogo
