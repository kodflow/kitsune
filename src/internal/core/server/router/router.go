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
	endpoints := strings.Split(req.Endpoint, "/")

	logger.Debug(endpoints)

	/*
		endpoint := r.endpoints[endpoints[0]]
		for _, endpointName := range endpoints[2:] {
			endpoint
		}
	*/

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
