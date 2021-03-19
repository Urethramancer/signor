package opt

import "errors"

var (
	ErrUsage     = errors.New("unknown options")
	ErrNoCommand = errors.New("no command specified")
)
