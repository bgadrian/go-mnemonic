package bip39

import "testing"

//TODO import and test all languages
func TestDictionaryIndex(t *testing.T) {
	type tuple struct {
		index int
		word  string
	}

	table := []tuple{
		//line in file - 1 being 0 index based
		{1 - 1, "abandon"},
		{655 - 1, "fade"},
		{2048 - 1, "zoo"},
	}

	for _, d := range table {
		index, err := dictionaryWordToIndex(d.word)
		if err != nil {
			t.Error(err)
		}

		if index != d.index {
			t.Errorf("wrong index for word %v, expected %v got %v",
				d.word, d.index, index)
		}

		word, err := dictionaryIndexToWord(d.index)
		if err != nil {
			t.Error(err)
		}

		if word != d.word {
			t.Errorf("wrong word for index %v, expected %v got %v",
				d.index, d.word, word)
		}
	}
}

func TestWrongWords(t *testing.T) {
	negativeIndex := []int{-1, 2049}
	negativeWord := []string{"xxx", "-"}

	for _, ni := range negativeIndex {
		_, err := dictionaryIndexToWord(ni)
		if err == nil {
			t.Errorf("index %v should return an error",
				ni)
		}
	}

	for _, nw := range negativeWord {
		_, err := dictionaryWordToIndex(nw)
		if err == nil {
			t.Errorf("word %v should return an error",
				nw)
		}
	}

}
