package endpoints

import (
	"github.com/kodflow/kitsune/src/internal/core/server/router"
	v1 "github.com/kodflow/kitsune/src/services/gateway/endpoints/v1"
)

var (
	ROOT *router.EndPoint = router.NewRootPoint()
)

func init() {
	ROOT.Sub(v1.EndPoint)
}
