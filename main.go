package main

import (
	"errors"
	"os"

	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
)

// Options for the app.
var Options struct {
	//Version string is returned.
	Version bool `opt:"" short:"V" long:"version" help:"Display the version string and exit."`
	opt.DefaultHelp
	// Config defaults to "config.json" in the same directory.
	Config   string      `opt:"required" short:"c" help:"The configuration file." default:"config.json" group:"Basics" placeholder:"FILE"`
	Start    StartCmd    `command:"start" help:"Start server."`
	Generate GenerateCmd `command:"generate" help:"Scaffolding, structure and utility function generation." aliases:"gen"`
}

// StartCmd options.
type StartCmd struct {
	Help   bool   `short:"h" long:"help" help:"Show usage."`
	Domain string `help:"The domain to accept connections on." placeholder:"DOMAIN"`
}

// Run Info.
func (start *StartCmd) Run(in []string) error {
	if start.Help {
		return errors.New(opt.ErrorUsage)
	}

	log.Default.Msg("Domain: '%s'", start.Domain)
	log.Default.Msg("Remaining: %v", in)
	return nil
}

func main() {
	a := opt.Parse(&Options)
	if Options.Help || len(os.Args) < 2 {
		a.Usage()
		return
	}

	var err error

	err = a.RunCommand(false)
	if err != nil {
		log.Default.Msg("Error running: %s", err.Error())
		os.Exit(2)
	}
}
