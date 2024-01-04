package status

import (
	"github.com/kodflow/kitsune/src/internal/core/server/router"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
)

var (
	EndPoint = router.NewEndPoint("status")
)

func init() {
	EndPoint.Get(func(req *generated.Request, res *generated.Response, next router.HandlerInterface) (*generated.Request, *generated.Response, error) {
		res.Body = []byte("STATUS")
		return req, res, nil
	})
}
