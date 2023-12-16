package multithread

import (
	"flag"
)

// child represents a flag to determine if it's a child.
var child *bool = flag.Bool("child", false, "determine if it's a child")

// IsWorker checks if the current process is a worker.
func IsWorker() bool {
	return *child
}
