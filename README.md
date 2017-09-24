# Go Mnemonic [![Build Status](https://travis-ci.org/bgadrian/go-mnemonic.svg?branch=master)](https://travis-ci.org/bgadrian/go-mnemonic) [![codecov](https://codecov.io/gh/bgadrian/go-mnemonic/branch/master/graph/badge.svg)](https://codecov.io/gh/bgadrian/go-mnemonic)

Golang implementation of Bitcoin & other Mnemonic algorithms used in block chains.

### Why?

Mainly for academic purposes, I want to learn Go & block chains. I only found 1 other implementation in Go so there is room for more.
The code has full unit tests and benchmarks, including all official vectors.


### [Bitcoin BIP39 - Mnemonic code for generating deterministic keys](./bip39/README.md)  [![Go Report Card](https://goreportcard.com/badge/github.com/bgadrian/go-mnemonic/bip39)](https://goreportcard.com/report/github.com/bgadrian/go-mnemonic/bip39)  [![GoDoc](https://godoc.org/golang.org/x/tools/cmd/godoc?status.svg)](https://godoc.org/github.com/bgadrian/go-mnemonic/bip39)

Implementation in Go (golang) based on [original specs](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki#Reference_Implementation).

> This BIP describes the implementation of a mnemonic code or mnemonic sentence -- a group of easy to remember words -- for the generation of deterministic wallets.

### Copyright 

Adrian B.G. 2017, no rights reserved.