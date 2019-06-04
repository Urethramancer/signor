package opt

import (
	"errors"
)

// Runner is the interface for tool commands to conform to.
type Runner interface {
	Run(args []string) error
}

// RunCommand and recurse.
func (a *Args) RunCommand(all bool) error {
	if a.execute == nil {
		return errors.New(ErrorNoCommand)
	}

	err := a.execute.executeCommand(all)
	if err != nil && err.Error() == ErrorUsage {
		a.execute.Args.Usage()
		// Swallow this error message
		return nil
	}

	return err
}
