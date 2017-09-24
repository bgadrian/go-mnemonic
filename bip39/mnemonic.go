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
	//entropy without the checksum
	ent        []byte
	passphrase string
	sentence   string
}

/*NewMnemonicRandom creates a new random (crypto safe) Mnemonic.Use size 128 for a 12 words code.*/
func NewMnemonicRandom(size int, passphrase string) (code *Mnemonic, e error) {
	//we generate ENT count of random bits
	ent, err := generateRandomEntropy(size)
	if err != nil {
		e = err
		return
	}

	code = &Mnemonic{}
	code.ent = ent
	code.passphrase = passphrase

	return
}

//NewMnemonicFromEntropy Generates a Mnemonic based on a known entropy (stored as hex bytes)
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

//NewMnemonicFromSentence Generates a menmonic based on a known code (list of words).
func NewMnemonicFromSentence(sentence string, passphrase string) (code *Mnemonic, e error) {
	words := strings.Split(sentence, " ")
	bitsCountWithCheksum := len(words) * wordBits
	checksumBitsCount := bitsCountWithCheksum % multiple
	bitsCount := bitsCountWithCheksum - checksumBitsCount

	e = validBitsCount(bitsCount)
	if e != nil {
		return
	}

	//ent as string of bits
	binWithChecksum := ""
	for _, word := range words {
		wordIndex, err := dictionaryWordToIndex(strings.Trim(word, ""))
		if err != nil {
			return nil, err
		}
		binWithChecksum = binWithChecksum + fmt.Sprintf("%011b", wordIndex)
	}

	if len(binWithChecksum) != bitsCountWithCheksum {
		return nil, fmt.Errorf("internal error, wrong checksum got %v bits, expected %v, sentence:'%v'",
			len(binWithChecksum), bitsCountWithCheksum, sentence)
	}

	//entropy without the checksum
	bin := binWithChecksum[:bitsCount]

	if len(bin) != bitsCount {
		return nil, fmt.Errorf("internal error, bits count for '%v' is wrong, got %v bits, exp %v",
			sentence, len(binWithChecksum), bitsCount)
	}
	//entropy as a string of bits
	ent := make([]byte, bitsCount/bitsInByte)

	var byteAsBinaryStr string
	for i := 0; i < len(ent); i++ {
		startIndex := i * bitsInByte
		endIndex := startIndex + bitsInByte
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
	//bits of the checksum, as string
	// checksum := binWithChecksum[len(binWithChecksum)-checksumSize:]
	//TODO add checksum to entropy

	code = &Mnemonic{}
	code.ent = ent
	code.passphrase = passphrase
	code.sentence = sentence
	return
}

//GetSentence Return the words from this Mnemonic.
func (m *Mnemonic) GetSentence() (string, error) {
	if len(m.sentence) != 0 {
		return m.sentence, nil
	}

	bin := ""
	for _, b := range m.ent {
		bin = bin + fmt.Sprintf("%08b", b)
	}

	checksum, err := generateChecksumEntropy(m.ent)
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
		if endIndex >= len(bin)-1 {
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

//GetSeed Returns the seed for this Mnemonic (as hex in string)
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

//NewSeed Based on a code (word list) returns the seed (hex bytes)
func NewSeed(mnecmonic, passphrase string) []byte {
	return pbkdf2.Key([]byte(mnecmonic), []byte("mnemonic"+passphrase), 2048, 64, sha512.New)
}

//GetEntropyStrHex get the entryope as hex in a string, for easy storage
func (m *Mnemonic) GetEntropyStrHex() (string, error) {
	if len(m.ent) == 0 {
		return "", errors.New("empty entropy")
	}

	return hex.EncodeToString(m.ent), nil
}

func generateRandomEntropy(bitsCount int) (ent []byte, err error) {
	err = validBitsCount(bitsCount)
	if err != nil {
		return
	}

	bytesCount := bitsCount / bitsInByte
	ent = make([]byte, bytesCount)
	_, err = rand.Read(ent)
	return
}

//Based on BIP39 specifications
func validBitsCount(bitsCount int) error {
	if bitsCount < minEnt || bitsCount > maxEnt || bitsCount%multiple != 0 {
		return fmt.Errorf(
			"entropy must between %v-%v and be divisible by %v, but got %v bits",
			minEnt, maxEnt, multiple, bitsCount)
	}
	return nil
}

/*generateChecksumEntropy A checksum is generated by taking the first
ENT / 32 bits of its SHA256 hash.*/
func generateChecksumEntropy(ent []byte) (string, error) {
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
