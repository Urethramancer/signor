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
	cmdGroups      map[string][]*Flag
	// groupOrder is in the order of group tags encountered.
	groupOrder    []string
	cmdGroupOrder []string
	Remaining     []string
	execute       *Flag
}

const (
	noGroup = "none"
)

// Usage printout.
func (a *Args) Usage() {
	var b stringer.Stringer
	b.WriteStrings("Usage:\n  ", a.Program)
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
	for _, gn := range a.cmdGroupOrder {
		flags := a.cmdGroups[gn]
		if gn == noGroup {
			if len(flags) > 0 {
				b.WriteString("\nCommands:\n")
			}
		} else {
			b.WriteStrings("\n", gn, ":\n")
		}
		for _, f := range flags {
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
	args.Parse(data, os.Args[1:], os.Args[0])
	return args
}

func newArgs(in []string) *Args {
	a := Args{
		short:         make(map[string]*Flag),
		long:          make(map[string]*Flag),
		commands:      make(map[string]*Flag),
		groups:        make(map[string][]*Flag),
		groupOrder:    []string{noGroup},
		cmdGroups:     make(map[string][]*Flag),
		cmdGroupOrder: []string{noGroup},
	}
	a.groups[noGroup] = make([]*Flag, 0)
	return &a
}

// Parse an option structure and slice of arguments.
func (a *Args) Parse(data interface{}, in []string, parent string) {
	a.Program = parent
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

	if f.Default != "" {
		f.setValue(f.Default)
	}

	switch f.field.Kind() {
	case reflect.Slice:
		f.IsSlice = true
	case reflect.Map:
		f.IsMap = true
	default:
	}

	f.IsCommand = f.CommandName != ""

	envvar := sf.Tag.Get("env")
	if envvar != "" {
		f.setValue(os.Getenv(envvar))
	}

	c := sf.Tag.Get("choices")
	if c != "" {
		f.Choices = strings.Split(c, ",")
	}
	for i, c := range f.Choices {
		f.Choices[i] = strings.ToLower(strings.TrimSpace(c))
	}

	// Get boolean options
	f.parseOpts(sf.Tag.Get("opt"))

	var g []*Flag
	var ok bool
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

		if f.Group == "" {
			g = a.cmdGroups[noGroup]
			g = append(g, f)
			a.cmdGroups[noGroup] = g
		} else {
			g, ok = a.cmdGroups[f.Group]
			if !ok {
				g = make([]*Flag, 0)
				a.cmdGroupOrder = append(a.cmdGroupOrder, f.Group)
			}
			g = append(g, f)
			a.cmdGroups[f.Group] = g
		}
	} else {
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
}

// parseArgs from CLI.
func (a *Args) parseArgs(args []string) {
	posDone := a.positionalList
	for len(args) > 0 {
		x := args[0]
		// We're done here - the "--" argument means to stop parsing
		if x == "--" {
			a.Remaining = args[1:]
			return
		}
		if strings.HasPrefix(x, "--") || strings.HasPrefix(x, "-") {
			if strings.HasPrefix(x, "--") {
				args = a.parseLong(args)
			} else {
				args = a.parseShort(args)
			}
		} else {
			if len(posDone) > 0 {
				p := posDone[0]
				posDone = posDone[1:]
				if p.IsSlice {
					p.field.Set(reflect.ValueOf(args))
					return
				}
				p.setValue(args[0])
			} else {
				f := a.commands[args[0]]
				if f != nil {
					f.parseCommand(args[1:], a.Program)
					a.execute = f
					return
				}
				a.Remaining = append(a.Remaining, args[0])
			}
			if len(args) > 0 {
				args = args[1:]
			}
		}
	}
}

func (a *Args) parseLong(args []string) []string {
	n := args[0][2:]
	if strings.Contains(n, "=") {
		sl := strings.Split(n, "=")
		n = sl[0]
		newargs := []string{sl[1]}
		args = append(newargs, args[1:]...)
	} else {
		args = args[1:]
	}
	f, ok := a.long[n]
	if !ok {
		return args
	}

	if f.field.Kind() == reflect.Bool {
		f.setBool(true)
		return args
	}

	return a.parseArg(args, f)
}

// parseShort sets any boolean flags encountered to true, and will
// parse the next argument if one of the options is a non-bool.
func (a *Args) parseShort(args []string) []string {
	flags := args[0][1:]
	for _, c := range flags {
		f := a.short[string(c)]
		if f != nil {
			if f.field.Kind() == reflect.Bool {
				f.setBool(true)
			} else {
				// We break off here, as non-bool options can only be the last one.
				return a.parseArg(args[1:], f)
			}
		}
	}
	return args[1:]
}

// isValidChoice checks if choice s is in list choices.
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
	if len(f.Choices) > 0 && len(args) > 1 && !isValidChoice(args[0], f.Choices) {
		def := f.Choices[0]
		if f.Default == "" {
			def = f.Default
		}
		switch f.field.Kind() {
		case reflect.Int:
			f.setInt(def)
		case reflect.String:
			f.setString(def)
		}
		if len(args) < 2 {
			return nil
		} else {
			return args[1:]
		}
	}

	if f.field.Kind() == reflect.Bool {
		f.setBool(true)
		return args[1:]
	}

	f.setValue(args[0])
	return args[1:]
}

// SetChoicesShort sets the selectable options based on the short option name.
func (a *Args) SetChoicesShort(name string, list []string) {
	sh, ok := a.short[name]
	if !ok {
		return
	}

	sh.Choices = list
}

// SetChoicesLong sets the selectable options based on the long option name.
func (a *Args) SetChoicesLong(name string, list []string) {
	a.long[name].Choices = list
}
