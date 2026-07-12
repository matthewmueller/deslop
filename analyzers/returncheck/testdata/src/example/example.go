package example

type Store interface {
	Get(id string) string
}

type store struct{}

func (s *store) Get(id string) string { return "" }

// Bad: returns an interface
func NewStore() Store { // want "constructor NewStore should return a concrete type"
	return &store{}
}

// Good: returns concrete type
func NewConcreteStore() *store {
	return &store{}
}

// Good: returning error is fine
func LoadConfig() (*store, error) {
	return &store{}, nil
}
