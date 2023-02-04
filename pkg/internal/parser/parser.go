package parser

import (
	"fmt"
	"github.com/francoganga/pagoda_bun/pkg/internal/lexer"
	"strconv"
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

func (p *Parser) parseDate() string {

	time_str := p.curToken.Literal

	if !p.expectPeek(lexer.SLASH) {
		return ""
	}

	time_str += p.curToken.Literal

	if !p.expectPeek(lexer.INT) {
		return ""
	}

	time_str += p.curToken.Literal

	if !p.expectPeek(lexer.SLASH) {
		return ""
	}

	time_str += p.curToken.Literal

	if !p.expectPeek(lexer.INT) {
		return ""
	}

	time_str += p.curToken.Literal

	return time_str
}

type ConsumoDto struct {
	Date        string
	Code        string
	Description string
	Amount      int
	Balance     int
}

func (p *Parser) parseAmount() int {

	str_code := ""

	if p.peekTokenIs(lexer.MINUS) {
		p.nextToken()
		str_code += p.curToken.Literal
	}

	if p.peekTokenIs(lexer.DOLLAR) {
		p.nextToken()
	}

	if !p.expectPeek(lexer.INT) {
		return 0
	}

	str_code += p.curToken.Literal

	if p.peekTokenIs(lexer.DOT) {
		p.nextToken()

		if !p.expectPeek(lexer.INT) {
			return 0
		}

		str_code += p.curToken.Literal
	}

	amount, err := strconv.Atoi(str_code)

	if err != nil {
		return 0
	}

	amount = amount * 100

	if p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		if !p.expectPeek(lexer.INT) {
			return 0
		}

		decimal, err := strconv.Atoi(p.curToken.Literal)

		if err != nil {
			return 0
		}

		if amount < 0 {
			amount -= decimal
		} else {
			amount += decimal
		}

	}

	return amount
}

func (p *Parser) ParseConsumo() *ConsumoDto {
	c := &ConsumoDto{}

	date := p.parseDate()

	c.Date = date

	if !p.expectPeek(lexer.INT) {
		return nil
	}

	c.Code = p.curToken.Literal

	if !p.expectPeek(lexer.DESC) {
		return nil
	}

	c.Description = p.curToken.Literal

	c.Amount = p.parseAmount()

	c.Balance = p.parseAmount()

	if p.peekTokenIs(lexer.DESC) {
		p.nextToken()
		c.Description += ": " + p.curToken.Literal
	}

	return c
}
