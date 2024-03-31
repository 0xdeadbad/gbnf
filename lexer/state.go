package lexer

import (
	"io"
	"unicode"
)

type StateFn func(*Lexer) (StateFn, error)

func lexToken(l *Lexer) (StateFn, error) {
	r, err := l.peekChar()
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

	return state, nil
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
	l.clearRuneTmpBuffer()

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

	return lexToken, nil
}

func lexEnclosedLeft(l *Lexer) (StateFn, error) {
	var tokenType TokenType
	var r1, expected rune

	r1, err := l.nextChar()
	if err != nil {
		return nil, err
	}
	l.clearRuneTmpBuffer()

	switch r1 {
	case '{':
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

	err = readNextCharWhile(l, func(r rune) bool {
		return r != expected
	})
	switch err {
	case nil:
		break
	case io.EOF:
		l.emitToken(tokenType)
		return lexToken, nil
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

	return lexToken, nil
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

	return lexToken, nil
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

	return lexToken, nil
}

func lexOr(l *Lexer) (StateFn, error) {
	r, err := l.nextChar()
	if err != nil {
		return nil, err
	}
	if r != '|' {
		return nil, ErrUnexpectedRune
	}

	l.emitToken(Or)

	return lexToken, nil
}

func lexWhitespace(l *Lexer) (StateFn, error) {
	if err := skipWhitespace(l); err != nil {
		return nil, err
	}

	return lexToken, nil
}
