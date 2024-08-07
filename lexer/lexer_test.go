package lexer_test

import (
	"testing"

	"github.com/kamilturek/monkey/lexer"
	"github.com/kamilturek/monkey/token"
)

func TestNextToken(t *testing.T) {
	t.Parallel()

	input := `
	let five = 5;
	let ten = 10;

	let add = fn(x, y) {
		x + y;
	};

	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;

	if (5 < 10) {
		return true;
	} else {
		return false;
	}

	10 == 10;
	10 != 9;

	return5;

	"foobar";
	"foo bar";
	`

	type expectedToken struct {
		expectedType    token.TokenType
		expectedLiteral string
	}

	// Pre-defined expected tokens
	LET := expectedToken{token.LET, "let"}
	INT5 := expectedToken{token.INT, "5"}
	INT9 := expectedToken{token.INT, "9"}
	INT10 := expectedToken{token.INT, "10"}
	SEMICOLON := expectedToken{token.SEMICOLON, ";"}

	tests := []expectedToken{
		LET,
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		INT5,
		SEMICOLON,
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		INT10,
		SEMICOLON,
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		SEMICOLON,
		{token.RBRACE, "}"},
		SEMICOLON,
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		SEMICOLON,
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		INT5,
		SEMICOLON,
		INT5,
		{token.LT, "<"},
		INT10,
		{token.GT, ">"},
		INT5,
		SEMICOLON,
		{token.IF, "if"},
		{token.LPAREN, "("},
		INT5,
		{token.LT, "<"},
		INT10,
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		SEMICOLON,
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		SEMICOLON,
		{token.RBRACE, "}"},
		INT10,
		{token.EQ, "=="},
		INT10,
		SEMICOLON,
		INT10,
		{token.NOT_EQ, "!="},
		INT9,
		SEMICOLON,
		{token.IDENT, "return5"},
		SEMICOLON,
		{token.STRING, "foobar"},
		SEMICOLON,
		{token.STRING, "foo bar"},
		SEMICOLON,
		{token.EOF, ""},
	}

	l := lexer.NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d]- tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d]- literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
