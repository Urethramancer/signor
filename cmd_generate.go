package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

type GenerateCmd struct {
	opt.DefaultHelp
	Travis TravisCmd `command:"travis" help:"Generate .travis.yml files for GitHub projects."`
	REST   RESTCmd   `command:"rest" help:"Generate REST code from tagged structures."`
	DB     DBCmd     `command:"database" help:"Generate database schema from tagged structures." aliases:"db"`
}

func (cmd *GenerateCmd) Run(in []string) error {
	return errors.New(opt.ErrorUsage)
}
