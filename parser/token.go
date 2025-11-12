package parser

const TokenTypeEndOfFile = "end_of_file"
const TokenTypeWhitespace = "whitespace"
const TokenTypeIdent = "ident"
const TokenTypeNumber = "number"
const TokenTypeString = "string"
const TokenTypeOpComma = ","
const TokenTypeOpDot = "."

// compare binary

const TokenTypeOpEqualEqual = "=="
const TokenTypeOpNotEqual = "!="
const TokenTypeOpGreaterEqual = ">="
const TokenTypeOpLessEqual = "<="
const TokenTypeOpGreater = ">"
const TokenTypeOpLess = "<"

// logical binary

const TokenTypeOpAndAnd = "&&"
const TokenTypeOpBarBar = "||"

// sum binary

const TokenTypeOpMinus = "-"
const TokenTypeOpPlus = "+"

// term binary

const TokenTypeOpSlash = "/"
const TokenTypeOpStar = "*"
const TokenTypeOpPercent = "%"

// unary

const TokenTypeOpNot = "!"

const TokenTypeOpLeftBrace = "{"
const TokenTypeOpRightBrace = "}"
const TokenTypeOpLeftBracket = "["
const TokenTypeOpRightBracket = "]"
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
