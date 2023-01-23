package lexer

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	INT    = "INT"
	SLASH  = "SLASH"
	DOLLAR = "DOLLAR"
	MINUS  = "MINUS"
	COMMA  = "COMMA"
	DOT    = "DOT"
	WORD   = "WORD"
	STAR   = "*"
	SEP    = "SEP"
	EOF    = "EOF"
)

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) NextToken() Token {
	var tok Token

	if l.hasMoreThanOneSpace() {
		l.skipWhitespace()
		return newToken(SEP, byte('_'))
	}

	if l.hasNewLine() {
		l.skipWhitespace()
		return newToken(SEP, byte('_'))
	}

	l.skipWhitespace()

	switch l.ch {
	case '-':
		tok = newToken(MINUS, l.ch)
	case '/':
		tok = newToken(SLASH, l.ch)
	case '$':
		tok = newToken(DOLLAR, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '.':
		tok = newToken(DOT, l.ch)
	case '*':
		tok = newToken(STAR, l.ch)

	case 0:
		tok.Literal = ""
		tok.Type = EOF

	default:

		if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		} else if isLetter(l.ch) {
			tok.Type = WORD
			tok.Literal = l.readWord()
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) hasMoreThanOneSpace() bool {
	if l.readPosition >= len(l.input) {
		return false
	}

	if l.position >= len(l.input) {
		return false
	}

	return l.input[l.position] == ' ' && l.input[l.readPosition] == ' '
}

func (l *Lexer) hasNewLine() bool {
	if l.position >= len(l.input) {
		return false
	}

	return l.input[l.position] == '\n'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isSymbol(ch byte) bool {

	if ch == 0 {
		return false
	}

	if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
		return false
	}

	return 0 <= ch && ch <= 127
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}

}

func (l *Lexer) peekTill(till byte, fun func(ch byte) bool) bool {
	cv := l.position
	nv := l.position + 1

	for l.input[cv] != till {

		if fun(l.input[cv]) {
			return true
		}

		if nv >= len(l.input) {
			break
		}

		cv = nv
		nv += 1
	}

	return false
}

func (l *Lexer) readWord() string {
	position := l.position

	for isSymbol(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}
