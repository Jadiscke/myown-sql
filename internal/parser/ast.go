package parser

import (
	"github.com/Jadiscke/myown-sql/internal/lexer"
)

type expressionKind uint

const (
	LiteralKind expressionKind = iota
)

type expression struct {
	Literal *lexer.Token
	Kind    expressionKind
}

type Ast struct {
	Statements []*Statement
}

type AstKind uint

const (
	SelectKind AstKind = iota
	CreateTableKind
	InsertKind
)

type Statement struct {
	SelectStatement      *SelectStatement
	CreateTableStatement *CreateTableStatement
	InsertStatement      *InsertStatement
	Kind                 AstKind
}

type InsertStatement struct {
	Table  lexer.Token
	Values *[]*expression
}

type columnDefinition struct {
	Name     lexer.Token
	Datatype lexer.Token
}

type CreateTableStatement struct {
	Name lexer.Token
	Cols *[]*columnDefinition
}

type SelectStatement struct {
	Item []*expression
	From lexer.Token
}
