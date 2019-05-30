package opt

// Runner is the interface for tool commands to conform to.
type Runner interface {
	Run(args []string) error
}

// RunCommand and recurse.
func (a *Args) RunCommand() error {
	if a.execute != nil {
		err := a.execute.executeCommand()
		if err != nil {
			return err
		}
	}

	return nil
}
