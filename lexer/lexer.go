package lexer

import (
	"io"
)

type Lexer struct {
	*CharReader
	Tokens        []*Token
	runeTmpBuffer []rune
	runeStk       *Stk[rune]
	stateStk      *Stk[StateFn]
	currentState  StateFn
	tokenChan     chan *Token
	pointer       int
	line          uint
	column        uint
}

func NewLexer(r io.ReadSeeker) *Lexer {
	return &Lexer{
		CharReader:    newCharReader(r),
		Tokens:        make([]*Token, 0),
		currentState:  lexToken,
		runeTmpBuffer: make([]rune, 0),
		tokenChan:     make(chan *Token, 3),
		runeStk:       NewStk[rune](),
		stateStk:      NewStk[StateFn](),
		pointer:       0,
		line:          0,
		column:        0,
	}
}

func (l *Lexer) NextToken() (*Token, error) {
	var err error
	var state StateFn
	for {
		select {
		case token, ok := <-l.tokenChan:
			if !ok && token == nil {
				return nil, nil
			}
			return token, nil
		default:
			if !l.stateStk.Empty() {
				state, err = l.stateStk.Dequeue()(l)
			} else {
				state, err = l.currentState(l)
			}
			if err != nil {
				if err == io.EOF {
					return nil, nil
				}
				return nil, err
			}
			if l.pointer < len(l.Tokens) {
				l.tokenChan <- l.Tokens[l.pointer]
				l.pointer++
			}
			l.clearRuneTmpBuffer()
			l.currentState = state
			if l.currentState == nil && l.stateStk.Empty() {
				l.tokenChan <- NewToken("", l.line, l.column, EndMark)
				close(l.tokenChan)
			}
		}
	}
}

func (l *Lexer) PeekToken() (*Token, error) {
	mark := l.Mark()
	defer l.Reset(mark)

	return l.NextToken()
}

func (l *Lexer) Mark() int {
	return l.pointer
}

func (l *Lexer) Reset(mark int) {
	l.pointer = mark
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
}

func (l *Lexer) emitToken(typ TokenType) {
	l.emitTokenOpts(string(l.runeTmpBuffer), l.line, l.column, typ)
}

func (l *Lexer) clearRuneTmpBuffer() {
	l.runeTmpBuffer = l.runeTmpBuffer[:0]
}

type Stk[T any] struct {
	stk []T
}

func NewStk[T any]() *Stk[T] {
	return &Stk[T]{stk: make([]T, 0)}
}

func (s *Stk[T]) Push(v T) {
	s.stk = append(s.stk, v)
}

func (s *Stk[T]) Pop() T {
	var v T
	if s.Empty() {
		return v
	}

	v = s.stk[len(s.stk)-1]
	s.stk = s.stk[:len(s.stk)-1]

	return v
}

func (s *Stk[T]) Peek() T {
	var v T
	if s.Empty() {
		return v
	}

	v = s.stk[len(s.stk)-1]

	return v
}

func (s *Stk[T]) Len() int {
	return len(s.stk)
}

func (s *Stk[T]) Empty() bool {
	return len(s.stk) == 0
}

func (s *Stk[T]) Clear() {
	s.stk = s.stk[:0]
}

func (s *Stk[T]) Enqueue(v T) {
	s.Push(v)
}

func (s *Stk[T]) Dequeue() T {
	var v T
	if s.Empty() {
		return v
	}

	v = s.stk[0]
	s.stk = s.stk[1:]

	return v
}

func (s *Stk[T]) Front() T {
	var v T
	if s.Empty() {
		return v
	}

	v = s.stk[0]

	return v
}

func (s *Stk[T]) Back() T {
	return s.Peek()
}
