package main

import (
	"errors"

	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
	"github.com/Urethramancer/signor/stringer"
	"github.com/Urethramancer/signor/structure"
)

// CmdGenConfig generates configuration file loading, saving and tool commands.
type CmdGenConfig struct {
	opt.DefaultHelp
	Input  string `help:"Input Go source file to read imports from." placeholder:"SOURCE"`
	Output string `help:"Output path." placeholder:"PATH"`
}

type Set map[string]bool

func (l Set) Add(s string) {
	l[s] = true
}

func (l Set) Remove(s string) {
	delete(l, s)
}

func (l Set) Has(s string) bool {
	_, ok := l[s]
	return ok
}

func (cmd *CmdGenConfig) Run(in []string) error {
	if cmd.Help || cmd.Output == "" {
		return errors.New(opt.ErrorUsage)
	}

	cfg := stringer.New()

	pkg, err := structure.NewPackage(cmd.Input)
	if err != nil {
		return err
	}

	cfg.WriteString(pkg.String())
	log.Default.Msg("%s", cfg.String())
	return nil
}
