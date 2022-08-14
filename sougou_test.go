package pydict

import (
	"os"
	"testing"
)

func TestSougouPinyinReader_Read(t *testing.T) {
	file, err := os.Open(`input.scel`)
	if err != nil {
		t.Fatalf("open sougou pinyin dictionary failed due to %v", err)
	}
	defer file.Close()

	r := NewSougouPinyinReader()
	dictionary, err := r.Read(file)
	if err != nil {
		t.Fatalf("read sougou pinyin dictionary failed due to %v", err)
	}

	t.Logf("pinyin: %d", len(dictionary.Pinyin))
}
