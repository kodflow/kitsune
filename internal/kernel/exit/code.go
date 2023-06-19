package exit

import (
	"fmt"
	"os"
)

const (
	Success int = 0
	Error   int = iota + 1
)

func App(code int, message ...string) {
	fmt.Println(message)
	os.Exit(code)
}
