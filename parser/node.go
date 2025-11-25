package parser

const NodeTypeIdent = "ident"
const NodeTypeNumber = "number"
const NodeTypeString = "string"
const NodeTypeUnary = "unary"
const NodeTypeBinary = "binary"
const NodeTypeCall = "call"
const NodeTypeSelector = "selector"
const NodeTypeParen = "paren"

func NewIdentNode(token *Token) *Node {
	return &Node{type_: NodeTypeIdent, token: token}
}

func NewNumberNode(token *Token) *Node {
	return &Node{type_: NodeTypeNumber, token: token}
}

func NewStringNode(token *Token) *Node {
	return &Node{type_: NodeTypeString, token: token}
}

func NewUnaryNode(op *Token, target *Node) *Node {
	return &Node{type_: NodeTypeUnary, op: op, x: target}
}

func NewBinaryNode(op *Token, lhs, rhs *Node) *Node {
	return &Node{type_: NodeTypeBinary, op: op, x: lhs, y: rhs}
}

func NewCallNode(callee *Node, args []*Node) *Node {
	return &Node{type_: NodeTypeCall, x: callee, s: args}
}

func NewSelectorNode(target *Node, key *Token) *Node {
	return &Node{type_: NodeTypeSelector, x: target, token: key}
}

func NewParenNode(n *Node) *Node {
	return &Node{type_: NodeTypeParen, x: n}
}

type Node struct {
	type_ string
	token *Token  // ident, number, string
	op    *Token  // unary, binary
	x     *Node   // unary, binary lhs, call callee
	y     *Node   // binary rhs
	s     []*Node // call args
}

func (n *Node) Type() string {
	return n.type_
}

func (n *Node) UnaryTarget() *Node {
	return n.x
}

func (n *Node) BinaryLhs() *Node {
	return n.x
}

func (n *Node) BinaryRhs() *Node {
	return n.y
}

func (n *Node) Callee() *Node {
	return n.x
}

func (n *Node) Args() []*Node {
	return n.s
}

func (n *Node) SelectorTarget() *Node {
	return n.x
}

func (n *Node) SelectorKey() string {
	return n.token.Text
}

func (n *Node) ParenTarget() *Node {
	return n.x
}

func (n *Node) Ident() string {
	return n.token.Text
}

func (n *Node) Number() string {
	return n.token.Text
}

func (n *Node) String() string {
	if n.token == nil {
		return ""
	}
	return n.token.Text
}

func (n *Node) Op() string {
	return n.op.Text
}

func (n *Node) Visit(f func(node *Node)) {
	f(n)
	if n.x != nil {
		n.x.Visit(f)
	}
	if n.y != nil {
		n.y.Visit(f)
	}
	for _, t := range n.s {
		t.Visit(f)
	}
}
