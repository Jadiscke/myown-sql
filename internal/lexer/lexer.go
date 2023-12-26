package lexer

import (
	"fmt"
	"strings"
)

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
	whereKeyword  keyword = "where"
)

type symbol string

const (
	semicolonSymbol  symbol = ";"
	asteriskSymbol   symbol = "*"
	commaSymbol      symbol = ","
	leftParenSymbol  symbol = "("
	rightParenSymbol symbol = ")"
	concatSymbol     symbol = "||"
)

type tokenKind uint

const (
	keywordKind tokenKind = iota
	symbolKind
	istringKind
	stringKind
	numericKind
	boolKind
	identifierKind
)

type Token struct {
	value string
	kind  tokenKind
	loc   location
}

type cursor struct {
	pointer uint
	loc     location
}

func (t *Token) Equals(other *Token) bool {
	return t.value == other.value && t.kind == other.kind
}

type lexer func(string, cursor) (*Token, cursor, bool)

func lex(source string) ([]*Token, error) {
	tokens := []*Token{}
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

func lexNumeric(source string, ic cursor) (*Token, cursor, bool) {
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

	return &Token{
		value: source[ic.pointer:cur.pointer],
		kind:  numericKind,
		loc:   ic.loc,
	}, cur, true

}

func lexCharacterDelimited(source string, ic cursor, delimiter byte) (*Token, cursor, bool) {
	cur := ic

	if len(source[cur.pointer:]) == 0 {
		return nil, ic, false
	}

	if source[cur.pointer] != delimiter {
		return nil, ic, false
	}

	cur.loc.col++
	cur.pointer++

	var value []byte

	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		character := source[cur.pointer]

		if character == delimiter {
			// To escape ' in SQL you should use ''
			// Example 'It''s a good day to be alive'
			if cur.pointer+1 >= uint(len(source)) || source[cur.pointer+1] != delimiter {
				return &Token{
					value: string(value),
					loc:   ic.loc,
					kind:  stringKind,
				}, cur, true
			} else {
				value = append(value, character)
				cur.pointer++
				cur.loc.col++
			}
		}

		value = append(value, character)
		cur.loc.col++
	}

	return nil, ic, false
}

func lexString(source string, ic cursor) (*Token, cursor, bool) {
	return lexCharacterDelimited(source, ic, '\'')
}

func longestMatch(source string, ic cursor, options []string) string {
	var value []byte
	var skipList []int
	var match string

	cur := ic

	for cur.pointer < uint(len(source)) {

		value = append(value, strings.ToLower(string(source[cur.pointer]))...)

		cur.pointer++

	match:
		for i, option := range options {
			for _, skip := range skipList {
				if i == skip {
					continue match
				}
			}
			// Deal with cases like INT vs INTO

			if option == string(value) {
				skipList = append(skipList, i)
				if len(option) > len(match) {
					match = option
				}

				continue
			}

			sharesPrefix := string(value) == option[:cur.pointer-ic.pointer]

			tooLong := len(value) > len(option)

			if tooLong || !sharesPrefix {
				skipList = append(skipList, i)
			}
		}

		if len(skipList) == len(options) {
			break
		}
	}

	return match
}

func lexSymbol(source string, ic cursor) (*Token, cursor, bool) {
	character := source[ic.pointer]

	cur := ic

	cur.pointer++
	cur.loc.col++

	switch character {
	// Remove whitespaces
	case '\n':
		cur.loc.line++
		cur.loc.col = 0
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		return nil, cur, true
	}

	symbols := []symbol{
		commaSymbol,
		leftParenSymbol,
		rightParenSymbol,
		semicolonSymbol,
		asteriskSymbol,
	}

	var options []string

	for _, sym := range symbols {
		options = append(options, string(sym))
	}

	match := longestMatch(source, ic, options)

	if match == "" {
		return nil, ic, false
	}

	cur.pointer = ic.pointer + uint(len(match))
	cur.loc.col = ic.loc.col + uint(len(match))

	return &Token{
		value: match,
		kind:  symbolKind,
		loc:   ic.loc,
	}, cur, true
}

func lexKeyword(source string, ic cursor) (*Token, cursor, bool) {
	cur := ic

	keywords := []keyword{
		selectKeyword,
		insertKeyword,
		valuesKeyword,
		tableKeyword,
		createKeyword,
		whereKeyword,
		fromKeyword,
		intoKeyword,
		textKeyword,
		intKeyword,
	}

	var options []string

	for _, kw := range keywords {
		options = append(options, string(kw))
	}

	match := longestMatch(source, ic, options)

	if match == "" {
		return nil, ic, false
	}

	cur.pointer = ic.pointer + uint(len(match))
	cur.loc.col = ic.loc.col + uint(len(match))

	return &Token{
		value: match,
		kind:  keywordKind,
		loc:   ic.loc,
	}, cur, true
}

func lexIdentifier(source string, ic cursor) (*Token, cursor, bool) {
	//Handle separetely if is a double-quoted identifier
	if token, newCursor, ok := lexCharacterDelimited(source, ic, '"'); ok {
		return token, newCursor, true
	}

	cur := ic

	character := source[cur.pointer]

	isAlphabetical := (character >= 'A' && character <= 'Z') || (character >= 'a' && character <= 'z')

	if !isAlphabetical {
		return nil, ic, false
	}

	cur.pointer++
	cur.loc.col++

	value := []byte{character}

	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		character = source[cur.pointer]

		isAlphabetical = (character >= 'A' && character <= 'Z') || (character >= 'a' && character <= 'z')
		isNumeric := character >= '0' && character <= '9'

		if isAlphabetical || isNumeric || character == '_' || character == '$' {
			value = append(value, character)

			cur.loc.col++
			continue
		}

		break
	}

	if len(value) == 0 {
		return nil, ic, false
	}

	return &Token{
		value: strings.ToLower(string(value)),
		loc:   ic.loc,
		kind:  identifierKind,
	}, cur, true
}