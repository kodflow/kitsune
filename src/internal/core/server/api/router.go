package api

// Router represents your API
type Router struct {
	Depreciated bool
	head        map[string][]*Handler
	get         map[string][]*Handler
	post        map[string][]*Handler
	put         map[string][]*Handler
	patch       map[string][]*Handler
	delete      map[string][]*Handler
}

// Make creates an instance of the API based on the provided configuration.
func Make(cfg *Config) APInterface {
	return &Router{
		Depreciated: cfg.Depreciated,
		head:        map[string][]*Handler{},
		get:         map[string][]*Handler{},
		post:        map[string][]*Handler{},
		put:         map[string][]*Handler{},
		patch:       map[string][]*Handler{},
		delete:      map[string][]*Handler{},
	}
}

// Head associates one or more handlers with the HEAD method.
func (r *Router) Head(url string, handlers ...*Handler) {
	r.head[url] = append(r.head[url], handlers...)
}

// Get associates one or more handlers with the GET method.
func (r *Router) Get(url string, handlers ...*Handler) {
	r.get[url] = append(r.get[url], handlers...)
}

// Post associates one or more handlers with the POST method.
func (r *Router) Post(url string, handlers ...*Handler) {
	r.post[url] = append(r.post[url], handlers...)
}

// Put associates one or more handlers with the PUT method.
func (r *Router) Put(url string, handlers ...*Handler) {
	r.put[url] = append(r.put[url], handlers...)
}

// Patch associates one or more handlers with the PATCH method.
func (r *Router) Patch(url string, handlers ...*Handler) {
	r.patch[url] = append(r.patch[url], handlers...)
}

// Delete associates one or more handlers with the DELETE method.
func (r *Router) Delete(url string, handlers ...*Handler) {
	r.delete[url] = append(r.delete[url], handlers...)
}
