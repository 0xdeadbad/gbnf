package lexer

import (
	"io"
)

type Lexer struct {
	*CharReader
	Tokens         []*Token
	tokenTmpBuffer []*Token
	runeTmpBuffer  []rune
	line           uint
	column         uint
}

func NewLexer(r io.ReadSeeker) *Lexer {
	return &Lexer{
		CharReader:     newCharReader(r),
		Tokens:         make([]*Token, 0),
		tokenTmpBuffer: make([]*Token, 0),
		runeTmpBuffer:  make([]rune, 0),
		line:           0,
		column:         0,
	}
}

func (l *Lexer) NextToken() error {

	_, err := lexToken(l)
	if err != nil {
		return err
	}

	return nil
}

func (l *Lexer) nextChar() (rune, error) {
	defer func(l *Lexer) {
		if l.buffer[0] == '\n' {
			l.line++
			l.column = 0
		} else {
			l.column++
		}
		l.runeTmpBuffer = append(l.runeTmpBuffer, l.buffer[0])
	}(l)

	return l.CharReader.nextChar()
}

func (l *Lexer) peekChar() (rune, error) {
	return l.CharReader.peekChar()
}

func (l *Lexer) emitTokenOpts(lexeme string, line, column uint, typ TokenType) {
	l.Tokens = append(l.Tokens, NewToken(lexeme, line, column, typ))
}

func (l *Lexer) emitToken(typ TokenType) {
	l.emitTokenOpts(string(l.runeTmpBuffer), l.line, l.column, typ)
}

func (l *Lexer) emitTmpTokenOpts(lexeme string, line, column uint, typ TokenType) {
	l.Tokens = append(l.tokenTmpBuffer, NewToken(lexeme, line, column, typ))
}

func (l *Lexer) emitTmpToken(typ TokenType) {
	l.emitTmpTokenOpts(string(l.runeTmpBuffer), l.line, l.column, typ)
}

func (l *Lexer) clearRuneTmpBuffer() {
	l.runeTmpBuffer = l.runeTmpBuffer[:0]
}

func (l *Lexer) clearTokenTmpBuffer() {
	l.tokenTmpBuffer = l.tokenTmpBuffer[:0]
}
