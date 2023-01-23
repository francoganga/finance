package parser

import (
	"fmt"
	"github.com/francoganga/pagoda_bun/pkg/internal/lexer"
	"testing"

	"github.com/ztrue/tracerr"
)

func TestParseDate(t *testing.T) {

	someDate := "05/08/22"

	input := fmt.Sprintf("%s           99999999", someDate)

	l := lexer.New(input)

	p := New(l)

	d, err := p.parseDate()

	if err != nil {
		tracerr.PrintSourceColor(err)
	}

	checkParserErrors(t, p)

	if d != someDate {
		t.Fatalf("Error parsing date, expected=%s, got=%s", someDate, d)
	}
}

func TestParseAmount(t *testing.T) {
	input := "$ 104.095,74"

	expectedAmount := 10409574

	l := lexer.New(input)

	p := New(l)

	a, err := p.parseAmount()

	if err != nil {
		tracerr.PrintSourceColor(err)
	}

	checkParserErrors(t, p)

	if a != expectedAmount {
		t.Fatalf("Error parsing amount, expected=%d, got=%d", expectedAmount, a)
	}

}

func TestParseConsumo(t *testing.T) {

	input := `05/07/21               10280171       Compra con tarjeta de debito                               -$ 650,00                                $ 104.095,74
    Mercadopago*recargatuenti - tarj nro. 1866`

	l := lexer.New(input)

	p := New(l)

	c, err := p.ParseConsumo()

	if err != nil {
		tracerr.PrintSourceColor(err)
		t.FailNow()
	}
	checkParserErrors(t, p)

	if c.Date != "05/07/21" {
		t.Fatalf("expected string to be=%s, got=%s", "05/07/2021", c.Date)
	}

	if c.Code != "10280171" {
		t.Fatalf("expected code to be=%s, got=%s", "10280171", c.Code)
	}

	if c.Description != "Compra con tarjeta de debito: Mercadopago*recargatuenti - tarj nro. 1866" {
		t.Fatalf("expected Description to be=%s, got=%s", "Compra con tarjeta de debito Mercadopago*recargatuenti - tarj nro. 1866", c.Description)
	}

	if c.Amount != -65000 {
		t.Fatalf("expected Amount to be=%d, got=%d", 6500, c.Amount)
	}

	if c.Balance != 10409574 {
		t.Fatalf("expected Balance to be=%d, got=%d", 10409574, c.Balance)
	}

}

func TestParseConsumo2(t *testing.T) {

	input := `25/08/21           25593863      Transferencia realizada                                             -$ 6.000,00                                $ 96.424,39
    A ganga carlos ignacio / varios - var / 201645877712`

	l := lexer.New(input)

	p := New(l)

	c, err := p.ParseConsumo()

	if err != nil {
		tracerr.PrintSourceColor(err)
		t.FailNow()
	}
	checkParserErrors(t, p)

	fmt.Printf("c=%v\n", c)

}

func TestParseConsumo3(t *testing.T) {
	input := `06/08/21               79378436       Acreditacion de haberes                                        $ 105.319,15                                 $ 185.136,86
    307113401661 210805007universidad nacional a jauretc`

	l := lexer.New(input)

	p := New(l)

	c, err := p.ParseConsumo()

	if err != nil {
		tracerr.PrintSourceColor(err)
		t.FailNow()
	}
	checkParserErrors(t, p)

	fmt.Printf("c=%v\n", c)
}

func TestParseConsumo4(t *testing.T) {

	input := `02/12/22 1714314                      Compra con tarjeta de debito                                       -$ 548,00                                $ 166.696,92
    Autoservicio santa ana - tarj nro. 1866`

	p := FromInput(input)

	c, err := p.ParseConsumo()

	if err != nil {
		tracerr.PrintSourceColor(err)
		t.FailNow()
	}
	checkParserErrors(t, p)

	fmt.Printf("c=%v\n", c)
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
