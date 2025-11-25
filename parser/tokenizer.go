package parser

import (
	"errors"
	"fmt"
)

func Tokenize(text string) ([]*Token, error) {
	if text == "" {
		return nil, errors.New("empty text")
	}
	tokenizer := &Tokenizer{text: text, la: text[0]}
	return tokenizer.Parse()
}

type Tokenizer struct {
	text string
	la   byte
	pos  int
}

func (t *Tokenizer) Parse() ([]*Token, error) {
	var ret []*Token
	for {
		tok, err := t.next()
		if err != nil {
			return nil, err
		}
		if tok.Type == TokenTypeWhitespace {
			continue
		}
		ret = append(ret, tok)
		if tok.Type == TokenTypeEndOfFile {
			break
		}
	}
	return ret, nil
}

func (t *Tokenizer) next() (*Token, error) {
	if t.la == 0 {
		return NewToken(TokenTypeEndOfFile, "EOF", t.pos, t.pos), nil
	} else if tok := t.op(); tok != nil {
		return tok, nil
	} else if tok = t.whitespace(); tok != nil {
		return tok, nil
	} else if tok = t.ident(); tok != nil {
		return tok, nil
	} else if tok = t.number(); tok != nil {
		return tok, nil
	} else if tok = t.string(); tok != nil {
		return tok, nil
	}
	return nil, fmt.Errorf("fail to tokenize '%s' at %d: \"%s\"", string(t.la), t.pos, t.text)
}

func (t *Tokenizer) isLetter(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}

func (t *Tokenizer) isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (t *Tokenizer) ident() *Token {
	if t.isLetter(t.la) {
		start := t.pos
		t.forward()
		for t.isLetter(t.la) || t.isDigit(t.la) {
			t.forward()
		}
		return t.newToken(TokenTypeIdent, start)
	}
	return nil
}

func (t *Tokenizer) number() *Token {
	if t.isDigit(t.la) {
		start := t.pos
		t.forward()
		for t.isDigit(t.la) {
			t.forward()
		}
		if t.la == '.' {
			t.forward()
			for t.isDigit(t.la) {
				t.forward()
			}
		}
		return t.newToken(TokenTypeNumber, start)
	}
	return nil
}

func (t *Tokenizer) string() *Token {
	if t.la == '\'' {
		start := t.pos
		t.forward()
		for {
			if t.la == '\\' {
				t.forward()
				t.forward()
			} else {
				t.forward()
				if t.la == '\'' || t.la == 0 {
					break
				}
			}
		}
		if t.la == '\'' {
			t.forward()
		}
		return t.newToken(TokenTypeString, start)
	}
	return nil
}

func (t *Tokenizer) whitespace() *Token {
	start := t.pos
	for t.la == ' ' {
		t.forward()
	}
	if start != t.pos {
		return t.newToken(TokenTypeWhitespace, start)
	}
	return nil
}

func (t *Tokenizer) op() *Token {
	var type_ string
	start := t.pos
	switch t.la {
	case '=':
		t.forward()
		if t.la == '=' {
			t.forward()
			type_ = TokenTypeOpEqualEqual
		}
	case '<':
		t.forward()
		if t.la == '=' {
			t.forward()
			type_ = TokenTypeOpLessEqual
		} else {
			type_ = TokenTypeOpLess
		}
	case '>':
		t.forward()
		if t.la == '=' {
			t.forward()
			type_ = TokenTypeOpGreaterEqual
		} else {
			type_ = TokenTypeOpGreater
		}
	case '!':
		t.forward()
		if t.la == '=' {
			t.forward()
			type_ = TokenTypeOpNotEqual
		}
	case '(':
		t.forward()
		type_ = TokenTypeOpLeftParen
	case ')':
		t.forward()
		type_ = TokenTypeOpRightParen
	case '.':
		t.forward()
		type_ = TokenTypeOpDot
	}
	if type_ != "" {
		return t.newToken(type_, start)
	}
	return nil
}

func (t *Tokenizer) newToken(type_ string, start int) *Token {
	return NewToken(type_, t.text[start:t.pos], start, t.pos)
}

func (t *Tokenizer) forward() {
	if t.pos < len(t.text) {
		t.pos++
		t.read()
	}
}

func (t *Tokenizer) read() {
	if t.pos >= len(t.text) {
		t.la = 0
	} else {
		t.la = t.text[t.pos]
	}
}
