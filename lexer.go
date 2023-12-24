package myownsql

import "fmt"

type location struct {
	line uint
	col  uint
}

type keyword string

const (
	selectKeyword keyword = "select"
	fromKeyword   keyword = "from"
	asKeyord      keyword = "as"
	tableKeyword  keyword = "table"
	createKeyword keyword = "create"
	insertKeyword keyword = "insert"
	intoKeyword   keyword = "into"
	valuesKeyword keyword = "values"
	intKeyword    keyword = "int"
	textKeyword   keyword = "text"
)

type symbol string

const (
	semicolonSymbol  symbol = ";"
	asteriskSymbol   symbol = "*"
	commaSymbol      symbol = ","
	leftParenSymbol  symbol = "("
	rightParenSymbol symbol = ")"
)

type tokenKind uint

const (
	keywordKind tokenKind = iota
	symbolKind
	identifierKind
	stringKind
	numericKind
)

type token struct {
	value string
	kind  tokenKind
	loc   location
}

type cursor struct {
	pointer uint
	loc     location
}

func (t *token) equals(other *token) bool {
	return t.value == other.value && t.kind == other.kind
}

type lexer func(string, cursor) (*token, cursor, bool)

func lex(source string) ([]*token, error) {
	tokens := []*token{}
	cur := cursor{}

lex:
	for cur.pointer < uint(len(source)) {
		lexers := []lexer{
			lexKeyword,
			lexSymbol,
			lexString,
			lexNumeric,
			lexIdentifier,
		}

		for _, lexer := range lexers {
			if token, newCursor, ok := lexer(source, cur); ok {
				cur = newCursor

				if token != nil {
					tokens = append(tokens, token)
				}

				continue lex
			}
		}
		hint := ""

		if len(tokens) > 0 {
			hint = " after " + tokens[len(tokens)-1].value
		}

		return nil, fmt.Errorf("Unable to lex token at %d:%d, token %s", cur.loc.line, cur.loc.col, hint)
	}

	return tokens, nil

}

func lexNumeric(source string, ic cursor) (*token, cursor, bool) {
	cur := ic

	periodFound := false
	expMarkerFound := false

	for ; cur.pointer < uint(len(source)); cur.pointer++ {

		character := source[cur.pointer]
		cur.loc.col++

		isDigit := character >= '0' && character <= '9'
		isPeriod := character == '.'
		isExpMarker := character == 'e'
		lastCharPosition := uint(len(source) - 1)

		// Must start with a digit or period
		if cur.pointer == ic.pointer {
			if !isDigit && !isPeriod {
				return nil, ic, false
			}

			periodFound = isPeriod
			continue
		}

		if isPeriod {
			if periodFound {
				return nil, ic, false
			}

			periodFound = true
			continue

		}

		if isExpMarker {
			if expMarkerFound {
				return nil, ic, false
			}

			//No periods allowed after expMarker

			periodFound = true
			expMarkerFound = true

			if cur.pointer == lastCharPosition {
				return nil, ic, false
			}

			nextCharacter := source[cur.pointer+1]

			if nextCharacter == '+' || nextCharacter == '-' {
				cur.pointer++
				cur.loc.col++
			}

			continue

		}

		if !isDigit {
			break
		}
	}

	if cur.pointer == ic.pointer {
		return nil, ic, false
	}

	return &token{
		value: source[ic.pointer:cur.pointer],
		kind:  numericKind,
		loc:   ic.loc,
	}, cur, true

}
