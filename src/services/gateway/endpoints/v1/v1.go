package v1

import (
	"github.com/kodflow/kitsune/src/internal/core/server/router"
	"github.com/kodflow/kitsune/src/services/gateway/endpoints/v1/hello"
	"github.com/kodflow/kitsune/src/services/gateway/endpoints/v1/status"
)

var (
	EndPoint *router.EndPoint = router.NewEndPoint("v1")
)

func init() {
	EndPoint.Sub(status.EndPoint)
	EndPoint.Sub(hello.EndPoint)
}
