package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/Urethramancer/signor/files"
	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
	"github.com/Urethramancer/signor/stringer"
	"github.com/Urethramancer/signor/structure"
)

type TravisCmd struct {
	Help  bool     `short:"h" long:"help" help:"Show usage."`
	Name  string   `short:"o" long:"output" help:"Filename to save the YAML file as." placeholder:"FILE"`
	Input []string `help:"Input Go source file to read imports from." placeholder:"SOURCE"`
}

func (tr *TravisCmd) Run(in []string) error {
	if tr.Help {
		return errors.New(opt.ErrorUsage)
	}

	ver, err := goversion()
	if err != nil {
		log.Default.Err("Couldn't run Go: %s", err.Error())
		os.Exit(2)
	}

	yml := stringer.New()
	yml.WriteStrings("language: go\n\ngo:\n  - ", ver, "\n")

	if tr.Input == nil || len(tr.Input) < 0 {
		return errors.New("not enough arguments")
	}

	pkg, err := structure.NewPackage(tr.Input...)
	if err != nil {
		return err
	}

	yml.WriteStrings("\ninstall:\n")
	for _, imp := range pkg.MergeExternalImports() {
		imp = strings.ReplaceAll(imp, "\"", "")
		yml.WriteStrings("    - go get ", imp, "\n")
	}

	yml.WriteString("\ninclude:\n  - os: linux\n")
	yml.WriteI("    go: ", '"', ver, ".x", '"', "\n")
	yml.WriteStrings(
		"    cache:\n      directories:\n",
		"        - $HOME/.cache/go-build\n",
		"        - $HOME/gopath/pkg/mod\n",
	)

	if tr.Name == "" {
		log.Default.Msg(yml.String())
	} else {
		return files.WriteFile(tr.Name, []byte(yml.String()))
	}
	return nil
}

func goversion() (string, error) {
	cmd := exec.Command("go", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	a := strings.Split(string(out), " ")
	if len(a) < 3 || !strings.HasPrefix(a[2], "go") {
		return "1.12", nil
	}
	in := strings.Split(a[2][2:], ".")
	in = in[:len(in)-1]
	ver := strings.Join(in, ".")
	return ver, nil
}
