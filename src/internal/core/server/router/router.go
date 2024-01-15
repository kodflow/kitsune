package router

import (
	"fmt"
	"strings"

	"github.com/kodflow/kitsune/src/internal/core/server/transport"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
)

func (r *Router) Register(epi *EndPoint) error {
	if epi == nil {
		return (fmt.Errorf("endpoint is not defined"))
	} else if !epi.isRootEndpoint() {
		return (fmt.Errorf("%v is not a root endpoint", epi.Endpoint))
	}

	r.endpoint = epi
	logger.Info("Register endpoint: ")
	for _, e := range epi.traverse() {
		url := e.URL()
		for _, method := range e.options {
			logger.Infof("%v %v", method, url)
		}
	}

	return nil
}

// Resolve process the request and applies appropriate handlers
//
// This function takes a router and an exchange object. It resolves the endpoint from the request,
// then applies the corresponding handlers based on the request method.
//
// Parameters:
// - exchange: *transport.Exchange The exchange object containing request and response.
//
// Returns:
// - err: error The error encountered during processing, if any.
func (r *Router) Resolve(exchange *transport.Exchange) {
	//time.Sleep(100 * time.Millisecond)
	req := exchange.Request()
	res := exchange.Response()

	if r.endpoint == nil {
		res.Status = 404
		return
	}

	// Simplify the extraction of endpoint names
	endpointNames := simplifyEndpointNames(req.Endpoint)

	var endpoint *EndPoint
	var ok bool

	// Find the appropriate endpoint
	for _, name := range endpointNames {
		endpoint, ok = r.getEndpoint(endpoint, name)
		if !ok {
			return
		}
	}

	// Process with the found endpoint
	if endpoint != nil {
		r.processEndpoint(endpoint, req, res)
	}
}

// simplifyEndpointNames trims and splits the endpoint string
//
// Parameters:
// - endpointStr: string The endpoint URL string.
//
// Returns:
// - []string The sliced parts of the endpoint.
func simplifyEndpointNames(endpointStr string) []string {
	return strings.Split(strings.Trim(endpointStr, "/"), "/")
}

// getEndpoint finds the nested endpoint based on the name
//
// Parameters:
// - currentEndpoint: *EndPoint The current endpoint to start the search from.
// - name: string The name of the next endpoint to find.
//
// Returns:
// - *EndPoint The found endpoint or nil.
// - bool Indicates if the endpoint was found.
func (r *Router) getEndpoint(currentEndpoint *EndPoint, name string) (*EndPoint, bool) {
	if currentEndpoint == nil {
		ep, ok := r.endpoint.subs[name]
		return ep, ok
	}
	ep, ok := currentEndpoint.subs[name]
	return ep, ok
}

// processEndpoint applies handlers for the endpoint based on the request method
//
// Parameters:
// - endpoint: *EndPoint The endpoint to process.
// - req: *Request The request object.
// - res: *Response The response object.
//
// Returns:
// - error The error encountered during processing, if any.
func (r *Router) processEndpoint(endpoint *EndPoint, req *generated.Request, res *generated.Response) error {
	if handlers, ok := endpoint.handlers[req.Method]; ok {
		for _, handler := range handlers {
			if err := handler(req, res); err != nil {
				return err
			}
		}
	}
	return nil
}

// Router represents your API.
// It manages the association of URL paths with their corresponding handlers based on
// HTTP methods like GET, POST, PUT, PATCH, and DELETE. The Router also keeps track of
// whether it is deprecated.
type Router struct {
	endpoint *EndPoint
}

// MakeRouter creates and returns a new instance of Router.
// This function initializes a Router with its default values.
//
// Returns:
// - *Router: A new instance of Router.
func MakeRouter() *Router {
	return &Router{}
}
