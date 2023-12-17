package writers

import (
	"io"
	"os"
	"path"

	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/kernel/storages/fs"
)

// TYPE is a custom type representing different writer types.
type TYPE uint8

// SOF (Success Or Failure) is a boolean type used to indicate the nature of the log (success or failure).
type SOF bool

// Constants representing different logging destinations and log types.
const (
	CONSOLE TYPE = 1 << iota // Log to console.
	FILE                     // Log to file.
	AWS                      // Log to AWS (placeholder, not implemented).

	SUCCESS = true  // Indicates success log.
	FAILURE = false // Indicates failure or error log.

	DEFAULT = CONSOLE | FILE // Default log destinations.
)

// File paths for standard error and standard output logs.
var (
	FILE_STDERR = path.Join(config.PATH_LOGS, config.BUILD_APP_NAME+".err") // File path for standard error logs.
	FILE_STDOUT = path.Join(config.PATH_LOGS, config.BUILD_APP_NAME+".out") // File path for standard output logs.

	CONSOLE_STDERR = os.Stderr // Standard error console output.
	CONSOLE_STDOUT = os.Stdout // Standard output console output.
)

// Make creates an io.Writer based on the specified TYPE and SOF.
// It can create a multi-writer that writes to console and/or file based on the provided flags.
//
// Parameters:
// - t: TYPE The type of writer(s) to create.
// - sof: SOF Indicates whether the writer is for success (true) or failure (false) logs.
//
// Returns:
// - io.Writer: A writer that logs to the specified destinations.
func Make(t TYPE, sof SOF) io.Writer {
	ws := []io.Writer{}

	// Add the appropriate console writer based on the TYPE and SOF.
	if t&CONSOLE != 0 {
		if sof {
			ws = append(ws, CONSOLE_STDOUT)
		} else {
			ws = append(ws, CONSOLE_STDERR)
		}
	}

	// Open the appropriate file based on the TYPE and SOF, and add it to the writer list.
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

	// Return a multi-writer that writes to all the specified writers.
	return io.MultiWriter(ws...)
}
