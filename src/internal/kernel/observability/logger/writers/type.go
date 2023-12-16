package writers

import (
	"io"
	"os"
	"path"

	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/kernel/storages/fs"
)

type TYPE uint8
type SOF bool

const (
	CONSOLE TYPE = 1 << iota
	FILE
	AWS

	SUCCESS = true
	FAILURE = false

	DEFAULT = CONSOLE | FILE
)

var (
	FILE_STDERR = path.Join(config.PATH_LOGS, config.BUILD_APP_NAME+".err")
	FILE_STDOUT = path.Join(config.PATH_LOGS, config.BUILD_APP_NAME+".out")

	CONSOLE_STDERR = os.Stderr
	CONSOLE_STDOUT = os.Stdout
)

func Make(t TYPE, sof SOF) io.Writer {
	ws := []io.Writer{}

	// If the CONSOLE flag is set, add the appropriate console writer to the list
	if t&CONSOLE != 0 {
		if sof {
			ws = append(ws, CONSOLE_STDOUT)
		} else {
			ws = append(ws, CONSOLE_STDERR)
		}
	}

	// If the FILE flag is set, open the appropriate file and add it to the list
	if t&FILE != 0 {
		var target string
		if sof {
			target = FILE_STDOUT
		} else {
			target = FILE_STDERR
		}

		if f, err := fs.OpenFile(target); err == nil {
			ws = append(ws, f)
		}
	}

	// Return a multi-writer that writes to all the writers in the list
	return io.MultiWriter(ws...)
}
