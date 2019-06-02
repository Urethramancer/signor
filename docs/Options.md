# Command line options and flags

For more detailed option usage, see the [tag document](Tags.md).

## How it works
The `opt` package is a fairly GNU-like command line option parser with additional tool command functionality. It interprets options as encountered, then switches context to a new set of options for each tool command encountered. Single-character options can be merged into a string, and parsing moves on to the next long/short/tool option string when a non-flag (non-boolean) option is encountered.

The philosophy of this package is to give the API user control. Some things which differ from other packages of its ilk:
- The `Usage()` command will need to be called manually. This leaves you free to implement the "-h" and "--help" options for whatever else in tool commands.
- Every tool command with a Run() implementation that is supplied on the command line is executed. If a command is just a holder of sub-commands, don't implement the interface. But the option is there to run early tests or other functionality you feel is necessary.

## Using it

A minimal example of use:
```go
var Options struct {
	Version bool `opt:"" short:"V" long:"version" help:"Display the version string and exit."`
	Help    bool `short:"h" long:"help" help:"Show usage."`
}

func main() {
	a := opt.Parse(&Options)
	if Options.Help || len(os.Args) < 2 {
		a.Usage()
		return
	}
}
```

This exits if the `-h` or `--help` flag is specified, showing pretty-printed usage options.

A tool command variant:
```go
var Options struct {
	Version bool `opt:"" short:"V" long:"version" help:"Display the version string and exit."`
	Help HelpCmd `short:"c" help:"Show help on a subject."`
}

type Help struct {
	Subject string `help:"The subject." placeholder:"SUBJECT"`
}

```

### Positional options

Options can be tagged only with placeholder and help tags, which will make them positional arguments.
