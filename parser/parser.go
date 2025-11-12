package parser

import (
	"errors"
	"fmt"
)

func Parse(tokens []*Token) (*Node, error) {
	if len(tokens) == 0 {
		return nil, errors.New("empty tokens")
	}
	ps := Parser{
		tokens: tokens,
		pos:    0,
		la:     tokens[0],
	}
	ret := ps.expr()
	if ret == nil || ps.la.Type != TokenTypeEndOfFile {
		tok := tokens[ps.max_]
		return nil, fmt.Errorf("fail to parse: \"%s\" at %d", tok.Text, tok.Start)
	}
	return ret, nil
}

type Parser struct {
	tokens []*Token
	pos    int
	max_   int
	la     *Token
}

func (p *Parser) expr() *Node {
	return p.ternary()
}

func (p *Parser) ternary() *Node {
	pos := p.pos
	if condition := p.binary(); condition != nil {
		if p.expect("?") != nil {
			if lhs := p.binary(); lhs != nil {
				if p.expect(":") != nil {
					if rhs := p.binary(); rhs != nil {
						return NewTernaryNode(condition, lhs, rhs)
					}
				}
			}
		} else {
			return condition
		}
	}
	p.reset(pos)
	return nil
}

func (p *Parser) binary() *Node {
	return p.compareBinary()
}

func (p *Parser) compareBinary() *Node {
	pos := p.pos
	var lhs *Node
	if lhs = p.logicalBinary(); lhs != nil {
		for {
			tmp := p.pos
			if op := p.expectOp("==", "!=", ">=", ">", "<=", "<"); op != nil {
				if rhs := p.logicalBinary(); rhs != nil {
					lhs = NewBinaryNode(op, lhs, rhs)
					continue
				}
			}
			p.reset(tmp)
			break
		}
	}
	if lhs != nil {
		return lhs
	}
	p.reset(pos)
	return nil
}

func (p *Parser) logicalBinary() *Node {
	pos := p.pos
	var lhs *Node
	if lhs = p.sumBinary(); lhs != nil {
		for {
			tmp := p.pos
			if op := p.expectOp("==", "!=", ">=", ">", "<=", "<"); op != nil {
				if rhs := p.sumBinary(); rhs != nil {
					lhs = NewBinaryNode(op, lhs, rhs)
					continue
				}
			}
			p.reset(tmp)
			break
		}
	}
	if lhs != nil {
		return lhs
	}
	p.reset(pos)
	return nil
}

func (p *Parser) sumBinary() *Node {
	pos := p.pos
	var lhs *Node
	if lhs = p.termBinary(); lhs != nil {
		for {
			tmp := p.pos
			if op := p.expectOp("+", "-"); op != nil {
				if rhs := p.termBinary(); rhs != nil {
					lhs = NewBinaryNode(op, lhs, rhs)
					continue
				}
			}
			p.reset(tmp)
			break
		}
	}
	if lhs != nil {
		return lhs
	}
	p.reset(pos)
	return nil
}

func (p *Parser) termBinary() *Node {
	pos := p.pos
	var lhs *Node
	if lhs = p.unary(); lhs != nil {
		for {
			tmp := p.pos
			if op := p.expectOp("*", "/", "%"); op != nil {
				if rhs := p.unary(); rhs != nil {
					lhs = NewBinaryNode(op, lhs, rhs)
					continue
				}
			}
			p.reset(tmp)
			break
		}
	}
	if lhs != nil {
		return lhs
	}
	p.reset(pos)
	return nil
}

func (p *Parser) expectOp(ops ...string) *Token {
	for _, op := range ops {
		if tok := p.expect(op); tok != nil {
			return tok
		}
	}
	return nil
}

func (p *Parser) unary() *Node {
	if op := p.expectOp("-", "!"); op != nil {
		if x := p.primary(); x != nil {
			return NewUnaryNode(op, x)
		}
	}
	return p.primary()
}

func (p *Parser) primary() *Node {
	pos := p.pos
	var lhs *Node
	if lhs = p.atom(); lhs != nil {
		for {
			tmp := p.pos
			if p.expect("(") != nil {
				var args []*Node
				for {
					if arg := p.expr(); arg != nil {
						args = append(args, arg)
					} else {
						break
					}
					if p.expect(",") != nil {
						continue
					} else {
						break
					}
				}
				if p.expect(")") != nil {
					lhs = NewCallNode(lhs, args)
					continue
				}
			}
			p.reset(tmp)
			if p.expect("[") != nil {
				if x := p.expr(); x != nil {
					if p.expect("]") != nil {
						lhs = NewIndexNode(lhs, x)
						continue
					}
				}
			}
			p.reset(tmp)
			if p.expect(".") != nil {
				if x := p.expect(TokenTypeIdent); x != nil {
					lhs = NewSelectorNode(lhs, x)
					continue
				}
			}
			break
		}
	}
	if lhs != nil {
		return lhs
	}
	p.reset(pos)
	if p.expect(".") != nil {
		if x := p.expect(TokenTypeIdent); x != nil {
			return NewSelectorNode(nil, x)
		}
		p.reset(pos)
	}
	return nil
}

func (p *Parser) atom() *Node {
	pos := p.pos
	if tok := p.expect(TokenTypeIdent); tok != nil {
		return NewIdentNode(tok)
	} else if tok = p.expect(TokenTypeNumber); tok != nil {
		return NewNumberNode(tok)
	} else if tok = p.expect(TokenTypeString); tok != nil {
		return NewStringNode(tok)
	} else if p.expect("(") != nil {
		n := p.expr()
		if n == nil {
			p.reset(pos)
			return nil
		}
		if p.expect(")") != nil {
			return NewParenNode(n)
		}
	} else if p.expect("[") != nil {
		var items []*Node
		n := p.expr()
		for n != nil {
			items = append(items, n)
			if p.expect(",") == nil {
				break
			}
			n = p.expr()
		}
		if p.expect("]") != nil {
			return NewArrayNode(items)
		}
		p.reset(pos)
		return nil
	} else if p.expect("{") != nil {
		var items []*Node
		n := p.pair()
		for n != nil {
			items = append(items, n)
			if p.expect(",") == nil {
				break
			}
			n = p.pair()
		}
		if p.expect("}") != nil {
			return NewObjectNode(items)
		}
		p.reset(pos)
		return nil
	}
	return nil
}

func (p *Parser) pair() *Node {
	pos := p.pos
	if k := p.expect(TokenTypeIdent); k != nil {
		if p.expect(":") != nil {
			if v := p.expr(); v != nil {
				return NewPairNode(k, v)
			}
		}
	}
	p.reset(pos)
	return nil
}

func (p *Parser) reset(pos int) {
	p.pos = pos
	p.read()
}

func (p *Parser) read() {
	if p.pos < len(p.tokens) {
		p.la = p.tokens[p.pos]
	}
}

func (p *Parser) forward() {
	if p.pos > p.max_ {
		p.max_ = p.pos
	}
	p.pos++
	p.read()
}

func (p *Parser) expect(t string) *Token {
	if p.la.Type == t || p.la.Text == t {
		ret := p.la
		p.forward()
		return ret
	}
	return nil
}
