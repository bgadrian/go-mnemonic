/*Package bip39 is an immutable class that represents a BIP39 Mnemonic code.
  See BIP39 specification for more info: https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
  A Mnemonic code is a a group of easy to remember words used for the generation
  of deterministic wallets. A Mnemonic can be used to generate a seed using
  an optional passphrase, for later generate a HDPrivateKey. */
package bip39

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const bitsInByte = 8
const minEnt = 128
const maxEnt = 256
const multiple = 32
const wordBits = 11

//Mnemonic ...
type Mnemonic struct {
	ent        []byte
	passphrase string
	sentence   string
}

/*NewMnemonicRandom generate a group of easy to remember words
 -- for the generation of deterministic wallets.
use size 128 for a 12 words code.*/
func NewMnemonicRandom(size int, passphrase string) (code *Mnemonic, e error) {
	//we generate ENT count of random bits
	ent, err := generateEntropy(size)
	if err != nil {
		e = err
		return
	}

	code = &Mnemonic{}
	code.ent = ent
	code.passphrase = passphrase

	return
}

//NewMnemonicFromEntropy ...
func NewMnemonicFromEntropy(ent []byte, passphrase string) (code *Mnemonic, err error) {
	bitsCount := len(ent) * bitsInByte
	err = validBitsCount(bitsCount)
	if err != nil {
		return
	}

	code = &Mnemonic{}
	code.ent = ent
	code.passphrase = passphrase
	return
}

//newMnemonicFromSentence ...
func newMnemonicFromSentence(sentence string, passphrase string) (code *Mnemonic, e error) {
	if SentenceValid(sentence) == false {
		return nil, errors.New("mnemonic is invalid")
	}

	words := strings.Split(sentence, " ")
	bitsCount := len(words) * wordBits
	e = validBitsCount(bitsCount)
	if e != nil {
		return
	}

	checksumSize := bitsCount % multiple
	groups := make([]int, len(words))

	for i, word := range words {
		wordIndex, err := dictionaryWordToIndex(word)
		if err != nil {
			return nil, err
		}
		groups[i] = wordIndex
	}

	binWithChecksum := ""
	for _, b := range groups {
		binWithChecksum = binWithChecksum + fmt.Sprintf("%08b", b)
	}

	if len(binWithChecksum) != bitsCount {
		return nil, fmt.Errorf("internal error, wrong checksum from %v",
			sentence)
	}

	//bits of the checksum, as string
	// checksum := binWithChecksum[len(binWithChecksum)-checksumSize:]
	//TODO check this for validity

	//entropy as a string of bits
	bin := binWithChecksum[:len(binWithChecksum)-checksumSize]
	ent := make([]byte, bitsCount/bitsInByte)

	var byteAsBinaryStr string
	for i := 0; i < len(ent); i += bitsInByte {
		startIndex := i * bitsInByte
		endIndex := startIndex + bitsInByte + 1
		if endIndex >= len(bin)-1 {
			byteAsBinaryStr = bin[startIndex:]
		} else {
			byteAsBinaryStr = bin[startIndex:endIndex]
		}
		asInt64, err := strconv.ParseInt(byteAsBinaryStr, 2, 64)
		if err != nil {
			return nil, err
		}
		ent[i] = byte(asInt64)
	}

	code = &Mnemonic{}
	code.ent = ent
	code.passphrase = passphrase
	code.sentence = sentence
	return
}

//GetSentence ...
func (m *Mnemonic) GetSentence() (string, error) {
	if len(m.sentence) != 0 {
		return m.sentence, nil
	}

	bin := ""
	for _, b := range m.ent {
		bin = bin + fmt.Sprintf("%08b", b)
	}

	checksum, err := checksumEntropy(m.ent)
	if err != nil {
		return "", err
	}

	bin = bin + checksum

	wordCount := len(bin) / wordBits
	if len(bin)%wordBits != 0 {
		err := fmt.Errorf("internal error, canot divide ENT to %v groups", wordBits)
		return "", err
	}

	groups := make([]int, wordCount)
	var str string
	for i := 0; i < wordCount; i++ {
		startIndex := i * wordBits
		endIndex := startIndex + wordBits
		if endIndex >= len(bin) {
			str = bin[startIndex:]
		} else {
			str = bin[startIndex:endIndex]
		}
		asInt, err := strconv.ParseInt(str, 2, 64)
		if err != nil {
			return "", err
		}
		groups[i] = int(asInt)
	}

	words := make([]string, wordCount)
	for i, wordIndex := range groups {
		words[i], err = dictionaryIndexToWord(wordIndex)
		if err != nil {
			return "", err
		}
	}

	m.sentence = strings.Join(words, " ")

	return m.sentence, nil
}

//GetSeed ...
func (m *Mnemonic) GetSeed() (seed string, e error) {

	sentence, err := m.GetSentence()
	if err != nil {
		e = err
		return
	}
	s := NewSeed(sentence, m.passphrase)
	seed = hex.EncodeToString(s)
	return
}

//NewSeed ...
func NewSeed(mnecmonic, passphrase string) []byte {
	return pbkdf2.Key([]byte(mnecmonic), []byte("mnemonic"+passphrase), 2048, 64, sha512.New)
}

func generateEntropy(bitsCount int) (ent []byte, err error) {
	err = validBitsCount(bitsCount)
	if err != nil {
		return
	}

	bytesCount := bitsCount / bitsInByte
	ent = make([]byte, bytesCount)
	_, err = rand.Read(ent)
	return
}

func validBitsCount(bitsCount int) error {
	if bitsCount < minEnt || bitsCount > maxEnt || bitsCount%multiple != 0 {
		return fmt.Errorf(
			"entropy must between %v-%v and be divisible by %v",
			minEnt, maxEnt, multiple)
	}
	return nil
}

/*checksumEntropy A checksum is generated by taking the first
ENT / 32 bits of its SHA256 hash.*/
func checksumEntropy(ent []byte) (string, error) {
	hash := sha256.New()
	_, err := hash.Write(ent)

	//sha256.Write never seems to return error
	if err != nil {
		return "", err
	}

	sum := hash.Sum(nil)
	bits := len(ent) * bitsInByte
	cs := bits / multiple

	hashbits := ""
	for _, b := range sum {
		hashbits = hashbits + fmt.Sprintf("%08b", b)
	}

	if len(hashbits) != 256 {
		return "", errors.New("internal error, sha256 doesnt have 256 bits")
	}
	return hashbits[:cs], nil
}

//SentenceValid ...
func SentenceValid(s string) bool {
	//TODO
	return true
}
