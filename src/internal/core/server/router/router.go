package router

import (
	"fmt"
	"strings"

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

// Resolve resolves a transport request and generates a transport response.
//
// Parameters:
// - req: *generated.Request - The transport request to be resolved.
// - res: *generated.Response - The transport response to be generated.
//
// Returns:
// - error: An error if there was an issue resolving the request, otherwise nil.
func (r *Router) Resolve(req *generated.Request, res *generated.Response) error {

	var endpoint *EndPoint
	var err error

	endpointNames := strings.Split(
		strings.TrimPrefix(
			strings.TrimSuffix(
				req.Endpoint, "/",
			), "/",
		), "/",
	)

	for _, endpointName := range endpointNames {
		if endpoint == nil && endpointName == "" {
			endpoint = r.endpoint
		} else if endpoint == nil {
			ep, ok := r.endpoint.subs[endpointName]
			if !ok {
				endpoint = nil
				break
			}
			endpoint = ep
		} else {
			ep, ok := endpoint.subs[endpointName]
			if !ok {
				endpoint = nil
				break
			}
			endpoint = ep
		}
	}

	logger.Debug(endpoint)

	if endpoint != nil {
		if condition, ok := endpoint.handlers[req.Method]; ok {
			for _, handler := range condition {
				err = handler(req, res)
			}
		}
	}

	return err
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
