package router

import "github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"

// HandlerInterface définit une interface pour le handler.
type HandlerInterface interface {
	Handle(*generated.Request, *generated.Response) (*generated.Request, *generated.Response, error)
}

// Handler définit le type de fonction qui implémente HandlerInterface.
type Handler func(req *generated.Request, res *generated.Response, next HandlerInterface) (*generated.Request, *generated.Response, error)

// Call implémente HandlerInterface pour une Handler.
func (f Handler) Call(req *generated.Request, res *generated.Response, next HandlerInterface) (*generated.Request, *generated.Response, error) {
	return f(req, res, next)
}
