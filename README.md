# Topographly

[![Build Status](https://travis-ci.com/benmaddison/topographly.svg?branch=master)](https://travis-ci.com/benmaddison/topographly)
[![Coverage Status](https://coveralls.io/repos/github/benmaddison/topographly/badge.svg?branch=master)](https://coveralls.io/github/benmaddison/topographly?branch=master)

An example program for creating versioned network topology datastructures in Go with
- GraphQL API
- Noms datastore
- YANG data modelling and validation

## Getting started

If you are running in docker, enable ipv6 support, as the http server will try
to bind to `http://[::]:8080`.

1. Install:
```bash
$ go get github.com/benmaddison/topographly/cmd/topographly
```
2. Run:
```bash
$ NOMS_VERSION_NEXT=1 $GOPATH/bin/topographly
```
3. Open a browser, and navigate to `http://localhost:8080/graphql`.

## Disclaimer
This program is a toy for experimenting with this combination of database,
schema language and API handling library. It is not intended for production use
for any purpose.
