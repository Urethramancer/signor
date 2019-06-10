package main

import (
	"os"

	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
)

// Options for the app.
var Options struct {
	opt.DefaultHelp
	//Version string is returned.
	Version bool `opt:"" short:"V" long:"version" help:"Display the version string and exit."`
	// Config defaults to "config.json" in the same directory.
	Config   string      `opt:"required" short:"c" help:"The configuration file." default:"config.json" group:"Basics" placeholder:"FILE"`
	Generate GenerateCmd `command:"generate" help:"Scaffolding, structure and utility function generation." aliases:"gen"`
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
