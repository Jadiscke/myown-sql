package parser

import (
	"fmt"

	"github.com/Jadiscke/myown-sql/internal/lexer"
)

func tokenFromKeyword(k lexer.Keyword) lexer.Token {
	return lexer.Token{
		Kind:  lexer.KeywordKind,
		Value: string(k),
	}
}

func tokenFromSymbol(s lexer.Symbol) lexer.Token {
	return lexer.Token{
		Kind:  lexer.SymbolKind,
		Value: string(s),
	}
}

func expectToken(tokens []*lexer.Token, cursor uint, t lexer.Token) bool {
	if cursor >= uint(len(tokens)) {
		return false
	}

	return t.Equals(tokens[cursor])
}

func helpMessage(tokens []*lexer.Token, cursor uint, msg string) {
	var c *lexer.Token

	if cursor < uint(len(tokens)) {
		c = tokens[cursor]
	} else {
		c = tokens[cursor-1]
	}

	fmt.Printf("[%d,%d]: %s, got: %s\n", c.Loc.Line, c.Loc.Col, msg, c.Value)

}

func parseToken(tokens []*lexer.Token, initialCursor uint, kind lexer.TokenKind) (*lexer.Token, uint, bool) {
	cursor := initialCursor

	if cursor >= uint(len(tokens)) {
		return nil, cursor, false
	}

	current := tokens[cursor]

	if current.Kind == kind {
		return current, cursor + 1, true
	}

	return nil, initialCursor, false
}

func parseExpression(tokens []*lexer.Token, initialCursor uint, _ lexer.Token) (*expression, uint, bool) {
	cursor := initialCursor

	kinds := []lexer.TokenKind{
		lexer.IdentifierKind,
		lexer.NumericKind,
		lexer.StringKind,
	}

	for _, kind := range kinds {
		t, newCursor, ok := parseToken(tokens, cursor, kind)

		if ok {
			return &expression{
				Literal: t,
				Kind:    LiteralKind,
			}, newCursor, true
		}
	}

	return nil, initialCursor, false
}

func parseExpressions(tokens []*lexer.Token, initialCursor uint, delimiters []lexer.Token) (*[]*expression, uint, bool) {
	cursor := initialCursor

	exps := []*expression{}

outer:
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		current := tokens[cursor]

		for _, delimiter := range delimiters {
			if delimiter.Equals(current) {
				break outer
			}
		}

		if len(exps) > 0 {
			if !expectToken(tokens, cursor, tokenFromSymbol(lexer.CommaSymbol)) {
				helpMessage(tokens, cursor, "Expected comma")

				return nil, initialCursor, false
			}

			cursor++
		}

		exp, newCursor, ok := parseExpression(tokens, cursor, tokenFromSymbol(lexer.CommaSymbol))

		if !ok {
			helpMessage(tokens, cursor, "Expected expression")
			return nil, initialCursor, false
		}

		cursor = newCursor

		exps = append(exps, exp)
	}

	return &exps, cursor, true

}

func Parse(source string) (*Ast, error) {
	tokens, err := lexer.Lex(source)

	if err != nil {
		return nil, err
	}

	a := Ast{}

	cursor := uint(0)

	for cursor < uint(len(tokens)) {
		statement, newCursor, ok := parseStatement(tokens, cursor)
	}
}

func parseStatement(tokens []*lexer.Token, initialCursor uint) (*Statement, uint, bool) {
	cursor := initialCursor

	semicolonToken := tokenFromSymbol(lexer.SemicolonSymbol)

	slct, newCursor, ok := parseSelectStatement(tokens, cursor, semicolonToken)

	if ok {
		return &Statement{
			Kind:            SelectKind,
			SelectStatement: slct,
		}, newCursor, true
	}

	inst, newCursor, ok := parseInsertStatement(tokens, cursor, semicolonToken)

	if ok {
		return &Statement{
			Kind:            InsertKind,
			InsertStatement: inst,
		}, newCursor, true
	}

	createTbl, newCursor, ok := parseCreateTableStatement(tokens, cursor, semicolonToken)

	if ok {
		return &Statement{
			Kind:                 CreateTableKind,
			CreateTableStatement: createTbl,
		}, newCursor, true
	}

	return nil, initialCursor, false

}

func parseSelectStatement(tokens []*lexer.Token, initialCursor uint, delimiter lexer.Token) (*SelectStatement, uint, bool) {
	cursor := initialCursor

	if !expectToken(tokens, cursor, tokenFromKeyword(lexer.SelectKeyword)) {
		return nil, initialCursor, false
	}

	cursor++

	slct := SelectStatement{}

	exps, newCursor, ok := parseExpressions(tokens, cursor, []lexer.Token{delimiter, tokenFromKeyword(lexer.FromKeyword)})

	if !ok {
		return nil, initialCursor, false
	}

	slct.Item = *exps
	cursor = newCursor

	if expectToken(tokens, cursor, tokenFromKeyword(lexer.FromKeyword)) {
		cursor++

		from, newCursor, ok := parseToken(tokens, cursor, lexer.IdentifierKind)

		if !ok {
			helpMessage(tokens, cursor, "Expected FROM token")

			return nil, initialCursor, false
		}

		slct.From = *from

		cursor = newCursor
	}
	return &slct, cursor, true
}

func parseInsertStatement(tokens []*lexer.Token, initialCursor uint, delimiter lexer.Token) (*InsertStatement, uint, bool) {
	cursor := initialCursor

	// Look for INSERT

	if !expectToken(tokens, cursor, tokenFromKeyword(lexer.InsertKeyword)) {
		return nil, initialCursor, false
	}

	return nil, initialCursor, false
}
