package user

import (
	"github.com/kodflow/kitsune/src/internal/core/server/api"
)

var V1 api.APInterface = api.Make(&api.Config{
	Depreciated: true,
	Version:     "v1",
})
