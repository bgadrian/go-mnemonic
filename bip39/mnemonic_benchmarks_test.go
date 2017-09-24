package bip39

import "testing"
import "encoding/hex"

func BenchmarkNewMnemonicRandom128(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewMnemonicRandom(128, "")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkNewMnemonicFromEntropy(b *testing.B) {
	ent, err := hex.DecodeString("80808080808080808080808080808080")
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewMnemonicFromEntropy(ent, "")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkNewMnemonicFromSentence(b *testing.B) {
	password := "letter advice cage absurd amount doctor acoustic avoid letter advice cage above"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := NewMnemonicFromSentence(password, "")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGetSentence(b *testing.B) {
	ent, _ := hex.DecodeString("80808080808080808080808080808080")
	m, _ := NewMnemonicFromEntropy(ent, "")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.GetSentence()
	}
}

func BenchmarkGetSeed(b *testing.B) {
	ent, _ := hex.DecodeString("80808080808080808080808080808080")
	m, _ := NewMnemonicFromEntropy(ent, "")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.GetSeed()
	}
}

func BenchmarkGetEntropyStrHex(b *testing.B) {
	password := "letter advice cage absurd amount doctor acoustic avoid letter advice cage above"
	m, _ := NewMnemonicFromSentence(password, "")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.GetEntropyStrHex()
	}
}
