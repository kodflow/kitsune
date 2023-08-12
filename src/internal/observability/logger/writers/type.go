package writers

import (
	"io"
	"os"
	"path"

	"github.com/kodmain/kitsune/src/internal/env"
	"github.com/kodmain/kitsune/src/internal/storages/fs"
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

func Make(t TYPE, sof SOF) io.Writer {
	ws := []io.Writer{}

	if t&CONSOLE != 0 {
		if sof {
			ws = append(ws, os.Stdout)
		} else {
			ws = append(ws, os.Stderr)
		}
	}

	if t&FILE != 0 {
		if sof {
			if f, err := fs.OpenFile(path.Join(env.PATH_LOGS, env.BUILD_APP_NAME+".log")); err == nil {
				ws = append(ws, f)
			}
		} else {
			if f, err := fs.OpenFile(path.Join(env.PATH_LOGS, env.BUILD_APP_NAME+".err")); err == nil {
				ws = append(ws, f)
			}
		}
	}

	return io.MultiWriter(ws...)
}
