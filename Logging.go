// Logging.go
package autositespeed

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogOptions struct {
	Enabled    bool
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

func initLogging(options *LogOptions) {
	log.SetOutput(&lumberjack.Logger{
		Filename:   options.FilePath,
		MaxSize:    options.MaxSize,
		MaxAge:     options.MaxAge,
		MaxBackups: options.MaxBackups,
	})
}

type logWriter struct {
}

func (lw *logWriter) Write(p []byte) (n int, err error) {
	log.Println(string(p))

	return len(p), nil
}
