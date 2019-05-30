//go:generate echo Generating stuff.
package main

import (
	"os"

	"github.com/Urethramancer/signor/opt"
	"github.com/Urethramancer/signor/server"
	"github.com/Urethramancer/slog"
)

// Options for the app.
var Options struct {
	//Version string is returned..
	Version bool `opt:"" short:"V" long:"version" help:"Display the version string and exit." group:"Basics"`
	Help    bool `short:"h" long:"help" help:"Show usage."`
	// Config defaults to "config.json" in the same directory.
	Config string `opt:"required" short:"c" help:"The configuration file." default:"config.json" group:"Basics" placeholder:"FILE"`
	Name   string `opt:"required" long:"name"`
	// Number is an integer.
	Number      int            `short:"n" group:"Maths" help:"An integer."`
	OneTwoThree int            `short:"o" help:"One, two or three." choices:"1,2,3" default:"1"`
	Colour      string         `short:"C" long:"colour" help:"A colour." choices:"red,green,blue" default:"green"`
	List        List           `command:"list" group:"Commands" help:"List stuff." aliases:"ls,l"`
	X           []string       `short:"X"`
	Y           []int          `short:"Y"`
	Z           map[string]int `short:"Z"`
	A           bool           `short:"a"`
	B           bool           `short:"b"`
}

// List options and sub-commands.
type List struct {
	All     All  `command:"all" group:"Commands"`
	Verbose bool `short:"v" help:"List more details."`
}

// Run List.
func (l *List) Run(in []string) error {
	slog.Msg("Running list")
	slog.Msg("Remaining args: %+v", in)
	return nil
}

// All has no options.
type All struct{}

// Run List All.
func (a *All) Run(in []string) error {
	slog.Msg("That's all, folks.")
	slog.Msg("Remaining args: %+v", in)
	return nil
}

func main() {
	a := opt.Parse(&Options)
	if Options.Help || len(os.Args) < 2 {
		a.Usage()
		return
	}

	srv := server.New("Test")
	srv.Start()
	// <-daemon.BreakChannel()
	srv.Stop()

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
