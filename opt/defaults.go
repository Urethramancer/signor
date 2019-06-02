package opt

// DefaultHelp can be embedded in your options struct to save some typing.
type DefaultHelp struct {
	Help bool `short:"h" help:"Show this help."`
}
