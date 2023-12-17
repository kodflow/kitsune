package api

// Handler is a function type that represents an API handler.
// Each handler is a function that returns an error, allowing for error handling in API operations.
type Handler = func() error

// APInterface defines the interface for an API router.
// It specifies methods for associating handlers with HTTP methods and URL paths.
type APInterface interface {
	Head(string, ...*Handler)   // Associate handlers with the HEAD HTTP method.
	Get(string, ...*Handler)    // Associate handlers with the GET HTTP method.
	Post(string, ...*Handler)   // Associate handlers with the POST HTTP method.
	Put(string, ...*Handler)    // Associate handlers with the PUT HTTP method.
	Patch(string, ...*Handler)  // Associate handlers with the PATCH HTTP method.
	Delete(string, ...*Handler) // Associate handlers with the DELETE HTTP method.
}

// Config represents the configuration for an API.
// It holds settings that affect the behavior of the API, such as deprecation status and version.
type Config struct {
	Depreciated bool   // Indicates if the API is deprecated.
	Version     string // The version of the API.
}
