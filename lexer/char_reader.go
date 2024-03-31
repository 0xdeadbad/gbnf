package lexer

import (
	"bufio"
	"io"
)

type CharReader struct {
	*bufio.Reader
	buffer [1]rune
}

func newCharReader(r io.Reader) *CharReader {
	return &CharReader{
		Reader: bufio.NewReader(r),
		buffer: [1]rune{0},
	}
}

func (l *CharReader) nextChar() (rune, error) {
	r, _, err := l.ReadRune()
	if err != nil {
		return 0, err
	}

	l.buffer[0] = r

	return r, nil
}

func (l *CharReader) peekChar() (rune, error) {
	r, _, err := l.ReadRune()
	if err != nil {
		return 0, err
	}

	err = l.UnreadRune()
	if err != nil {
		return 0, err
	}

	return r, nil
}

type ErrCharReader string

const (
	ErrCharReaderZeroBytes ErrCharReader = "zero bytes read"
	ErrUnexpectedRune      ErrCharReader = "unexpected rune"
	ErrCondFailed          ErrCharReader = "condition failed"
)

func (e ErrCharReader) Error() string {
	return string(e)
}

func (e ErrCharReader) String() string {
	return string(e)
}
