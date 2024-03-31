package lexer

import (
	"io"
	"unicode"
)

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

type StateFn func(*Lexer) (StateFn, error)

func lexToken(l *Lexer) (StateFn, error) {
	r, err := l.nextChar()
	if err != nil {
		return nil, err
	}

	var state StateFn
	switch r {
	case ':':
		state = lexAssignment
	case '|':
		state = lexOr
	case ' ', '\t', '\n':
		state = lexWhitespace
	case '"', '\'', '<', '{', '[':
		state = lexEnclosedLeft
	default:
		state = lexTerminalSymbol
	}

	for state != nil {
		state, err = state(l)
		if err != nil {
			return nil, err
		}
		// Clear the temporary buffer after each state
		l.clearRuneTmpBuffer()
	}

	return nil, nil
}

func lexAction(l *Lexer) (StateFn, error) {
	if err := skipWhitespace(l); err != nil {
		return nil, err
	}

	if err := readNextCharWhile(l, func(r rune) bool {
		return r != '('
	}); err != nil {
		return nil, err
	}

	l.emitToken(Action)

	if err := advanceChar(l); err != nil {
		return nil, err
	}

	l.emitToken(ParenLeft)

	return lexActionArg, nil
}

func lexActionArg(l *Lexer) (StateFn, error) {
	if err := skipWhitespace(l); err != nil {
		return nil, err
	}

	var last rune
	if err := readNextCharWhile(l, func(r rune) bool {
		last = r
		return r != ')' && r != ',' && r != ' ' && r != '\t'
	}); err != nil {
		return nil, err
	}

	switch last {
	case ')':
		if len(l.runeTmpBuffer) > 0 {
			l.emitToken(ActionArg)
		}
		break
	case ',', ' ', '\t':
		if len(l.runeTmpBuffer) > 0 {
			l.emitToken(ActionArg)
		}
		if err := advanceCharN(l, 2); err != nil {
			return nil, err
		}
		return lexActionArg, nil
	}

	if err := advanceChar(l); err != nil {
		return nil, err
	}

	l.emitTokenOpts(string(l.runeTmpBuffer[len(l.runeTmpBuffer)-1]), l.line, l.column, ParenRight)

	return nil, nil
}

func lexEnclosedLeft(l *Lexer) (StateFn, error) {
	var tokenType TokenType
	var r1, expected rune

	r1 = l.buffer[0]
	l.clearRuneTmpBuffer()

	switch r1 {
	case '{':
		if err := advanceChar(l); err != nil {
			return nil, err
		}
		return lexAction, nil
	case '<':
		tokenType = NonTerminalSymbol
		expected = '>'
	case '[':
		tokenType = TerminalSymbol
		expected = ']'
	case '"', '\'':
		tokenType = String
		expected = r1
	default:
		return nil, ErrUnexpectedRune
	}

	err := readNextCharWhile(l, func(r rune) bool {
		return r != expected
	})
	switch err {
	case nil:
		break
	case io.EOF:
		l.emitToken(tokenType)
		return nil, nil
	default:
		return nil, err
	}

	l.emitToken(tokenType)

	err = expectChar(l, expected)
	if err != nil {
		return nil, err
	}

	err = advanceIfChar(l, func(r rune) bool {
		return r == expected
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func lexTerminalSymbol(l *Lexer) (StateFn, error) {
	err := readNextCharWhile(l, func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	})
	switch err {
	case nil:
		break
	case io.EOF:
		l.emitToken(TerminalSymbol)
		return nil, nil
	default:
		return nil, err
	}

	err = advanceIfChar(l, func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	})
	if err != nil {
		return nil, err
	}

	l.emitToken(TerminalSymbol)

	return nil, nil
}

func lexAssignment(l *Lexer) (StateFn, error) {
	r, err := l.nextChar()
	if err != nil {
		return nil, err
	}
	if r != ':' {
		return nil, ErrUnexpectedRune
	}

	r, err = l.nextChar()
	if err != nil {
		return nil, err
	}
	if r != '=' {
		return nil, ErrUnexpectedRune
	}

	l.emitToken(Assignment)

	return nil, nil
}

func lexOr(l *Lexer) (StateFn, error) {
	l.emitToken(Or)

	return nil, nil
}

func lexWhitespace(l *Lexer) (StateFn, error) {
	if err := skipWhitespace(l); err != nil {
		return nil, err
	}

	return lexToken, nil
}