package bip39

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"
)

func ExampleNewMnemonicRandom() {
	/*
		128 bits -> 12 words
		160 bits -> 15 words
		192 bits -> 18 words
		224 bits -> 21 words
		256 bits -> 24 words
	*/
	//generate a new random 12 english words password (mnemonic)
	newRandomMnemonic, err := NewMnemonicRandom(128, "")
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

func ExampleNewMnemonicFromEntropy() {
	//we ignore the errors returned for Example simplicity sake, don't do this on, production
	entHexBytes := []byte{94, 127, 194, 94, 217, 163, 84, 91, 112, 158, 206, 144, 80, 5, 219, 134}

	mnemonic, _ := NewMnemonicFromEntropy(entHexBytes, "")
	password, _ := mnemonic.GetSentence()
	seed, _ := mnemonic.GetSeed()

	fmt.Printf("entropy can be stored as '%v'\n", hex.EncodeToString(entHexBytes))
	fmt.Printf("password is '%v'\n", password)
	fmt.Printf("seed is '%v'\n", seed)

	// Output:
	//entropy can be stored as '5e7fc25ed9a3545b709ece905005db86'
	//password is 'fury wrap nut rebuild crystal color second supply motion lens ivory around'
	//seed is 'b5f87c4020dde5e83f73dd89ceb49b2700437008fe6593dd675ea856be1d687e3bc17a4cf9c070f7e469d704942f137ae4eea7ad2c9189edffd991e0075b44ee'
}

func ExampleNewMnemonicFromSentence() {
	//we ignore the errors returned for Example simplicity sake, don't do this on, production
	password := "fury wrap nut rebuild crystal color second supply motion lens ivory around"

	mnemonic, _ := NewMnemonicFromSentence(password, "")
	seed, _ := mnemonic.GetSeed()
	ent, _ := mnemonic.GetEntropyStrHex()

	fmt.Printf("seed is '%v'\n", seed)
	fmt.Printf("entropy can be stored as '%v'\n", ent)

	// Output:
	//seed is 'b5f87c4020dde5e83f73dd89ceb49b2700437008fe6593dd675ea856be1d687e3bc17a4cf9c070f7e469d704942f137ae4eea7ad2c9189edffd991e0075b44ee'
	//entropy can be stored as '5e7fc25ed9a3545b709ece905005db86'
}
