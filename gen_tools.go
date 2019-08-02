package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Urethramancer/signor/files"
	"github.com/Urethramancer/signor/opt"
	"github.com/Urethramancer/signor/stringer"
)

type CmdGenTools struct {
	opt.DefaultHelp
	Index    string   `short:"i" long:"index" help:"Generate 'index' top-level command to list the supplied commands." placeholder:"NAME"`
	Main     string   `short:"m" long:"main" help:"Generate the option parser call code in its own file"`
	Output   string   `short:"o" long:"output" help:"Directory to save output files in. Current directory will be used if not specified." placeholder:"DIR" default:"cmd"`
	Package  string   `short:"p" long:"package" help:"Package name." placeholder:"NAME" default:"cmd"`
	Commands []string `help:"Command to generate a stub for. Aliases may be specified in the format 'command=alias1,alias2'." placeholder:"COMMAND"`
}

type cmdList struct {
	Name    string
	Aliases string
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
	cmd.Commands = stringer.RemoveDuplicateStringsAndTitle(cmd.Commands)
	out := stringer.New()
	var commands []cmdList
	for _, x := range cmd.Commands {
		c, l := splitCommandAliases(x)
		cl := cmdList{c, l}
		commands = append(commands, cl)
	}

	if cmd.Index != "" {
		_, err = out.WriteStrings("package ", cmd.Package, "\n\n")
		if err != nil {
			return err
		}

		err = generateIndex(out, cmd.Index, commands)
		if err != nil {
			return err
		}

		err = saveSource(out, cmd.Output, cmd.Index)
		if err != nil {
			return err
		}

		out.Reset()
	}

	for _, x := range commands {
		_, err = out.WriteStrings("package ", cmd.Package, "\n\n")
		if err != nil {
			return err
		}

		err = generateToolCommand(out, x)
		if err != nil {
			return err
		}

		err = saveSource(out, cmd.Output, x.Name)
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

func generateIndex(s *stringer.Stringer, name string, commands []cmdList) error {
	name = strings.ToLower(name)
	name = strings.Title(name)
	_, err := s.WriteStrings(
		"import (", "\n",
		"\t\"errors\"", "\n\n",
		"\t\"github.com/Urethramancer/signor/opt\"", "\n", ")\n\n",
		"// Cmd", name, " subcommands.\n",
		"type Cmd", name, " struct {\n", "\topt.DefaultHelp\n")
	if err != nil {
		return err
	}

	for _, x := range commands {
		_, err = s.WriteStrings("\t", x.Name, "\tCmd", x.Name, "\t`command:\"", strings.ToLower(x.Name), "\" help:\"<command help>\"")
		if err != nil {
			return err
		}

		if x.Aliases != "" {
			_, err = s.WriteStrings(" aliases:\"", x.Aliases, "\"")
			if err != nil {
				return err
			}
		}
		_, err = s.WriteString("`\n")
		if err != nil {
			return err
		}
	}

	_, err = s.WriteString("}\n\n")
	if err != nil {
		return err
	}

	return generateRun(s, name)
}

func generateRun(s *stringer.Stringer, name string) error {
	_, err := s.WriteStrings(
		"// Run ", strings.ToLower(name), "\n",
		"func (cmd *Cmd", name, ") Run(in []string) error {\n",
		"\tif cmd.Help {\n",
		"\t\treturn errors.New(opt.ErrorUsage)\n\t}\n\n",
		"\treturn nil\n",
		"}\n\n")
	return err
}

func generateToolCommand(s *stringer.Stringer, cmd cmdList) error {
	_, err := s.WriteStrings(
		"import (", "\n",
		"\t\"errors\"", "\n\n",
		"\t\"github.com/Urethramancer/signor/opt\"", "\n", ")\n\n",
		"// Cmd", cmd.Name, " options.\n",
		"type Cmd", cmd.Name, " struct {\n", "\topt.DefaultHelp\n}\n\n")
	if err != nil {
		return err
	}

	return generateRun(s, cmd.Name)
}

func splitCommandAliases(command string) (string, string) {
	a := strings.SplitN(command, "=", 2)
	cmd := a[0]
	if len(a) == 1 {
		return cmd, ""
	}

	return cmd, a[1]
}
