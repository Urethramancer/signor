package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

type GenerateCmd struct {
	opt.DefaultHelp
	Travis TravisCmd `command:"travis" help:"Generate .travis.yml files for GitHub projects."`
}

func (g *GenerateCmd) Run(in []string) error {
	if g.Help {
		return errors.New(opt.ErrorUsage)
	}

	return nil
}
