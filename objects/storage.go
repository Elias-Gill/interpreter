package objects

type Storage struct {
	identifiers map[string]Object
	outer       *Storage // outer environment
}

func NewStorage() *Storage {
	return &Storage{
		identifiers: make(map[string]Object),
		outer:       nil,
	}
}

func NewEnclosedStorage(outer *Storage) *Storage {
	return &Storage{
		identifiers: make(map[string]Object),
		outer:       outer,
	}
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
