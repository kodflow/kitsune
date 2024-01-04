package router

import "github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"

// Handler définit le type de fonction qui implémente HandlerInterface.
type Handler func(req *generated.Request, res *generated.Response) error
