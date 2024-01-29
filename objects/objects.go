package objects

import "fmt"

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ = "INTEGER"
	BOOL_OBJ    = "BOOL"
	NULL_OBJ    = "NULL"
	RETURN_OBJ  = "NULL"
)

// --- Object types ---

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

type NULL struct{}

func (b *NULL) Type() ObjectType {
	return NULL_OBJ
}
func (b *NULL) Inspect() string {
	return NULL_OBJ
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
