package log

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

// Event line in a log file. Status, warnings, failures etc.
type Event struct {
	strings.Builder
	// Level of seriousness. App-specific, but generally 0 = informational,
	// and higher numbers increase criticality.
	Level uint `json:"level,omitempty"`
	// PID is a process identifier, if relevant.
	PID int `json:"pid,omitempty"`
	// Time is the timestamp of the event.
	Time time.Time `json:"timestamp,omitempty"`
	// Name of app or sub-system where event started.
	Name string `json:"name,omitempty"`
	// Hostname of the system the event originated from.
	Hostname string `json:"hostname,omitempty"`
	// Source is app-specific.
	Source string `json:"source,omitempty"`
	// Message is the human-readable form of the event message.
	Message string `json:"message,omitempty"`
	// Extra strings for whatever.
	Extra []string `json:"extra,omitempty"`
}

// Log format keywords
const (
	fmtLevel   = "%level"
	fmtPID     = "%pid"
	fmtTime    = "%time"
	fmtName    = "%name"
	fmtHost    = "%host"
	fmtSource  = "%src"
	fmtMessage = "%msg"
	fmtExtra   = "%extra"
)

// Fmt creates a log event string from the provided format.
// Use Event.String() to look it up again without reparsing.
func (e *Event) Fmt(f string) string {
	e.Reset()
	for len(f) > 0 {
		c := f[0]
		if c == '%' {
			var key string
			key, f = e.parseKeyword(f)
			switch key {
			case fmtLevel:
				e.WriteString(fmt.Sprintf("%d", e.Level))
			case fmtPID:
				e.WriteString(fmt.Sprintf("%d", e.PID))
			case fmtTime:
				e.WriteString(NowString())
			case fmtName:
				e.WriteString(fmt.Sprintf("%s", e.Name))
			case fmtHost:
				e.WriteString(fmt.Sprintf("%s", e.Hostname))
			case fmtSource:
				e.WriteString(fmt.Sprintf("%s", e.Source))
			case fmtMessage:
				e.WriteString(fmt.Sprintf("%s", e.Message))
			case fmtExtra:
				s := strings.Join(e.Extra, ",")
				e.WriteString(fmt.Sprintf("%s", s))
			default:
				// This skips fmt keywords, which could be useful.
				e.WriteString(key)
			}
		} else {
			e.WriteByte(f[0])
			f = f[1:]
		}
	}
	e.WriteString("\n")
	return e.String()
}

// parseKeyword returns the parsed keyword and the rest of the input string.
func (e *Event) parseKeyword(f string) (string, string) {
	var b strings.Builder
	if len(f) == 0 {
		return "", f
	}

	b.WriteByte(f[0])
	in := f[1:]
	loop := true
	for len(in) > 0 && loop {
		if !unicode.IsLetter(rune(in[0])) {
			loop = false
		} else {
			b.WriteByte(in[0])
			in = in[1:]
		}
	}
	return b.String(), in
}
