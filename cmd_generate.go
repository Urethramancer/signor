package main

import (
	"errors"

	"github.com/Urethramancer/signor/opt"
)

type GenerateCmd struct {
	Travis TravisCmd `command:"travis" help:"Generate .travis.yml files for GitHub projects."`
}

func (g *GenerateCmd) Run(in []string) error {
	return errors.New(opt.ErrorUsage)
}
