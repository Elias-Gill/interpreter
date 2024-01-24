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
