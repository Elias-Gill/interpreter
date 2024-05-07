package objects

import (
	"fmt"

	"github.com/sl2.0/ast"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ = "INTEGER"
	STRING_OBJ  = "STRING"
	BOOL_OBJ    = "BOOL"
	NULL_OBJ    = "NULL"
	ERROR_OBJ   = "ERROR"
	RETURN_OBJ  = "RETURN"
	FUNC_OBJ    = "FUNCTION"
)

// --- Primitive data types ---

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOL_OBJ
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}

type String struct {
	Value string
}

func (i *String) Type() ObjectType {
	return INTEGER_OBJ
}
func (i *String) Inspect() string {
	return i.Value
}

// --- Complex data types ---

type ErrorObject struct {
	error string
}

func (b *ErrorObject) Type() ObjectType {
	return ERROR_OBJ
}

func NewError(format string, message ...interface{}) Object {
	return &ErrorObject{error: fmt.Sprintf(format, message...)}
}

func (b *ErrorObject) Inspect() string {
	return b.error
}

type ReturnObject struct {
	Value Object
}

func (r *ReturnObject) Type() ObjectType {
	return RETURN_OBJ
}
func (r *ReturnObject) Inspect() string {
	return r.Value.Inspect()
}

type FunctionObject struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        Storage
}

func (f *FunctionObject) Type() ObjectType {
	return RETURN_OBJ
}
func (f *FunctionObject) Inspect() string {
	s := "("
	for _, param := range f.Parameters {
		s += param.ToString() + " "
	}
	s += ")"

	return s + "\n" + f.Body.ToString()
}
