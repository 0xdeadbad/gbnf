package lexer

import (
	"io"
)

type Lexer struct {
	*CharReader
	Tokens        []*Token
	runeTmpBuffer []rune
	currentState  StateFn
	tokenChan     chan *Token
	line          uint
	column        uint
}

func NewLexer(r io.ReadSeeker) *Lexer {
	return &Lexer{
		CharReader:    newCharReader(r),
		Tokens:        make([]*Token, 0),
		currentState:  lexToken,
		runeTmpBuffer: make([]rune, 0),
		tokenChan:     make(chan *Token, 2),
		line:          0,
		column:        0,
	}
}

func (l *Lexer) NextToken() (*Token, error) {
	var err error = nil
	var state StateFn
	for {
		select {
		case token := <-l.tokenChan:
			return token, nil
		default:
			state, err = l.currentState(l)
			if err != nil {
				return nil, err
			}
			l.clearRuneTmpBuffer()
			l.currentState = state
		}
	}
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
	token := NewToken(lexeme, line, column, typ)
	l.Tokens = append(l.Tokens, token)

	l.tokenChan <- token
}

func (l *Lexer) emitToken(typ TokenType) {
	l.emitTokenOpts(string(l.runeTmpBuffer), l.line, l.column, typ)
}

func (l *Lexer) clearRuneTmpBuffer() {
	l.runeTmpBuffer = l.runeTmpBuffer[:0]
}
