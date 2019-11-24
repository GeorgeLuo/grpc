package utils

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// TODO: log level ignore setting

// this class manages the write output of commands to a file.

// StaticLogger is a map of taskID to previous processes. The abstraction is defined
// to enforce thread-safety of read and write to CommandWrapper objects.
type StaticLogger struct {
	Filename *string
	mutex    *sync.Mutex
	file     *os.File
	prepend  string
	Level    Level
}

// NewStaticLogger is initialize a static logger.
func NewStaticLogger(logDir string, fn string) StaticLogger {
	filename := logDir + fn + "." + time.Now().Format("20060102150405")
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error initializing logger:\n %s\n", err.Error())
	}
	return StaticLogger{
		Filename: &filename,
		mutex:    &sync.Mutex{},
		file:     f,
		prepend:  "",
		Level:    D,
	}
}

// SetGlobalPrepend sets the global string prepended to loglines.
// timed will enable the logging of datetime before every line.
func (logger *StaticLogger) SetGlobalPrepend(prepend ...string) {
	logger.prepend = strings.Join(prepend, " | ")
}

// Level defines the log levels
type Level int

const (
	// D = Debug
	D Level = iota
	// I = Info
	I
	// W = Warning
	W
	// E = Error
	E
	// F = Fatal
	F
)

// SetLevel sets the level level at which a log will be printed.
// if level is set to W, on W, E, F level logs will ever be printed.
func (logger *StaticLogger) SetLevel(level Level) {
	logger.Level = level
}

// TODO: input checking

// WriteDtLine is an operation write a new line into the logfile
func (logger *StaticLogger) WriteDtLine(line string) {
	var logLine string
	logLine += time.Now().Format("2006-01-02T15:04:05Z07:00") + " | "
	if logger.prepend != "" {
		logLine += logger.prepend + " | "
	}
	logger.WriteLine(logLine + line + "\n")
}

// WriteLine is the base write instruction.
func (logger *StaticLogger) WriteLine(line string) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	_, err := logger.file.WriteString(line)
	if err != nil {
		fmt.Printf("error writing line '%s':\n %s\n", line, err.Error())
	}
	logger.file.Sync()
}

// WriteDtDebug is the write instruction that prepends debug notation.
func (logger *StaticLogger) WriteDtDebug(line string) {
	if logger.Level <= D {
		logger.WriteDtLine("DEBUG | " + line)
	}
}

// WriteDtInfo is the write instruction that prepends info notation.
func (logger *StaticLogger) WriteDtInfo(line string) {
	if logger.Level <= I {
		logger.WriteDtLine("INFO | " + line)
	}
}

// WriteDtWarn is the write instruction that prepends warning notation.
func (logger *StaticLogger) WriteDtWarn(line string) {
	if logger.Level <= W {
		logger.WriteDtLine("WARN | " + line)
	}
}

// WriteDtError is the write instruction that prepends error notation.
func (logger *StaticLogger) WriteDtError(line string) {
	if logger.Level <= E {
		logger.WriteDtLine("ERROR | " + line)
	}
}

// WriteDtFatal is the write instruction that prepends fatal notation.
func (logger *StaticLogger) WriteDtFatal(line string) {
	if logger.Level <= F {
		logger.WriteDtLine("FATAL | " + line)
	}
}

// Close is an operation close the logger
func (logger *StaticLogger) Close() {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	err := logger.file.Close()
	if err != nil {
		fmt.Printf("error closing logger:\n %s\n", err.Error())
	}
}
