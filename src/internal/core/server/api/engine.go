package api

import (
	"github.com/kodmain/kitsune/src/internal/core/server/handler"
)

type Engine struct {
	Handler  func(b []byte) []byte
	versions map[string]*Router
}

func MakeEngine() *Engine {
	return &Engine{
		Handler:  handler.TCPHandler,
		versions: map[string]*Router{},
	}
}
