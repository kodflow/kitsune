package multithread

import (
	"flag"
)

var child *bool = flag.Bool("child", false, "determine if it's a child")

func IsWorker() bool {
	return *child
}
