package api

type Handler = func() error
type APInterface interface {
	//Use(string, ...*Handler)
	Head(string, ...*Handler)
	Get(string, ...*Handler)
	Post(string, ...*Handler)
	Put(string, ...*Handler)
	Patch(string, ...*Handler)
	Delete(string, ...*Handler)
	//All(string, ...*Handler)
}

type Config struct {
	Depreciated bool
	Version     string
}
