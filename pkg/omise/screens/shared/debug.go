package shared

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

var (
	debugWriter io.Writer
	debugMutex  sync.Mutex
	debugMode   bool
)

// InitDebug initializes debug mode if DEBUG environment variable is set
// Messages will be dumped to /tmp/bento-debug.log
func InitDebug() {
	if os.Getenv("DEBUG") == "1" || os.Getenv("DEBUG") == "true" {
		file, err := os.OpenFile("/tmp/bento-debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open debug log: %v\n", err)
			return
		}

		debugWriter = file
		debugMode = true

		fmt.Fprintln(debugWriter, "=== Debug Mode Enabled ===")
		fmt.Fprintf(debugWriter, "To tail this file: tail -f /tmp/bento-debug.log\n\n")
	}
}

// CloseDebug closes the debug log file
func CloseDebug() {
	debugMutex.Lock()
	defer debugMutex.Unlock()

	if debugWriter != nil {
		if closer, ok := debugWriter.(io.Closer); ok {
			closer.Close()
		}
		debugWriter = nil
		debugMode = false
	}
}

// DebugMsg logs a message to the debug file if debug mode is enabled
func DebugMsg(msg interface{}, context string) {
	if !debugMode {
		return
	}

	debugMutex.Lock()
	defer debugMutex.Unlock()

	if debugWriter != nil {
		if context != "" {
			fmt.Fprintf(debugWriter, "\n=== %s ===\n", context)
		}
		spew.Fdump(debugWriter, msg)
		fmt.Fprintln(debugWriter, "")
	}
}

// DebugPrintf logs a formatted message to the debug file
func DebugPrintf(format string, args ...interface{}) {
	if !debugMode {
		return
	}

	debugMutex.Lock()
	defer debugMutex.Unlock()

	if debugWriter != nil {
		fmt.Fprintf(debugWriter, format, args...)
		fmt.Fprintln(debugWriter, "")
	}
}

// IsDebugMode returns whether debug mode is enabled
func IsDebugMode() bool {
	return debugMode
}
