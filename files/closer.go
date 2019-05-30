package files

import (
	"os"

	"github.com/Urethramancer/signor/log"
)

// Closer holds open files to close them all in sequence, logging any errors.
type Closer struct {
	list []*os.File
	l    *log.Logger
}

// NewCloser returns a pointer to a closer.
func NewCloser(files ...*os.File) *Closer {
	c := &Closer{
		l: log.Default,
	}
	for _, f := range files {
		c.AddFile(f)
	}
	return c
}

// SetLogger to an alternative logger.
func (c *Closer) SetLogger(l *log.Logger) {
	c.l = l
}

// AddFile adds a file pointer to the closer's list.
func (c *Closer) AddFile(f *os.File) *Closer {
	c.list = append(c.list, f)
	return c
}

// Close and remove all files in the list, and log any errors from the Close() call with or without a timestamp.
func (c *Closer) Close(t bool) {
	if len(c.list) == 0 {
		return
	}

	var err error
	for _, f := range c.list {
		name := f.Name()
		err = f.Close()
		if err != nil {
			if t {
				c.l.TErr("Error closing %s: %s", name, err.Error())
			} else {
				c.l.Err("Error closing %s: %s", name, err.Error())
			}
		}
	}

	// Clear the list for reuse.
	c.list = []*os.File{}
}
