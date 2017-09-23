package bip39

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"
)

func TestGenerateEntropy(t *testing.T) {
	positive := []int{128, 160, 192, 224, 256}
	negative := []int{96, 127, 257, 288}
	var err error

	for _, p := range positive {
		assertEntropy(p, t)
	}

	for _, n := range negative {
		_, err = generateEntropy(n)
		if err == nil {
			t.Errorf("generateEntropy shouldn't work with size %v", n)
		}

		_, err = NewMnemonicRandom(n, "")
		if err == nil {
			t.Errorf("NewMnemonicRandom shouldn't work with size %v", n)
		}
	}

	negativeEnt := [][]byte{
		[]byte{},
		[]byte{0},
		[]byte{0, 0},
	}

	for _, n := range negativeEnt {
		_, err := NewMnemonicFromEntropy(n, "")
		if err == nil {
			t.Errorf("NewMnemonicFromEntropy shouldn't work with %v", n)
		}
	}
}

func TestRandomGeneration(t *testing.T) {
	a, err := NewMnemonicRandom(128, "")
	assertErr(err, t)

	b, err := NewMnemonicRandom(128, "")
	assertErr(err, t)

	senA, err := a.GetSentence()
	assertErr(err, t)

	senB, err := b.GetSentence()
	assertErr(err, t)

	if len(senA) == 0 {
		t.Error("sentence A is empty")
	}

	if len(senB) == 0 {
		t.Error("sentence B is empty")
	}

	if senA == senB {
		t.Error("two random senteces are the same")
	}

	seedA, err := a.GetSeed()
	assertErr(err, t)

	seedB, err := b.GetSeed()
	assertErr(err, t)

	if len(seedA) == 0 {
		t.Error("seed A is empty")
	}

	if len(seedB) == 0 {
		t.Error("seed B is empty")
	}

	if seedA == seedB {
		t.Error("two random seeds are the same")
	}
}

func TestVectors(t *testing.T) {
	vectors, err := ioutil.ReadFile("vectors.json")
	if err != nil {
		t.Error(err)
	}

	type vector struct {
		Arr [][]string `json:"english"`
	}
	var data vector
	err = json.Unmarshal(vectors, &data)
	if err != nil {
		t.Error(err)
	}

	if len(data.Arr) == 0 {
		t.Error(errors.New("no vectors to test"))
	}

	for _, v := range data.Arr {
		entropyHex := v[0]
		mnemonic := v[1]
		seed := v[2]
		// what is v[3] ?

		ent, err := hex.DecodeString(entropyHex)
		if err != nil {
			t.Error(err)
		}

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

		// fmt.Printf("seed pass %v", seed)
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
