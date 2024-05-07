package storage

import "github.com/sl2.0/objects"

type Env struct {
	identifiers map[string]objects.Object
}

func NewEnv() *Env {
	return &Env{
		identifiers: make(map[string]objects.Object),
	}
}

func (e *Env) Get(ident string) (objects.Object, bool) {
	value, ok := e.identifiers[ident]
	return value, ok
}

func (e *Env) Set(ident string, obj objects.Object) objects.Object {
	e.identifiers[ident] = obj
	return obj
}
