package lexer

import (
	"fmt"
	"testing"
)

func TestNextToken(t *testing.T) {

	input := `05/07/21               10280171       Compra con tarjeta de debito                               -$ 650,00                                $ 104.095,74
    Mercadopago*recargatuenti - tarj nro. 1866`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{INT, "05"},
		{SLASH, "/"},
		{INT, "07"},
		{SLASH, "/"},
		{INT, "21"},
		{INT, "10280171"},
		{DESC, "Compra con tarjeta de debito"},
		{MINUS, "-"},
		{DOLLAR, "$"},
		{INT, "650"},
		{COMMA, ","},
		{INT, "00"},
		{DOLLAR, "$"},
		{INT, "104"},
		{DOT, "."},
		{INT, "095"},
		{COMMA, ","},
		{INT, "74"},
		{DESC, "Mercadopago*recargatuenti - tarj nro. 1866"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		fmt.Printf("tok=%v\n", tok)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] [position: %d] - tokentype wrong. expected=%q, got=%q", i, l.position, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}

// func TestReadSentence(t *testing.T) {
//
// 	input := `05/07/21               Compra con tarjeta de debito                               -$ 650,00                                $ 104.095,74
//     Mercadopago*recargatuenti - tarj nro. 1866`
//
// 	l := New(input)
//
// 	tok := l.NextToken()
//
// 	fmt.Printf("tok=%v\n", tok)
// 	if tok.Type != DATE {
// 		t.Fatalf("tokentype wrong. expected=%q, got=%q", DATE, tok.Type)
// 	}
//
// }
	

	

