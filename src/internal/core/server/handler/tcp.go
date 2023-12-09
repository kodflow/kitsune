package handler

import (
	"github.com/kodmain/kitsune/src/internal/core/server/router"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"google.golang.org/protobuf/proto"
)

func TCPHandler(b []byte) []byte {
	req, res := transport.New()

	err := proto.Unmarshal(b, req)
	if err != nil {
		return transport.Empty
	}

	err = router.Resolve(req, res)
	if err != nil {
		return transport.Empty
	}

	b, err = proto.Marshal(res)
	if err != nil {
		return transport.Empty
	}

	return b
}
