package pydict

import (
	"errors"
	"io"
)

type Reader interface {
	Read(rs io.ReadSeeker) (*Dictionary, error)
}

type Dictionary struct {
	Pinyin    []string
	Word      map[string][]*Word
	WordCount int
}

func (d *Dictionary) Merge(dict *Dictionary) {
	d.mergePinyin(dict.Pinyin)
	d.mergeWord(dict.Word)
}

func (d *Dictionary) mergePinyin(pinyin []string) {
	pinyinDict := make(map[string]bool, len(d.Pinyin))
	for _, word := range d.Pinyin {
		pinyinDict[word] = true
	}

	for _, word := range pinyin {
		if pinyinDict[word] {
			continue
		}

		d.Pinyin = append(d.Pinyin, word)
	}
}

func (d *Dictionary) mergeWord(word map[string][]*Word) {
	for pinyin, words := range word {
		_, found := d.Word[pinyin]
		if !found {
			d.Word[pinyin] = make([]*Word, 1)
		}

		d.Word[pinyin] = append(d.Word[pinyin], words...)
		d.WordCount += len(words)
	}
}

type Word struct {
	Text  string
	No    uint16
	Extra []byte
}

var (
	ErrInvalidFormat = errors.New("invalid dictionary format")
)
