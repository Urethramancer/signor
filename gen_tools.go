package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Urethramancer/signor/files"

	"github.com/Urethramancer/signor/stringer"

	"github.com/Urethramancer/signor/opt"
)

type CmdGenTools struct {
	opt.DefaultHelp
	Index    string   `short:"i" long:"index" help:"Generate 'index' top-level command to list the supplied commands."`
	Output   string   `short:"o" long:"output" help:"Directory to save output files in. Current directory will be used if not specified." placeholder:"DIR" default:"cmd"`
	Package  string   `short:"p" long:"package" help:"Package name." placeholder:"NAME" default:"cmd"`
	Commands []string `help:"Command to generate a stub for." placeholder:"COMMAND"`
}

func (cmd *CmdGenTools) Run(in []string) error {
	if cmd.Help || len(cmd.Commands) == 0 {
		return errors.New(opt.ErrorUsage)
	}

	err := os.MkdirAll(cmd.Output, 0755)
	if err != nil {
		return err
	}

	sort.Strings(cmd.Commands)
	cmd.Commands = prepToolList(cmd.Commands)
	out := stringer.New()
	if cmd.Index != "" {
		out.WriteStrings("package ", cmd.Package, "\n\n")
		generateIndex(out, cmd.Index, cmd.Commands)
		err = saveSource(out, cmd.Output, cmd.Index)
		if err != nil {
			return err
		}

		out.Reset()
	}

	for _, x := range cmd.Commands {
		out.WriteStrings("package ", cmd.Package, "\n\n")
		generateToolCommand(out, x)
		err = saveSource(out, cmd.Output, x)
		if err != nil {
			return err
		}

		out.Reset()
	}
	return nil
}

func saveSource(s *stringer.Stringer, dir, fn string) error {
	fn = fmt.Sprintf("cmd_%s.go", strings.ToLower(fn))
	fn = filepath.Join(dir, fn)
	return files.WriteFile(fn, []byte(s.String()))
}

func generateIndex(s *stringer.Stringer, name string, commands []string) {
	s.WriteStrings(
		"import (", "\n",
		"\t\"errors\"", "\n\n",
		"\t\"github.com/Urethramancer/signor/opt\"", "\n", ")\n\n",
	)

	name = strings.ToLower(name)
	name = strings.Title(name)
	s.WriteStrings("// Cmd", name, " subcommands.\n")
	s.WriteStrings("type Cmd", name, " struct {\n", "\topt.DefaultHelp\n")
	for _, x := range commands {
		s.WriteStrings("\t", x, "\tCmd", x, "\t`command:\"", strings.ToLower(x), "\" help:\"<command help>\"`", "\n")
	}
	s.WriteString("}\n\n")
	generateRun(s, name)
}

func generateRun(s *stringer.Stringer, name string) {
	s.WriteStrings("// Run ", strings.ToLower(name), "\n")
	s.WriteStrings("func (cmd *Cmd", name, ") Run(in []string) error {\n")
	s.WriteString("\tif cmd.Help {\n")
	s.WriteString("\t\treturn errors.New(opt.ErrorUsage)\n\t}\n\n")
	s.WriteString("\treturn nil\n")
	s.WriteString("}\n\n")
}

func generateToolCommand(s *stringer.Stringer, name string) {
	s.WriteStrings(
		"import (", "\n",
		"\t\"errors\"", "\n\n",
		"\t\"github.com/Urethramancer/signor/opt\"", "\n", ")\n\n",
	)

	s.WriteStrings("// Cmd", name, " options.\n")
	s.WriteStrings("type Cmd", name, " struct {\n", "\topt.DefaultHelp\n}\n\n")
	generateRun(s, name)
}

func prepToolList(a []string) []string {
	m := map[string]bool{}
	if len(a) < 2 {
		return a
	}

	for _, x := range a {
		x = strings.ToLower(x)
		x = strings.Title(x)
		m[x] = true
	}

	l := []string{}
	for k := range m {
		l = append(l, k)
	}

	sort.Strings(l)
	return l
}
