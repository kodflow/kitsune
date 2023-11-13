package user

import "github.com/kodmain/kitsune/src/internal/core/server/api"

var V2 api.APInterface = api.Make(&api.Config{
	Version: "v2",
})
