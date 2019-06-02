package opt

import (
	"errors"
)

// Runner is the interface for tool commands to conform to.
type Runner interface {
	Run(args []string) error
}

// RunCommand and recurse.
func (a *Args) RunCommand() error {
	if a.execute == nil {
		return errors.New(ErrorNoCommand)
	}

	err := a.execute.executeCommand()
	if err != nil && err.Error() == ErrorUsage {
		a.execute.Args.Usage()
		// Swallow this error message
		return nil
	}

	return err
}
