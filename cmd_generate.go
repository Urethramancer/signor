package main

type GenerateCmd struct {
	Travis TravisCmd `command:"travis" help:"Generate .travis.yml files for GitHub projects."`
}
