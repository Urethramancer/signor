package opt

import (
	"os"
	"reflect"
	"strings"

	"github.com/Urethramancer/signor/log"
	"github.com/Urethramancer/signor/stringer"
)

// Args gets options and commands parsed into it.
type Args struct {
	st             reflect.Value
	Program        string
	short          map[string]*Flag
	long           map[string]*Flag
	commands       map[string]*Flag
	commandlist    []*Flag
	positionalList []*Flag
	groups         map[string][]*Flag
	// groupOrder is in the order of group tags encountered.
	groupOrder []string
	Remaining  []string
	execute    *Flag
}

const (
	noGroup = "none"
)

// Usage printout.
func (a *Args) Usage() {
	var b stringer.Stringer
	b.WriteStrings("Usage:\n  ", os.Args[0])
	// Invocation
	c := len(a.short) + len(a.long)
	if c > 0 {
		if c > 1 {
			b.WriteString(" [OPTION]...")
		} else {
			b.WriteString(" [OPTION]")
		}
	}

	if a.execute != nil {
		b.WriteStrings(" ", a.execute.Name)
	}
	if len(a.commandlist) > 0 {
		b.WriteString(" [COMMAND]")
	}

	for _, p := range a.positionalList {
		b.WriteStrings(" [", p.Placeholder, "]")
		if p.IsSlice {
			b.WriteString("...")
		}
	}

	b.WriteString("\n")

	// Groups
	for _, gn := range a.groupOrder {
		flags := a.groups[gn]
		if gn == noGroup {
			if len(flags) > 0 {
				b.WriteString("\nApplication options:\n")
			}
		} else {
			b.WriteStrings("\n", gn, ":\n")
		}
		for _, f := range flags {
			fullFieldUsage(&b, f)
		}
	}

	// Positional arguments
	if len(a.positionalList) > 0 {
		b.WriteString("\nPositional arguments:\n")
		for _, f := range a.positionalList {
			fullFieldUsage(&b, f)
		}
	}

	// Commands
	if len(a.commandlist) > 0 {
		b.WriteString("\nCommands:\n")
		for _, f := range a.commandlist {
			fullFieldUsage(&b, f)
		}
	}

	log.Default.Msg(b.String())
}

func fullFieldUsage(b *stringer.Stringer, f *Flag) {
	vars, help := f.UsageString()
	b.WriteStrings(vars, "\t\t")
	if len(vars) < 16 {
		b.WriteString("\t")
	}
	if len(vars) < 8 {
		b.WriteString("\t")
	}
	b.WriteStrings(help, "\n")
}

// Parse the command line for arguments and tool commands.
func Parse(data interface{}) *Args {
	args := newArgs(os.Args)
	args.Parse(data, os.Args[1:])
	args.Program = os.Args[0]
	return args
}

func newArgs(in []string) *Args {
	a := Args{
		short:      make(map[string]*Flag),
		long:       make(map[string]*Flag),
		commands:   make(map[string]*Flag),
		groups:     make(map[string][]*Flag),
		groupOrder: []string{noGroup},
	}
	a.groups[noGroup] = make([]*Flag, 0)
	return &a
}

// Parse an option structure and slice of arguments.
func (a *Args) Parse(data interface{}, in []string) {
	a.parseOpts(data)
	a.parseArgs(in)
}

//Parse available options.
func (a *Args) parseOpts(data interface{}) {
	a.st = reflect.ValueOf(data).Elem()
	t := a.st.Type()
	for i := 0; i < a.st.NumField(); i++ {
		f := t.Field(i).Type
		if t.Field(i).Anonymous && f.Kind() == reflect.Struct && f.Field(0).Tag.Get("command") == "" {
			a.parseField(f.Field(0))
		} else {
			a.parseField(t.Field(i))
		}
	}
}

