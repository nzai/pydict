package pydict

import (
	"encoding/binary"
	"io"
	"strings"
	"unicode/utf16"

	"github.com/nzai/bio"
)

const (
	SougouPinyinPinyinOffset = 0x1540
	SougouPinyinWordOffset   = 0x2628
)

type SougouPinyinReader struct{}

func NewSougouPinyinReader() Reader {
	return &SougouPinyinReader{}
}

func (s SougouPinyinReader) Read(rs io.ReadSeeker) (*Dictionary, error) {
	br := bio.NewBinaryReaderOrder(rs, binary.LittleEndian)

	_, err := rs.Seek(SougouPinyinPinyinOffset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	pinyin, err := s.readPinyin(br)
	if err != nil {
		return nil, err
	}

	_, err = rs.Seek(SougouPinyinWordOffset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	word, count, err := s.readWord(br, pinyin)
	if err != nil {
		return nil, err
	}

	return &Dictionary{Pinyin: pinyin, Word: word, WordCount: count}, nil
}

func (s SougouPinyinReader) readPinyin(br *bio.BinaryReader) ([]string, error) {
	// version
	_, err := br.Bytes(4)
	if err != nil {
		return nil, err
	}

	// header and 4 bytes magic words
	pinyins := make([]string, 0)

	for offset := SougouPinyinPinyinOffset + 4; offset < SougouPinyinWordOffset; {
		index, err := br.UInt16()
		if err != nil {
			return nil, err
		}
		offset += 2

		if int(index) != len(pinyins) {
			break
		}

		pinyin, size, err := s.readUnicode16LeString(br)
		if err != nil {
			return nil, err
		}

		offset += size
		pinyins = append(pinyins, pinyin)
	}

	return pinyins, nil
}

func (s SougouPinyinReader) readWord(br *bio.BinaryReader, pinyins []string) (map[string][]*Word, int, error) {
	words := make(map[string][]*Word)
	count := 0
	for {
		groups, pinyin, err := s.readWordGroup(br, pinyins)
		if err == io.EOF || len(groups) == 0 {
			break
		}

		if err != nil {
			return nil, 0, err
		}

		words[pinyin] = groups
		count += len(groups)
	}

	return words, count, nil
}

func (s SougouPinyinReader) readWordGroup(br *bio.BinaryReader, pinyins []string) ([]*Word, string, error) {
	wordCount, err := br.UInt16()
	if err != nil {
		return nil, "", err
	}

	pinyinCount, err := br.UInt16()
	if err != nil {
		return nil, "", err
	}

	wordPinyins := make([]string, 0, pinyinCount)
	for index := 0; index < int(pinyinCount)/2; index++ {
		pinyinIndex, err := br.UInt16()
		if err != nil {
			return nil, "", err
		}

		wordPinyins = append(wordPinyins, pinyins[pinyinIndex])
	}
	pinyin := strings.Join(wordPinyins, " ")

	words := make([]*Word, 0, wordCount)
	var text string
	var no uint16
	var extra []byte
	for index := 0; index < int(wordCount); index++ {
		text, _, err = s.readUnicode16LeString(br)
		if err != nil {
			return nil, "", err
		}

		extraSize, err := br.UInt16()
		if err != nil {
			return nil, "", err
		}

		extra, err = br.Bytes(int(extraSize))
		if err != nil {
			return nil, "", err
		}

		no = br.ByteOrder().Uint16(extra[:2])

		words = append(words, &Word{
			Text:  text,
			No:    no,
			Extra: extra[2:],
		})
	}

	return words, pinyin, nil
}

func (s SougouPinyinReader) readUnicode16LeString(br *bio.BinaryReader) (string, int, error) {
	length, err := br.UInt16()
	if err != nil {
		return "", 0, err
	}

	buffer, err := br.Bytes(int(length))
	if err != nil {
		return "", 0, err
	}

	points := make([]uint16, 0, len(buffer)/2)
	for index := 0; index < len(buffer); index += 2 {
		points = append(points, binary.LittleEndian.Uint16(buffer[index:index+2]))
	}

	return string(utf16.Decode(points)), int(length), nil
}
