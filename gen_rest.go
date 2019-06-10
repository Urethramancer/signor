package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

// RESTCmd generator struct.
type RESTCmd struct {
	// Help flag.
	Help bool `short:"h" long:"help" help:"Show usage."`
	// Input is the source file to generate from.
	Input string `help:"Input file with structures for REST endpoints." placeholder:"INFILE"`
	// Output is the Go file to generate.
	Output string `help:"Output Go source file." placeholder:"OUTFILE"`
}

func (cmd *RESTCmd) Run(in []string) error {
	return errors.New(opt.ErrorUsage)
}
