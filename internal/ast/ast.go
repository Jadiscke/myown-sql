package ast

import (
	"github.com/Jadiscke/myown-sql/internal/lexer"
)

type token lexer.Token

type expressionKind uint

const (
	LiteralKind expressionKind = iota
)

type expression struct {
	Literal *token
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
	Table  token
	Values *[]*expression
}

type columnDefinition struct {
	Name     token
	Datatype token
}

type CreateTableStatement struct {
	Name token
	Cols *[]*columnDefinition
}

type SelectStatement struct {
	Item []*expression
	From token
}
