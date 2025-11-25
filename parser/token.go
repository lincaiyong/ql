package parser

const TokenTypeEndOfFile = "end_of_file"
const TokenTypeWhitespace = "whitespace"
const TokenTypeIdent = "ident"
const TokenTypeNumber = "number"
const TokenTypeString = "string"
const TokenTypeOpDot = "."
const TokenTypeOpEqualEqual = "=="
const TokenTypeOpNotEqual = "!="
const TokenTypeOpGreaterEqual = ">="
const TokenTypeOpLessEqual = "<="
const TokenTypeOpGreater = ">"
const TokenTypeOpLess = "<"
const TokenTypeOpLeftParen = "("
const TokenTypeOpRightParen = ")"

func NewToken(type_, text string, start, end int) *Token {
	return &Token{type_, text, start, end}
}

type Token struct {
	Type  string
	Text  string
	Start int
	End   int
}
