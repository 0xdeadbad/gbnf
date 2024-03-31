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

func readNextCharWhile(l *Lexer, f func(r rune) bool) error {
	r, err := l.peekChar()
	if err != nil {
		return err
	}
	for f(r) {
		r, err = l.nextChar()
		if err != nil {
			return err
		}
		r, err = l.peekChar()
		if err != nil {
			return err
		}
	}

	return nil
}

func advanceIfChar(l *Lexer, f func(r rune) bool) error {
	r, err := l.peekChar()
	if err != nil {
		return err
	}
	if f(r) {
		_, err = l.nextChar()
		if err != nil {
			return err
		}
	}

	return nil
}

func expectChar(l *Lexer, r rune) error {
	if r1, err := l.peekChar(); err != nil {
		return err
	} else {
		if r1 != r {
			return ErrUnexpectedRune
		}
	}

	return nil
}

func advanceChar(l *Lexer) error {
	if _, err := l.nextChar(); err != nil {
		return err
	}

	return nil
}

func advanceCharN(l *Lexer, n int) error {
	for i := 0; i < n; i++ {
		if err := advanceChar(l); err != nil {
			return err
		}
	}

	return nil
}

func advanceClr(l *Lexer) error {
	if err := advanceChar(l); err != nil {
		return err
	}
	l.clearRuneTmpBuffer()

	return nil
}

func expectCond(l *Lexer, r rune) error {
	if r1, err := l.peekChar(); err != nil {
		return err
	} else {
		if r1 != r {
			return ErrCondFailed
		}
	}

	return nil
}

func skipWhitespace(l *Lexer) error {
	err := readNextCharWhile(l, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n'
	})

	// Discard the last whitespace character
	// In some cases the next character isn't a whitespace character,
	// for instance if there's only one space character, so we check it before discarding it
	err = advanceIfChar(l, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n'
	})
	if err != nil {
		return err
	}

	return nil
}
