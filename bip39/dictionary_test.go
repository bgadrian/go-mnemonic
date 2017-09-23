package bip39

import "testing"

//TODO import and test all languages
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
