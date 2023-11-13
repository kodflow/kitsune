package api

type APInterface interface {
	Head(string, ...*Handler)
	Get(string, ...*Handler)
	Post(string, ...*Handler)
	Put(string, ...*Handler)
	Patch(string, ...*Handler)
	Delete(string, ...*Handler)
}

type Handler = func() error

type Config struct {
	Depreciated bool
	Version     string
}

// API represents your API.
type API struct {
	Depreciated bool
	Version     string
	head        map[string][]*Handler
	get         map[string][]*Handler
	post        map[string][]*Handler
	put         map[string][]*Handler
	patch       map[string][]*Handler
	delete      map[string][]*Handler
}

// Head associates one or more handlers with the HEAD method.
func (api *API) Head(url string, handlers ...*Handler) {
	api.head[url] = append(api.head[url], handlers...)
}

// Get associates one or more handlers with the GET method.
func (api *API) Get(url string, handlers ...*Handler) {
	api.get[url] = append(api.get[url], handlers...)
}

// Post associates one or more handlers with the POST method.
func (api *API) Post(url string, handlers ...*Handler) {
	api.post[url] = append(api.post[url], handlers...)
}

// Put associates one or more handlers with the PUT method.
func (api *API) Put(url string, handlers ...*Handler) {
	api.put[url] = append(api.put[url], handlers...)
}

// Patch associates one or more handlers with the PATCH method.
func (api *API) Patch(url string, handlers ...*Handler) {
	api.patch[url] = append(api.patch[url], handlers...)
}

// Delete associates one or more handlers with the DELETE method.
func (api *API) Delete(url string, handlers ...*Handler) {
	api.delete[url] = append(api.delete[url], handlers...)
}

// Make creates an instance of the API based on the provided configuration.
func Make(cfg *Config) APInterface {
	return &API{
		Depreciated: cfg.Depreciated,
		Version:     cfg.Version,
		head:        map[string][]*Handler{},
		get:         map[string][]*Handler{},
		post:        map[string][]*Handler{},
		put:         map[string][]*Handler{},
		patch:       map[string][]*Handler{},
		delete:      map[string][]*Handler{},
	}
}
