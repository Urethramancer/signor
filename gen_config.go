package main

import (
	"errors"
	"strings"

	"github.com/Urethramancer/cross"
	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/opt"
	"github.com/Urethramancer/signor/stringer"
	"github.com/Urethramancer/signor/structure"
)

// CmdGenConfig generates configuration file loading, saving and tool commands.
type CmdGenConfig struct {
	opt.DefaultHelp
	Input  string `help:"Input Go source file to read imports from." placeholder:"SOURCE"`
	Output string `help:"Output path." placeholder:"PATH" default:"config"`
}

type SubStructs map[string]string

func (ss SubStructs) Add(k, v string) {
	ss[k] = v
}

func (ss SubStructs) Remove(k string) {
	delete(ss, k)
}

func (ss SubStructs) Has(k string) bool {
	_, ok := ss[k]
	return ok
}

type cfgFile struct {
	options []*cfgOption
}

func (c *cfgFile) Add(o *cfgOption) {
	c.options = append(c.options, o)
}

type cfgOption struct {
	Name string
	Type string
	Tag  string
}

func (cmd *CmdGenConfig) Run(in []string) error {
	if cmd.Help || cmd.Input == "" || cmd.Output == "" {
		return errors.New(opt.ErrorUsage)
	}

	if cross.DirExists(cmd.Output) {
		return errors.New("output directory already exists, aborting")
	}

	pkg, err := structure.NewPackage(cmd.Input)
	pkg.Name = cmd.Output
	if err != nil {
		return err
	}

	m := log.Default.Msg
	commands := stringer.New()
	handlers := stringer.New()
	embedded := make(SubStructs)
	stlist := []string{}
	stmap := make(map[string][]*cfgOption)
	for _, s := range pkg.Structs {
		var comment string
		list := []*cfgOption{}
		if embedded.Has(s.Name) {
			list = stmap[embedded[s.Name]]
		}
		for _, f := range s.Fields {
			if f.IsComment {
				comment = f.Name
			} else {
				switch f.Value {
				case "string", "float32", "float64", "int", "bool":
					f.MakeTags(true, false)
					opt := createOption(f, comment)
					list = append(list, opt)

				default:
					embedded.Add(f.Name, s.Name)
					f.MakeTags(true, false)
				}
			}
		}
		if embedded.Has(s.Name) {
			stmap[embedded[s.Name]] = list
		} else {
			stmap[s.Name] = list
			stlist = append(stlist, s.Name)
		}
	}

	printHeader(cmd.Output + ".go:")
	m("// Package %s loads and saves the %s structure.", pkg.Name)
	m("%s", pkg.String())
	printHeader("Config functions")

	printHeader("commands.go")
	for _, st := range stlist {
		commands.WriteStrings("type ", st, "GetCommands struct {\n")
		commands.WriteStrings(
			"\tGet", st,
			"\tCmdGet", st,
			"\t`",
			"command:\"", strings.ToLower(st), "\"",
			" help:\"Get configuration variables from ", st, ".", "\"",
			"`", "\n",
		)
		commands.WriteString("}\n\n")

		commands.WriteStrings("type ", st, "SetCommands struct {\n")
		commands.WriteStrings(
			"\tSet", st,
			"\tCmdSet", st,
			"\t`",
			"command:\"", strings.ToLower(st), "\"",
			" help:\"Set configuration variables from ", st, ".", "\"",
			"`", "\n",
		)
		commands.WriteString("}\n\n")

		commands.WriteStrings("// CmdGet", st, " options.\n")
		commands.WriteStrings("type CmdGet", st, " struct {\n")
		for _, opt := range stmap[st] {
			commands.WriteStrings("\t", opt.Name, "\tGet", opt.Name, "\t", opt.Tag, "\n")
		}
		commands.WriteStrings("}\n\n")

		commands.WriteStrings("// CmdSet", st, " options.\n")
		commands.WriteStrings("type CmdSet", st, " struct {\n")
		for _, opt := range stmap[st] {
			commands.WriteStrings("\t", opt.Name, "\tSet", opt.Name, "\t", opt.Tag, "\n")
		}
		commands.WriteStrings("}\n\n")
	}
	m("%s", commands.String())
	printHeader("handlers.go")
	m("%s", handlers.String())
	return nil
}

func printHeader(s string) {
	h := stringer.New()
	h.WriteStrings(
		strings.Repeat("*", len(s)+4),
		"\n* ", s, " *\n",
		strings.Repeat("*", len(s)+4),
	)
	log.Default.Msg("%s", h.String())
}

func createOption(f *structure.Field, comment string) *cfgOption {
	opt := &cfgOption{
		Type: f.Value,
	}
	t := stringer.New()
	f.Name = strings.ToLower(f.Name)
	t.WriteI("`command:\"", f.Name, "\"")
	f.Name = strings.Title(f.Name)
	opt.Name = f.Name
	if comment != "" {
		t.WriteI(" help:", "\"", strings.TrimSpace(comment[2:]), "\"")
	}
	t.WriteI(" placeholder:", "\"", strings.ToUpper(f.Name), "\"")
	opt.Tag = t.String()
	t.WriteString("`")
	return opt
}
