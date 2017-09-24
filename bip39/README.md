# Bitcoin BIP39 - Mnemonic code for generating deterministic keys  [![Build Status](https://travis-ci.org/bgadrian/go-mnemonic.svg?branch=master)](https://travis-ci.org/bgadrian/go-mnemonic) [![codecov](https://codecov.io/gh/bgadrian/go-mnemonic/branch/master/graph/badge.svg)](https://codecov.io/gh/bgadrian/go-mnemonic) [![Go Report Card](https://goreportcard.com/badge/github.com/bgadrian/go-mnemonic)](https://goreportcard.com/report/github.com/bgadrian/go-mnemonic/bip39)  [![GoDoc](https://godoc.org/golang.org/x/tools/cmd/godoc?status.svg)](https://godoc.org/github.com/bgadrian/go-mnemonic/bip39)

Implementation in Go (golang) based on [original specs](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki#Reference_Implementation).

> This BIP describes the implementation of a mnemonic code or mnemonic sentence -- a group of easy to remember words -- for the generation of deterministic wallets. It consists of two parts: generating the mnemonic, and converting it into a binary seed. This seed can be later used to generate deterministic wallets using BIP-0032 or similar methods


### Example
A basic example, more on go docs.
```go
package main

import (
	"fmt"
	"log"
	"strings"
	"github.com/bgadrian/go-mnemonic/bip39"
)

func main() {
	/*
		128 bits -> 12 words
		160 bits -> 15 words
		192 bits -> 18 words
		224 bits -> 21 words
		256 bits -> 24 words
	*/
	//generate a new random 12 english words password (mnemonic)
	newRandomMnemonic, err := bip39.NewMnemonicRandom(128, "")
	if err != nil {
		log.Panic(err)
	}

	password, err := newRandomMnemonic.GetSentence()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("the password has %v words\n", len(strings.Split(password, " ")))

	// Output:
	//the password has 12 words
}
```

### TODO

* create a Mnemonic from a seed
* implement all word lists

Currently the library only supports the English dictionary, if enough interests is shown it's easy to implement the rest of the [word lists](https://github.com/bitcoin/bips/tree/master/bip-0039).
