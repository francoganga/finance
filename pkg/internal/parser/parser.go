package parser

import (
	"fmt"
	"github.com/francoganga/pagoda_bun/pkg/internal/lexer"
	"strconv"
	"strings"

	"github.com/ztrue/tracerr"
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  lexer.Token
	peekToken lexer.Token
}

func FromInput(input string) *Parser {
	return New(lexer.New(input))
}

func New(l *lexer.Lexer) *Parser {

	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) expectPeek2(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseDate() (string, error) {

	time_str := p.curToken.Literal

	if !p.expectPeek2(lexer.SLASH) {
		return "", tracerr.Errorf("expected next token to be %s, got %s instead", lexer.SLASH, p.peekToken.Type)
	}

	time_str += p.curToken.Literal

	if !p.expectPeek2(lexer.INT) {
		return "", tracerr.Errorf("expected next token to be %s, got %s instead", lexer.INT, p.peekToken.Type)
	}

	time_str += p.curToken.Literal

	if !p.expectPeek2(lexer.SLASH) {
		return "", tracerr.Errorf("expected next token to be %s, got %s instead", lexer.SLASH, p.peekToken.Type)
	}

	time_str += p.curToken.Literal

	if !p.expectPeek2(lexer.INT) {
		return "", tracerr.Errorf("expected next token to be %s, got %s instead", lexer.INT, p.peekToken.Type)
	}

	time_str += p.curToken.Literal

	return time_str, nil
}

type ConsumoDto struct {
	Date        string `json:"date"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	Balance     int    `json:"balance"`
}

// TODO refactor removing dots and commas for simplicity
func (p *Parser) parseAmount() (int, error) {

	str_code := ""

	if p.peekTokenIs(lexer.MINUS) {
		p.nextToken()
		str_code += p.curToken.Literal
	}

	if p.peekTokenIs(lexer.DOLLAR) {
		p.nextToken()
	}

	if !p.expectPeek2(lexer.INT) {
		return 0, tracerr.Errorf("expected next token to be=%s, got=%s", lexer.INT, p.peekToken.Type)
	}

	str_code += p.curToken.Literal

	if p.peekTokenIs(lexer.DOT) {
		p.nextToken()

		if !p.expectPeek(lexer.INT) {
			return 0, tracerr.Errorf("expected next token to be=%s, got=%s", lexer.INT, p.peekToken.Type)
		}

		str_code += p.curToken.Literal
	}

	amount, err := strconv.Atoi(str_code)

	if err != nil {
		return 0, tracerr.Errorf("Could not parse string=%s to int", str_code)
	}

	amount = amount * 100

	if p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		if !p.expectPeek2(lexer.INT) {
			return 0, tracerr.Errorf("expected next token to be=%s, got=%s", lexer.INT, p.peekToken.Type)
		}

		decimal, err := strconv.Atoi(p.curToken.Literal)

		if err != nil {
			return 0, tracerr.Errorf("Could not parse string=%s to int", p.curToken.Literal)
		}

		if amount < 0 {
			amount -= decimal
		} else {
			amount += decimal
		}

	}

	return amount, nil
}

func (p *Parser) parseDescription() (string, error) {

	desc := ""

	for p.curToken.Type != lexer.SEP && p.curToken.Type != lexer.EOF {
		fmt.Printf("%+v != lexer.SEP && %+v != lexer.EOF \n", p.curToken.Type, p.curToken.Type)
		fmt.Printf("curToek=%v\n", p.curToken)

		fmt.Printf("nextToken=%v\n", p.peekToken)

		desc += p.curToken.Literal + " "

		fmt.Println("after concat??")
		p.nextToken()

		fmt.Println("hace el nextToken????")
	}

	fmt.Println("after for")

	return strings.TrimRight(desc, " "), nil
}

func (p *Parser) ParseConsumo() (*ConsumoDto, error) {
	c := &ConsumoDto{}

	date, err := p.parseDate()

	if err != nil {
		return &ConsumoDto{}, err
	}

	fmt.Println("after date parse")

	c.Date = date

	if !p.expectPeek2(lexer.SEP) {
		return &ConsumoDto{}, tracerr.Errorf("expected next token to be %s, got %s instead", lexer.SEP, p.peekToken.Type)
	}

	p.nextToken()

	c.Code = p.curToken.Literal

	fmt.Println("after code")

	if !p.expectPeek2(lexer.SEP) {
		return &ConsumoDto{}, tracerr.Errorf("expected next token to be %s, got %s instead", lexer.SEP, p.peekToken.Type)
	}

	p.nextToken()

	desc, err := p.parseDescription()

	if err != nil {
		return &ConsumoDto{}, err
	}

	fmt.Println("after desc1")

	c.Description = desc

	a, err := p.parseAmount()

	if err != nil {
		return &ConsumoDto{}, err
	}

	fmt.Println("after amount")

	c.Amount = a

	if !p.expectPeek2(lexer.SEP) {
		return &ConsumoDto{}, tracerr.Errorf("expected next token to be %s, got %s instead", lexer.SEP, p.peekToken.Type)
	}

	p.nextToken()

	b, err := p.parseAmount()

	if err != nil {
		return &ConsumoDto{}, err
	}

	fmt.Println("after balance")

	c.Balance = b

	if !p.expectPeek2(lexer.SEP) {
		return &ConsumoDto{}, tracerr.Errorf("expected next token to be %s, got %s instead", lexer.SEP, p.peekToken.Type)
	}

	p.nextToken()

	desc2, err := p.parseDescription()

	if err != nil {
		return &ConsumoDto{}, tracerr.Errorf("expected next token to be %s, got %s instead", lexer.SEP, p.peekToken.Type)
	}

	fmt.Println("after desc2")

	c.Description += ": " + desc2

	// if p.peekTokenIs(lexer.DESC) {
	// 	p.nextToken()
	// 	c.Description += ": " + p.curToken.Literal
	// }

	return c, nil
}
