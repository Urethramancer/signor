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
	ConfigPath string `opt:"required" short:"c" help:"The configuration file." default:"config.json" group:"Basics" placeholder:"FILE"`

	// Tool commands
	Travis CmdTravis `command:"travis" help:"Generate .travis.yml files for GitHub projects."`
	REST   CmdREST   `command:"rest" help:"Generate REST code from tagged structures."`
	DB     CmdDB     `command:"database" help:"Generate database schema from tagged structures." aliases:"db"`
	Config CmdConfig `command:"config" aliases:"cfg" help:"Generate database schema from tagged structures."`
	Tools  CmdTools  `command:"tools" help:"Generate tool command stubs."`
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
