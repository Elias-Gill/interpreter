package objects

import "fmt"

type Storage struct {
	identifiers map[string]Object
	outer       *Storage // outer environment
	lvl         int      // to meassure recursion lvl
}

func NewStorage() *Storage {
	return &Storage{
		identifiers: make(map[string]Object),
		outer:       nil,
		lvl:         0,
	}
}

func NewEnclosedStorage(outer *Storage) (*Storage, error) {
	lvl := outer.lvl + 1
	if lvl > 200 {
		return nil, fmt.Errorf("Max level of recursion reached")
	}

	return &Storage{
		identifiers: make(map[string]Object),
		outer:       outer,
		lvl:         lvl,
	}, nil
}

func (e *Storage) Get(ident string) (Object, bool) {
	value, ok := e.identifiers[ident]

	if !ok && e.outer != nil {
		return e.outer.Get(ident)
	}

	return value, ok
}

func (e *Storage) Set(ident string, obj Object) Object {
	e.identifiers[ident] = obj
	return obj
}
