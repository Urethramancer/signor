// Package log contains the Logger, a complex (and possibly complicated) structure for
// logging of data from a longer-running process to different destinations.
// TODO: Define exactly what to do with JSON and RPC output.
package log

import (
	"fmt"
	"os"
	"strings"
)

// Default Logger object.
var Default *Logger

func init() {
	Default = NewLogger()
}

// Logger structure for configurable output.
type Logger struct {
	// msgF is the format of ordinary message. These sometimes don't need all the details.
	msgF string
	// errF is the format of errors. Users generally want every detail you can provide.
	errF     string
	servers  []string
	outFiles []*os.File
	validOut []bool
	logDst   byte
}

const (
	// O_FILE for stdout+stderr or a filename.
	O_FILE = 1
	// O_JSON for a JSON log server.
	O_JSON = 2
	// O_RPC for a gRPC server.
	O_RPC = 4
)

const (
	// DetailedFormat is the default log format which uses only the most essential fields.
	DetailedFormat = "%host %time: Level %level crisis from %name on %src: %msg"
)

// NewLogger creates a logger with some reasonable defaults for printing to stdout/stderr.
func NewLogger() *Logger {
	l := Logger{
		msgF:     DetailedFormat,
		errF:     DetailedFormat,
		servers:  make([]string, 2),
		outFiles: []*os.File{os.Stdout, os.Stderr},
		validOut: make([]bool, 2),
		logDst:   O_FILE,
	}
	return &l
}

// CloseFiles closes any open non-stdout/stderr files and replaces them with stdout.
func (l *Logger) CloseFiles() {
	for i := 0; i < 2; i++ {
		if l.validOut[i] {
			l.outFiles[i].Close()
			l.outFiles[i] = os.Stdout
		}
	}
}

// Msg prints arbitrary formatted messages to the configured message output(s).
func (l *Logger) Msg(f string, v ...interface{}) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(f, v...))
	b.WriteString("\n")
	if l.logDst&O_FILE == O_FILE {
		fmt.Fprint(l.outFiles[0], b.String())
	}
}

// TMsg prints arbitrary formatted messages to the configured message output(s),
// starting with a timestamp.
func (l *Logger) TMsg(f string, v ...interface{}) {
	var b strings.Builder
	b.WriteString(NowString())
	b.WriteRune(':')
	b.WriteString(fmt.Sprintf(f, v...))
	b.WriteString("\n")
	if l.logDst&O_FILE == O_FILE {
		fmt.Fprint(l.outFiles[0], b.String())
	}
}

// Err prints arbitrary formatted errors to the configured error output(s).
func (l *Logger) Err(f string, v ...interface{}) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(f, v...))
	b.WriteString("\n")
	if l.logDst&O_FILE == O_FILE {
		fmt.Fprint(l.outFiles[1], b.String())
	}
}

// TErr prints arbitrary formatted errors to the configured error output(s),
// starting with a timestamp.
func (l *Logger) TErr(f string, v ...interface{}) {
	var b strings.Builder
	b.WriteString(NowString())
	b.WriteRune(':')
	b.WriteString(fmt.Sprintf(f, v...))
	b.WriteString("\n")
	if l.logDst&O_FILE == O_FILE {
		fmt.Fprint(l.outFiles[1], b.String())
	}
}

// Log an event to an appropriate output in a configured format for that log level.
// Level 0 defaults to stdout, anything else to stderr.
func (l *Logger) Log(e *Event) {
	if l.logDst&O_FILE == O_FILE {
		if e.Level == 0 {
			fmt.Fprint(l.outFiles[0], e.Fmt(l.msgF))
		} else {
			fmt.Fprint(l.outFiles[1], e.Fmt(l.errF))
		}
	}
}

// SetFmt for messages and errors to the same format.
func (l *Logger) SetFmt(s string) {
	l.SetLogFmt(s)
	l.SetELogFmt(s)
}

// SetLogFmt sets the output format for informational event logs.
func (l *Logger) SetLogFmt(s string) {
	l.msgF = s
}

// SetELogFmt sets the output format for error event logs.
func (l *Logger) SetELogFmt(s string) {
	l.errF = s
}

// SetLogOut sets the output methods for messages and errors.
// files - filenames, or blank for stdout and stderr
// servers - host:port strings for remote logging destinations
// Specify O_FILE and blank files to use stdout and stderr.
// This can be combined with either O_JSON or O_RPC.
func (l *Logger) SetLogOut(log byte, files, servers []string) {
	l.logDst = log
	l.outFiles[0] = os.Stdout
	l.outFiles[1] = os.Stderr
	if files == nil || len(files) < 2 {
		return
	}

	var err error
	var f *os.File
	for i := 0; i < 2; i++ {
		if files[i] != "" {
			f, err = os.OpenFile(files[i], os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
			if err == nil {
				l.outFiles[i] = f
				l.validOut[i] = true
			}
		}
	}
}

// Warn is meant to be deferred with closing operations which might return an error.
// If t is true, the output will be timestamped with the default format of the logger.
// Any error returns 1 to the operating system, which is considered a warning/minor error.
func (l *Logger) Warn(err error, t bool) {
	if err == nil {
		return
	}

	if t {
		l.TErr("Error: %s", err.Error())
	} else {
		l.Err("Error: %s", err.Error())
	}
	os.Exit(1)
}

// Fail is meant to be deferred with closing operations which might return an error.
// If t is true, the output will be timestamped with the default format of the logger.
// Any error returns 2 to the operating system, which is considered a major error.
func (l *Logger) Fail(err error, t bool) {
	if err == nil {
		return
	}

	if t {
		l.TErr("Error: %s", err.Error())
	} else {
		l.Err("Error: %s", err.Error())
	}
	os.Exit(2)
}
