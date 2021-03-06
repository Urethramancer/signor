package main

import (
	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
	"github.com/Urethramancer/signor/structure"
)

// CmdDB generates database access functions.
type CmdDB struct {
	// Help flag.
	Help bool `short:"h" long:"help" help:"Show usage."`

	// Package name.
	Package string `short:"p" long:"package" help:"Package name." placeholder:"STRING" default:"database"`
	// User table flag.
	User bool `short:"u" long:"user" help:"Set this as a user table for authentication, and generate special code."`
	// Type of database. Only Postgres is supported for now.
	Type string `short:"t" long:"type" help:"Type of database." choices:"pg"`

	// Input is the source file to generate from.
	Input string `help:"Input file with structures for database tables." placeholder:"INFILE"`
	// Output is the directory to generate output files in.
	Output string `help:"Output directory." placeholder:"PATH"`
}

// Run the DB code generator.
func (cmd *CmdDB) Run(in []string) error {
	if cmd.Help || cmd.Output == "" {
		return opt.ErrUsage
	}

	m := log.Default.Msg
	pkg, err := structure.NewPackage(cmd.Input)
	if err != nil {
		return err
	}

	pkg.Name = cmd.Package

	m("%s", pkg.String())
	return nil
}
