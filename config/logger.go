package config

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger sets up a rotating file logger and returns an io.Writer + *log.Logger
// Use the returned writer for other libraries (e.g. Fiber logger middleware, zap/stdlog wrapper)
func InitLogger(path string) (io.Writer, *log.Logger) {
	// ensure log directory exists
	// if path includes directory, create it (best-effort)
	// note: keep simple — if permission error occurs, fallback to stdout
	dir := ""
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			dir = path[:i]
			break
		}
	}
	if dir != "" {
		_ = os.MkdirAll(dir, 0o755)
	}

	lj := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    10, // megabytes
		MaxBackups: 7,
		MaxAge:     28,   // days
		Compress:   true, // compressed rotated files
	}

	// write to both lumberjack and stdout
	mw := io.MultiWriter(lj, os.Stdout)
	logger := log.New(mw, "", log.LstdFlags|log.Lshortfile)
	return mw, logger
}
