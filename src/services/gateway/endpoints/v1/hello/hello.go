package hello

import (
	"github.com/kodflow/kitsune/src/internal/core/server/router"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
	"github.com/kodflow/kitsune/src/services/gateway/endpoints/v1/hello/world"
)

var (
	EndPoint = router.NewEndPoint("hello")
)

func init() {
	EndPoint.Get(func(req *generated.Request, res *generated.Response) error {
		res.Body = []byte("Hello World")
		return nil
	})

	EndPoint.Sub(world.EndPoint)
}
