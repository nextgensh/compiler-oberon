package compiler

import (
	"github.com/proebsting/cs553s2013/lexer"
)

type lexType struct {
	stream  chan lexer.Token
	current lexer.Token
}

func (o *lexType) Peek() (val int) {
	return o.current.Enum()
}

func (o *lexType) Tok() lexer.Token {
	return o.current
}

func (o *lexType) Tok_value() string {
	return o.current.Value()
}

func (o *lexType) Advance() {
	o.current = <-o.stream
}

func (o *lexType) ErrorMessage(s string) (val string) {
	return s
}
