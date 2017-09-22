package bip39

import (
	"encoding/hex"
	"testing"
)

func TestGenerateEntropy(t *testing.T) {

	positive := []int{128, 160, 192, 224, 256}
	negative := []int{96, 127, 257, 288}

	for _, p := range positive {
		assertEntropy(p, t)
	}

	for _, n := range negative {
		_, err := generateEntropy(n)
		if err == nil {
			t.Errorf("generateEntropy shouldn't work with size %v", n)
		}
	}
}

func TestGenerateMnemonic(t *testing.T) {
	entropyHex := "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f"
	mnemonic := "legal winner thank year wave sausage worth useful legal winner thank yellow"
	seed := "2e8905819b8723fe2c1d161860e5ee1830318dbf49a83bd451cfb8440c28bd6fa457fe1296106559a3c80937a1c1069be3a3a5bd381ee6260e8d9739fce1f607"

	ent, err := hex.DecodeString(entropyHex)
	if err != nil {
		t.Error(err)
	}

	// ent, _ = generateEntropy(128)

	r, err := NewMnemonicFromEntropy(ent, "TREZOR")
	if err != nil {
		t.Error(err)
	}

	s, err := r.GetSentence()
	if err != nil {
		t.Error(err)
	}

	if s != mnemonic {
		t.Errorf("exp %v got %v for %v",
			mnemonic, s, entropyHex)
	}
	se, err := r.GetSeed()
	if err != nil {
		t.Error(err)
	}

	if se != seed {
		t.Errorf("exp %v got %v for %v",
			seed, se, entropyHex)
	}
}

// func unhexlify(hex string) (binary string, e error) {
// 	if len(hex)%2 != 0 {
// 		return "", errors.New("param must be an even number of hex digits")
// 	}

// 	for i := 0; i < len(hex); i += 2 {
// 		char := fmt.Sprintf("%s", hex[i:i+2])
// 		v64, err := strconv.ParseInt(char, 16, 64)
// 		if err != nil {
// 			return "", err
// 		}
// 		binary = binary + fmt.Sprintf("%08b", v64)
// 	}

// 	return
// }

func TestDictionary(t *testing.T) {
	type tuple struct {
		dict []string
		line int
		word string
	}

	en, err := dictionary()
	if err != nil {
		t.Error(err)
	}

	table := []tuple{
		{en, 1, "abandon"},
		{en, 655, "fade"},
		{en, 2048, "zoo"},
	}

	for _, tuple := range table {
		index := tuple.line - 1
		if tuple.dict[index] != tuple.word {
			t.Errorf("wrong word for line %v, expected %v got %v",
				tuple.line, tuple.word, tuple.dict[index])
		}
	}
}

func assertEntropy(size int, t *testing.T) {
	a, err := generateEntropy(size)
	assertErr(err, t)
	if emptyBytes(a) {
		t.Errorf("generateEntropy empty for %v", size)
	}

	count := len(a) * bitsInByte
	if count != size {
		t.Errorf("generateEntropy wrong number of bites for %v, exp: %v got %v",
			size, size, count)
	}
}

func assertErr(err error, t *testing.T) {
	if err == nil {
		return
	}
	t.Error(err)
}

func emptyBytes(slice []byte) bool {
	for _, b := range slice {
		if b != 0 {
			return false
		}
	}
	return true
}
