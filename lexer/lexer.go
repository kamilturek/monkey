package lexer

import "github.com/kamilturek/monkey/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
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

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{
				Type:    token.EQ,
				Literal: string(ch) + string(l.ch),
			}
		} else {
			tok = token.Token{
				Type:    token.ASSIGN,
				Literal: string(l.ch),
			}
		}
	case '+':
		tok = token.Token{
			Type:    token.PLUS,
			Literal: string(l.ch),
		}
	case '-':
		tok = token.Token{
			Type:    token.MINUS,
			Literal: string(l.ch),
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: string(ch) + string(l.ch),
			}
		} else {
			tok = token.Token{
				Type:    token.BANG,
				Literal: string(l.ch),
			}
		}
	case '*':
		tok = token.Token{
			Type:    token.ASTERISK,
			Literal: string(l.ch),
		}
	case '/':
		tok = token.Token{
			Type:    token.SLASH,
			Literal: string(l.ch),
		}
	case '<':
		tok = token.Token{
			Type:    token.LT,
			Literal: string(l.ch),
		}
	case '>':
		tok = token.Token{
			Type:    token.GT,
			Literal: string(l.ch),
		}
	case '(':
		tok = token.Token{
			Type:    token.LPAREN,
			Literal: string(l.ch),
		}
	case ')':
		tok = token.Token{
			Type:    token.RPAREN,
			Literal: string(l.ch),
		}
	case '{':
		tok = token.Token{
			Type:    token.LBRACE,
			Literal: string(l.ch),
		}
	case '}':
		tok = token.Token{
			Type:    token.RBRACE,
			Literal: string(l.ch),
		}
	case ',':
		tok = token.Token{
			Type:    token.COMMA,
			Literal: string(l.ch),
		}
	case ';':
		tok = token.Token{
			Type:    token.SEMICOLON,
			Literal: string(l.ch),
		}
	case 0:
		tok = token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	default:
		if isLetter(l.ch) {
			identifier := l.readIdentifier()
			tokenType := token.LookupIdent(identifier)

			return token.Token{
				Type:    tokenType,
				Literal: identifier,
			}
		}

		if isDigit(l.ch) {
			return token.Token{
				Type:    token.INT,
				Literal: l.readNumber(),
			}
		}

		tok = token.Token{
			Type:    token.ILLEGAL,
			Literal: string(l.ch),
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) readIdentifier() string {
	startPosition := l.position

	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	return l.input[startPosition:l.position]
}

func (l *Lexer) readNumber() string {
	startPosition := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[startPosition:l.position]
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}
