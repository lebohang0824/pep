package lexer

import "fmt"

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	APP      TokenType = "APP"
	META     TokenType = "META"
	END      TokenType = "END"
	ACTION   TokenType = "ACTION"
	FEATURE  TokenType = "FEATURE"
	PARAMS   TokenType = "PARAMS"
	PARAM    TokenType = "PARAM"
	RULES    TokenType = "RULES"
	RULE     TokenType = "RULE"
	IF       TokenType = "IF"
	EMPTY    TokenType = "EMPTY"
	REJECT   TokenType = "REJECT"
	EVENTS   TokenType = "EVENTS"
	TRIGGER  TokenType = "TRIGGER"
	ENUM     TokenType = "ENUM"
	FUNCTION TokenType = "FUNCTION"
	RETURN   TokenType = "RETURN"
	NOT      TokenType = "NOT"

	IDENT   TokenType = "IDENT"
	STRING  TokenType = "STRING"
	INTEGER TokenType = "INTEGER"
	TRUE    TokenType = "TRUE"
	FALSE   TokenType = "FALSE"

	COLON  TokenType = "COLON"
	LBRACE TokenType = "LBRACE"
	RBRACE TokenType = "RBRACE"
	LPAREN TokenType = "LPAREN"
	RPAREN TokenType = "RPAREN"
	COMMA  TokenType = "COMMA"
	EQ_EQ  TokenType = "EQ_EQ"
	DOT    TokenType = "DOT"
	GTE    TokenType = "GTE"
	GT     TokenType = "GT"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%s, %q, line=%d, col=%d}", t.Type, t.Literal, t.Line, t.Column)
}

type Lexer struct {
	input    string
	position int
	readPos  int
	ch       byte
	line     int
	column   int
}

var keywords = map[string]TokenType{
	"app":     APP,
	"meta":    META,
	"end":     END,
	"action":  ACTION,
	"feature": FEATURE,
	"params":  PARAMS,
	"param":   PARAM,
	"rules":   RULES,
	"rule":    RULE,
	"if":      IF,
	"empty":   EMPTY,
	"reject":  REJECT,
	"events":  EVENTS,
	"trigger":  TRIGGER,
	"enum":     ENUM,
	"function": FUNCTION,
	"return":   RETURN,
	"not":      NOT,
	"true":     TRUE,
	"false":    FALSE,
	"type":     IDENT,
	"required": IDENT,
	"default":  IDENT,
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.position = l.readPos
	l.readPos++
	l.column++
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case 0:
		tok.Type = EOF
		tok.Literal = ""
		tok.Line = l.line
		tok.Column = l.column
	case ':':
		tok = l.makeToken(COLON, ":")
	case '{':
		tok = l.makeToken(LBRACE, "{")
	case '}':
		tok = l.makeToken(RBRACE, "}")
	case '(':
		tok = l.makeToken(LPAREN, "(")
	case ')':
		tok = l.makeToken(RPAREN, ")")
	case ',':
		tok = l.makeToken(COMMA, ",")
	case '.':
		tok = l.makeToken(DOT, ".")
	case '"':
		tok = l.readString()
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(GTE, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(GT, string(l.ch))
		}
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.makeToken(EQ_EQ, string(ch)+string(l.ch))
		} else {
			tok = l.makeToken(ILLEGAL, string(l.ch))
		}
	default:
		if isLetter(l.ch) {
			tok = l.readIdentifier()
			return tok
		} else if isDigit(l.ch) {
			tok = l.readNumber()
			return tok
		} else {
			tok = l.makeToken(ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) makeToken(t TokenType, lit string) Token {
	return Token{Type: t, Literal: lit, Line: l.line, Column: l.column}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() Token {
	startCol := l.column
	startPos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	literal := l.input[startPos:l.position]
	tokType, ok := keywords[literal]
	if !ok {
		tokType = IDENT
	}
	return Token{Type: tokType, Literal: literal, Line: l.line, Column: startCol}
}

func (l *Lexer) readNumber() Token {
	startCol := l.column
	startPos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return Token{Type: INTEGER, Literal: l.input[startPos:l.position], Line: l.line, Column: startCol}
}

func (l *Lexer) readString() Token {
	startCol := l.column
	l.readChar() // skip opening "
	startPos := l.position
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' && l.peekChar() == '"' {
			l.readChar()
		}
		l.readChar()
	}
	literal := l.input[startPos:l.position]
	// leave ch at closing " so NextToken's readChar advances past it
	return Token{Type: STRING, Literal: literal, Line: l.line, Column: startCol}
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
