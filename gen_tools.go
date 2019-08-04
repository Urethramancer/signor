package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mgutz/str"

	"github.com/Urethramancer/signor/files"
	"github.com/Urethramancer/signor/opt"
	"github.com/Urethramancer/signor/stringer"
)

type CmdGenTools struct {
	opt.DefaultHelp
	Index    string   `short:"i" long:"index" help:"Generate 'index' top-level command to list the supplied commands." placeholder:"NAME"`
	Main     string   `short:"m" long:"main" help:"Generate the option parser call code in its own file" placeholder:"FILENAME"`
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

		err = saveSource(out, cmd.Output, "cmd_", cmd.Index)
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

		err = generateToolCommand(out, cmd.Index, x)
		if err != nil {
			return err
		}

		err = saveSource(out, cmd.Output, cmd.Index, x.Name)
		if err != nil {
			return err
		}

		out.Reset()
	}

	if cmd.Main != "" {
		_, err = out.WriteStrings("package ", cmd.Package, "\n\n")
		if err != nil {
			return err
		}

		name := strings.ToLower(cmd.Main)
		name = str.ChompRight(name, ".go")
		name = strings.Title(name)
		err = generateMain(out, name)
		if err != nil {
			return err
		}

		name = filepath.Join(cmd.Output, cmd.Main)
		err = files.WriteFile(name, []byte(out.String()))
		if err != nil {
			return err
		}
	}

	return nil
}

func saveSource(s *stringer.Stringer, dir, sub, fn string) error {
	if sub == "" {
		fn = fmt.Sprintf("cmd_%s.go", strings.ToLower(fn))
	} else {
		fn = fmt.Sprintf("%s_%s.go", sub, strings.ToLower(fn))
	}
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
		_, err = s.WriteStrings("\t", x.Name, "\tCmd", name, x.Name, "\t`command:\"", strings.ToLower(x.Name), "\" help:\"<command help>\"")
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

	return generateRun(s, "", name)
}

func generateRun(s *stringer.Stringer, sub, name string) error {
	sub = strings.ToLower(sub)
	sub = strings.Title(sub)
	_, err := s.WriteStrings(
		"// Run ", strings.ToLower(name), "\n",
		"func (cmd *Cmd", sub, name, ") Run(in []string) error {\n",
		"\tif cmd.Help {\n",
		"\t\treturn errors.New(opt.ErrorUsage)\n\t}\n\n",
		"\treturn nil\n",
		"}\n\n")
	return err
}

func generateToolCommand(s *stringer.Stringer, sub string, cmd cmdList) error {
	sub = strings.ToLower(sub)
	sub = strings.Title(sub)
	_, err := s.WriteStrings(
		"import (", "\n",
		"\t\"errors\"", "\n\n",
		"\t\"github.com/Urethramancer/signor/opt\"", "\n", ")\n\n",
		"// Cmd", sub, cmd.Name, " options.\n",
		"type Cmd", sub, cmd.Name, " struct {\n", "\topt.DefaultHelp\n}\n\n")
	if err != nil {
		return err
	}

	return generateRun(s, sub, cmd.Name)
}

func splitCommandAliases(command string) (string, string) {
	a := strings.SplitN(command, "=", 2)
	cmd := a[0]
	if len(a) == 1 {
		return cmd, ""
	}

	return cmd, a[1]
}

func generateMain(s *stringer.Stringer, command string) error {
	s.WriteStrings(
		optHeaderTemplate,
		"// Options holds all the tool commands.\n",
		"var Options struct {\n",
		"\topt.DefaultHelp\n",
	)
	s.WriteStrings(
		"\t", command, "\tCmd", command,
		"`command:\"", strings.ToLower(command), "\" help:\"", "<insert help here>\"`\n",
	)
	s.WriteStrings("}\n\n", optTemplate)
	return nil
}

var optHeaderTemplate = `import (
	"os"

	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
)

`

var optTemplate = `// ParseOptions could be renamed to main(), or called from main.
func ParseOptions() {
	a := opt.Parse(&Options)
	if Options.Help || len(os.Args) < 2 {
		a.Usage()
		return
	}

	var err error
	err = a.RunCommand(false)
	if err != nil {
		log.Default.Msg("Error running: %s", err.Error())
		os.Exit(2)
	}
}
`