func (a *Args) parseField(sf reflect.StructField) {
	field := a.st.FieldByName(sf.Name)

	if !field.IsValid() {
		return
	}

	f := &Flag{
		field:       field,
		Name:        sf.Name,
		Help:        sf.Tag.Get("help"),
		Short:       sf.Tag.Get("short"),
		Long:        sf.Tag.Get("long"),
		Group:       sf.Tag.Get("group"),
		Placeholder: sf.Tag.Get("placeholder"),
		CommandName: sf.Tag.Get("command"),
		Default:     sf.Tag.Get("default"),
	}

	switch f.field.Kind() {
	case reflect.Slice:
		f.IsSlice = true
	case reflect.Map:
		f.IsMap = true
	default:
	}

	f.IsCommand = f.CommandName != ""

	c := sf.Tag.Get("choices")
	if c != "" {
		f.Choices = strings.Split(c, ",")
	}
	for i, c := range f.Choices {
		f.Choices[i] = strings.TrimSpace(c)
	}

	// Get boolean options
	f.parseOpts(sf.Tag.Get("opt"))

	if f.IsCommand {
		a.commandlist = append(a.commandlist, f)
		c = sf.Tag.Get("aliases")
		if c != "" {
			f.Aliases = strings.Split(c, ",")
		}
		for i, c := range f.Aliases {
			f.Aliases[i] = strings.TrimSpace(c)
		}
		a.commands[f.CommandName] = f
		for _, x := range f.Aliases {
			a.commands[x] = f
		}
	} else {
		var g []*Flag
		var ok bool
		if f.Long == "" && f.Short == "" && f.Placeholder != "" {
			a.positionalList = append(a.positionalList, f)
		} else {
			if f.Group == "" {
				g = a.groups[noGroup]
				g = append(g, f)
				a.groups[noGroup] = g
			} else {
				g, ok = a.groups[f.Group]
				if !ok {
					g = make([]*Flag, 0)
					a.groupOrder = append(a.groupOrder, f.Group)
				}
				g = append(g, f)
				a.groups[f.Group] = g
			}

			if f.Short != "" {
				a.short[f.Short] = f
			}
			if f.Long != "" {
				a.long[f.Long] = f
			}
		}
	}

	return
}

// parseArgs from CLI.
func (a *Args) parseArgs(args []string) {
	posDone := a.positionalList
	for i := 0; i < len(args); i++ {
		x := args[i]
		// We're done here - the "--" argument means to stop parsing
		if x == "--" {
			a.Remaining = args[i+1:]
			return
		}
		if strings.HasPrefix(x, "--") || strings.HasPrefix(x, "-") {
			if strings.HasPrefix(x, "--") {
				a.parseLong(args[i:])
			} else {
				args = a.parseShort(args[i:])
			}
		} else {
			if len(posDone) > 0 {
				p := posDone[0]
				posDone = posDone[1:]
				if p.IsSlice {
					p.field.Set(reflect.ValueOf(args[i:]))
					return
				}
				p.setValue(args[i])
			} else {
				f := a.commands[args[i]]
				if f != nil {
					f.parseCommand(args[i+1:])
					a.execute = f
					return
				}
				a.Remaining = append(a.Remaining, args[i])
			}
		}
	}
}

func (a *Args) parseLong(args []string) {
	n := args[0][2:]
	f, ok := a.long[n]
	if !ok {
		return
	}

	a.parseArg(args, f)
}

func (a *Args) parseShort(args []string) []string {
	flags := args[0][1:]
	for _, c := range flags {
		f := a.short[string(c)]
		if f != nil {
			if f.field.Kind() == reflect.Bool {
				f.setBool(true)
			} else {
				return a.parseArg(args, f)
			}
		}
	}
	return args
}

func isValidChoice(s string, choices []string) bool {
	if len(choices) == 0 {
		return true
	}

	for _, c := range choices {
		if s == c {
			return true
		}
	}

	return false
}

func (a *Args) parseArg(args []string, f *Flag) []string {
	if len(f.Choices) > 0 && !isValidChoice(args[1], f.Choices) {
		def := f.Choices[0]
		if f.Default != "" {
			def = f.Default
		}
		switch f.field.Kind() {
		case reflect.Int:
			f.setInt(def)
		case reflect.String:
			f.setString(def)
		}
		return args[1:]
	}

	if f.field.Kind() == reflect.Bool {
		f.setBool(true)
		return args[1:]
	}

	// Last arg is malformed, so we're done
	if len(args) < 2 {
		return args
	}

	f.setValue(args[1])
	return args[2:]
}
