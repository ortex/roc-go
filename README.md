# Go bindings for Roc

[![GoDoc](https://godoc.org/github.com/roc-project/roc-go/roc?status.svg)](https://godoc.org/github.com/roc-project/roc-go/roc) [![Build Status](https://travis-ci.org/roc-project/roc-go.svg?branch=master)](https://travis-ci.org/roc-project/roc-go) [![Coverage Status](https://coveralls.io/repos/github/roc-project/roc-go/badge.svg?branch=master)](https://coveralls.io/github/roc-project/roc-go?branch=master)

_Work in progress!_

## Install

```
go get github.com/roc-project/roc-go/roc
```

## Dependencies
You will need to have libroc and libroc-devel (headers) installed. Refer to official build [instructions](https://roc-project.github.io/roc/docs/building.html) on how to install libroc. There is no official distribution for any os as of now, you will need to install from source.

## Build
Will not produce any lib/binary. It's just to check for the syntax errors:

`make build`

## Test

`make test`
