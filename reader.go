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

type Word struct {
	Text  string
	No    uint16
	Extra []byte
}

var (
	ErrInvalidFormat = errors.New("invalid dictionary format")
)
