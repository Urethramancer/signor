package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

type GenerateCmd struct {
	opt.DefaultHelp
	Travis TravisCmd `command:"travis" help:"Generate .travis.yml files for GitHub projects."`
	REST   RESTCmd   `command:"rest" help:"Generate REST code from tagged structures."`
}

func (cmd *GenerateCmd) Run(in []string) error {
	return errors.New(opt.ErrorUsage)
}
