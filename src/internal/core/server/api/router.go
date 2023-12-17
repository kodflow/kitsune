package api

// Router represents your API.
// It manages the association of URL paths with their corresponding handlers based on
// HTTP methods like GET, POST, PUT, PATCH, and DELETE. The Router also keeps track of
// whether it is deprecated.
type Router struct {
	Depreciated bool                  // Indicates if the router is deprecated.
	head        map[string][]*Handler // Map of HEAD method handlers.
	get         map[string][]*Handler // Map of GET method handlers.
	post        map[string][]*Handler // Map of POST method handlers.
	put         map[string][]*Handler // Map of PUT method handlers.
	patch       map[string][]*Handler // Map of PATCH method handlers.
	delete      map[string][]*Handler // Map of DELETE method handlers.
}

// Make creates an instance of the API based on the provided configuration.
// It initializes the Router with the deprecation status and empty handler maps for
// each HTTP method.
//
// Parameters:
// - cfg: *Config Configuration data used to set up the router.
//
// Returns:
// - APInterface: An instance of Router fulfilling the APInterface.
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
// This method allows adding handlers to the router that respond to HEAD requests
// for a specific URL path.
//
// Parameters:
// - url: string The URL path to associate with the handlers.
// - handlers: ...*Handler The handlers to be associated with the URL path.
func (r *Router) Head(url string, handlers ...*Handler) {
	r.head[url] = append(r.head[url], handlers...)
}

// Get associates one or more handlers with the GET method.
// Similar to Head, this method allows adding handlers that respond to GET requests.
//
// Parameters:
// - url: string The URL path to associate with the handlers.
// - handlers: ...*Handler The handlers to be associated with the URL path.
func (r *Router) Get(url string, handlers ...*Handler) {
	r.get[url] = append(r.get[url], handlers...)
}

// Post associates one or more handlers with the POST method.
// This method is used for associating handlers that handle POST requests.
//
// Parameters:
// - url: string The URL path to associate with the handlers.
// - handlers: ...*Handler The handlers to be associated with the URL path.
func (r *Router) Post(url string, handlers ...*Handler) {
	r.post[url] = append(r.post[url], handlers...)
}

// Put associates one or more handlers with the PUT method.
// It allows for handlers that manage PUT requests to be linked to a specific URL path.
//
// Parameters:
// - url: string The URL path to associate with the handlers.
// - handlers: ...*Handler The handlers to be associated with the URL path.
func (r *Router) Put(url string, handlers ...*Handler) {
	r.put[url] = append(r.put[url], handlers...)
}

// Patch associates one or more handlers with the PATCH method.
// This method facilitates adding handlers for PATCH requests to the router.
//
// Parameters:
// - url: string The URL path to associate with the handlers.
// - handlers: ...*Handler The handlers to be associated with the URL path.
func (r *Router) Patch(url string, handlers ...*Handler) {
	r.patch[url] = append(r.patch[url], handlers...)
}

// Delete associates one or more handlers with the DELETE method.
// It allows for setting up handlers to respond to DELETE requests on a given URL path.
//
// Parameters:
// - url: string The URL path to associate with the handlers.
// - handlers: ...*Handler The handlers to be associated with the URL path.
func (r *Router) Delete(url string, handlers ...*Handler) {
	r.delete[url] = append(r.delete[url], handlers...)
}

// MakeRouter creates and returns a new instance of Router.
// This function initializes a Router with its default values.
//
// Returns:
// - *Router: A new instance of Router.
func MakeRouter() *Router {
	return &Router{
		head:   map[string][]*Handler{},
		get:    map[string][]*Handler{},
		post:   map[string][]*Handler{},
		put:    map[string][]*Handler{},
		patch:  map[string][]*Handler{},
		delete: map[string][]*Handler{},
	}
}
