//go:generate echo Generating stuff.
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
	Help    bool `short:"h" long:"help" help:"Show usage."`
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

	err := a.RunCommand()
	if err != nil {
		log.Default.Msg("Error running: %s", err.Error())
		os.Exit(2)
	}

	// srv := server.New("Test")
	// srv.Start()
	// <-daemon.BreakChannel()
	// srv.Stop()

	// var err error
	// var pkg structure.Package
	// pkg, err = structure.NewPackage("server/site.go", "structure/package.go")
	// if err != nil {
	// 	slog.TError("Error loading: %s", err.Error())
	// 	return
	// }

	// pkg.MakeTags(true, true)
	// slog.Msg("%s", pkg.String())
	// slog.Msg("%s", pkg.ProtoString())

	// err = a.RunCommand()
	// if err != nil {
	// 	slog.Msg("Error: %s", err.Error())
	// }

	// e := &log.Event{
	// 	Level:    4,
	// 	PID:      1242,
	// 	Time:     time.Now(),
	// 	Name:     "disk monitor",
	// 	Hostname: "localhost",
	// 	Source:   "/dev/md0",
	// 	Message:  "Disk full. Less than 1MB free.",
	// 	Extra:    nil,
	// }

	// l := log.NewLogger()
	// l.Log(e)
}
