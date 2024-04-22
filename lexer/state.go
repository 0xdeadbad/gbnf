package lexer

import (
	"io"
	"unicode"
)

type StateFn func(*Lexer) (StateFn, error)

func lexToken(l *Lexer) (StateFn, error) {
	r, err := l.peekChar()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	var state StateFn
	switch r {
	case '!':
		state = lexEndOfRule
	case ':':
		state = lexAssignment
	case '|':
		state = lexOr
	case '=':
		state = lexAssign
	case '.':
		state = lexSequence
	case '&':
		state = lexAnd
	case ' ', '\t', '\n':
		state = lexWhitespace
	case '"', '\'', '<', '{', '[', '(':
		state = lexEnclosedLeft
	default:
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) {
			state = lexTerminalSymbol
		} else {
			return nil, ErrUnexpectedRune
		}
	}

	return state, nil
}

func lexAnd(l *Lexer) (StateFn, error) {
	if err := advanceChar(l); err != nil {
		return nil, err
	}

	l.emitToken(And)

	return lexToken, nil
}

func lexEndOfRule(l *Lexer) (StateFn, error) {
	if err := advanceChar(l); err != nil {
		return nil, err
	}

	if err := expectChar(l, '!'); err != nil {
		if err == ErrUnexpectedRune {
			l.emitToken(Not)
			return lexToken, nil
		}
		return nil, err
	}

	if err := advanceIfChar(l, func(r rune) bool {
		return r == '!'
	}); err != nil {
		return nil, err
	}

	l.emitToken(EndOfRule)

	return lexToken, nil
}

func lexSequence(l *Lexer) (StateFn, error) {

	if err := advanceChar(l); err != nil {
		return nil, err
	}
	if err := advanceIfChar(l, func(r rune) bool {
		return r == '.'
	}); err != nil {
		return nil, err
	}
	if err := advanceIfChar(l, func(r rune) bool {
		return r == '.'
	}); err != nil {
		return nil, err
	}

	l.emitToken(Sequence)

	return lexToken, nil
}

func lexAssign(l *Lexer) (StateFn, error) {
	if err := advanceChar(l); err != nil {
		return nil, err
	}

	l.emitToken(Assign)

	return lexToken, nil
}

func lexAction(l *Lexer) (StateFn, error) {
	if err := skipWhitespace(l); err != nil {
		return nil, err
	}

	l.clearRuneTmpBuffer()

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

	// l.emitToken(ParenLeft)

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

	// l.emitTokenOpts(string(l.runeTmpBuffer[len(l.runeTmpBuffer)-1]), l.line, l.column, ParenRight)

	return lexToken, nil
}

func lexEnclosedLeft(l *Lexer) (StateFn, error) {
	r1, err := l.nextChar()
	if err != nil {
		return nil, err
	}
	l.clearRuneTmpBuffer()

	switch r1 {
	case '{':
		return lexAction, nil
	case '"', '\'', '<':
		if r1 == '<' {
			l.runeStk.Push('>')
		} else {
			l.runeStk.Push(r1)
		}
		return lexString, nil
	case '[', '(':
		l.runeStk.Push(r1)
		return lexGroup, nil
	}

	return nil, ErrUnexpectedRune
}

func lexEnclosedRight(l *Lexer) (StateFn, error) {
	opening := l.runeStk.Pop()

	var expected rune

	switch opening {
	case '(':
		expected = ')'
	case '[':
		expected = ']'
	default:
		return nil, ErrUnexpectedRune
	}

	if err := expectChar(l, expected); err != nil {
		return nil, err
	}

	if err := advanceChar(l); err != nil {
		return nil, err
	}

	return lexToken, nil
}

func lexString(l *Lexer) (StateFn, error) {
	expected := l.runeStk.Pop()
	if err := readNextCharWhile(l, func(r rune) bool {
		return r != expected
	}); err != nil {
		return nil, err
	}

	if err := expectChar(l, expected); err != nil {
		return nil, err
	}

	if expected == '>' {
		l.emitToken(NonTerminalSymbol)
	} else {
		l.emitToken(TerminalSymbol)
	}

	if err := advanceIfChar(l, func(r rune) bool {
		return r == expected
	}); err != nil {
		return nil, err
	}

	return lexToken, nil
}

func lexGroup(l *Lexer) (StateFn, error) {
	opening := l.runeStk.Pop()
	var closing rune

	switch opening {
	case '[':
		l.emitTokenOpts("[", l.line, l.column, BracketLeft)
		closing = ']'
	case '(':
		l.emitTokenOpts("(", l.line, l.column, ParenLeft)
		closing = ')'
	default:
		return nil, ErrUnexpectedRune
	}

	l.runeStk.Push(closing)

	return lexToken, nil
}

func lexTerminalSymbol(l *Lexer) (StateFn, error) {
	r, err := l.nextChar()
	if err != nil {
		return nil, err
	}

	if unicode.IsLetter(r) || unicode.IsDigit(r) {
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
	} else if !unicode.IsPunct(r) && !unicode.IsSymbol(r) {
		return nil, ErrUnexpectedRune
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

	l.emitToken(ProdRule)

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
