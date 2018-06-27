# TablePrinter

_tableprinter_ is an intelligent value-to-table formatter and writer. It uses a customized version of [olekukonko/tablewriter](https://github.com/kataras/tablewriter) to render a result table.

![](color.png)

It checks every in data and transforms those data(structure values, slices, maps, single lists that may contain different type of values such as go standard values like `int`, `string` even a value that implements the `fmt.Stringer` interface) to a table formatted text and outputs it to an `io.Writer`. Like `encoding/json#Unmarshal` but for tables.

[![build status](https://img.shields.io/travis/kataras/tableprinter/master.svg?style=flat-square)](https://travis-ci.org/kataras/tableprinter) [![report card](https://img.shields.io/badge/report%20card-a%2B-ff3333.svg?style=flat-square)](http://goreportcard.com/report/kataras/tableprinter)

## Installation

The only requirement is the [Go Programming Language](https://golang.org/dl), at least version **1.10+**.

```sh
$ go get -u github.com/kataras/tableprinter
```

## Versioning

Current: **v0.0.1**

Read more about Semantic Versioning 2.0.0

- http://semver.org/
- https://en.wikipedia.org/wiki/Software_versioning
- https://wiki.debian.org/UpstreamGuide#Releases_and_Versions

## License

Distributed under MIT, See [LICENSE](LICENSE) for more information.