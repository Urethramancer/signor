package opt

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/Urethramancer/signor/stringer"
)

// Flag or command option.
type Flag struct {
	field       reflect.Value
	Name        string
	Help        string
	Short       string
	Long        string
	Group       string
	Placeholder string
	CommandName string
	Default     string
	Choices     []string
	Aliases     []string
	Args        *Args
	command     reflect.Value
	filter      reflect.Value
	err         error
	IsCommand   bool
	IsSlice     bool
	IsMap       bool
	Required    bool
}

func (f *Flag) UsageString() (string, string) {
	var vars, help stringer.Stringer
	vars.WriteString("  ")
	if f.Short != "" {
		vars.WriteStrings("-", f.Short)
	}

	if f.Long != "" {
		if f.Short != "" {
			vars.WriteString(", ")
		}
		vars.WriteStrings("--", f.Long)
	}

	if f.Placeholder != "" {
		vars.WriteStrings(" ", f.Placeholder)
	}

	help.WriteString(f.Help)

	if f.CommandName != "" {
		vars.WriteString(f.CommandName)
		if len(f.Aliases) > 0 {
			help.WriteStrings(
				" (Aliases: ",
				strings.Join(f.Aliases, ", "),
				")",
			)
		}
	}

	if len(f.Choices) > 0 {
		help.WriteStrings(
			" (Restricted to: ",
			strings.Join(f.Choices, ", "),
			")",
		)
	}

	if f.Default != "" {
		help.WriteStrings(" (Default: ", f.Default, ")")
	}
	return vars.String(), help.String()
}

func (f *Flag) parseOpts(opt string) {
	opts := strings.Split(opt, ",")
	if len(opts) == 0 {
		return
	}

	for _, o := range opts {
		switch o {
		case "required":
			f.Required = true
		}
	}
}

func (f *Flag) setValue(s string) {
	switch f.field.Kind() {
	case reflect.String:
		f.setString(s)
	case reflect.Int:
		f.setInt(s)
	case reflect.Float32:
		f.setFloat32(s)
	case reflect.Float64:
		f.setFloat64(s)
	case reflect.Slice:
		f.setSlice(s)
	case reflect.Map:
		f.setMap(s)
	}
}

// setBool from bool
func (f *Flag) setBool(b bool) {
	f.field.SetBool(b)
}

// setString from string
func (f *Flag) setString(s string) {
	f.field.SetString(s)
}

// setInt from string
func (f *Flag) setInt(s string) {
	n, _ := strconv.Atoi(s)
	f.field.SetInt(int64(n))
}

// setFloat32 from string
func (f *Flag) setFloat32(s string) {
	n, _ := strconv.ParseFloat(s, 32)
	f.field.SetFloat(n)
}

// setFloat64 from string
func (f *Flag) setFloat64(s string) {
	n, _ := strconv.ParseFloat(s, 64)
	f.field.SetFloat(n)
}

// setSlice from string, creating a new slice.
func (f *Flag) setSlice(s string) {
	a := strings.Split(s, ",")
	switch f.field.Type().Elem().Kind() {
	case reflect.String:
		f.field.Set(reflect.ValueOf(a))
	case reflect.Int:
		var ints []int
		for _, x := range a {
			n, _ := strconv.Atoi(x)
			ints = append(ints, n)
		}
		f.field.Set(reflect.ValueOf(ints))
	}
}

// setMap from string, creating a new map.
func (f *Flag) setMap(s string) {
	kt := f.field.Type().Key()
	vt := f.field.Type().Elem()
	mt := reflect.MapOf(kt, vt)
	a := strings.Split(s, ",")
	f.field.Set(reflect.MakeMapWithSize(mt, len(a)))
	for _, m := range a {
		pair := strings.SplitN(m, "=", 2)
		f.field.SetMapIndex(val(pair[0], kt.Kind()), val(pair[1], vt.Kind()))
	}
}

// val from string. The most useful ones for maps in CLI options
// are int and string, so that's all we're supporting for now.
func val(s string, kind reflect.Kind) reflect.Value {
	switch kind {
	case reflect.Int:
		n, _ := strconv.Atoi(s)
		return reflect.ValueOf(n)
	default:
		return reflect.ValueOf(s)
	}
}

// parseCommand with the remaining args.
func (f *Flag) parseCommand(args []string, parent string) {
	f.Args = newArgs(args)
	iface := f.field.Addr()
	p := stringer.New()
	p.WriteStrings(parent, " ", f.CommandName)
	f.Args.Parse(iface.Interface(), args, p.String())
	f.command = iface.MethodByName("Run")
	f.filter = iface.MethodByName("Filter")
}

// executeCommand specified on command line. Returns the next command, if any, or an error.
func (f *Flag) executeCommand(all bool) error {
	f.err = nil
	if !all && f.Args != nil && f.Args.execute != nil {
		return f.Args.RunCommand(all)
	}

	if f.filter.IsValid() {
		f.filter.Call([]reflect.Value{reflect.ValueOf(f.Args)})
	}

	if f.command.IsValid() && f.command.Kind() == reflect.Func {
		ret := f.command.Call([]reflect.Value{reflect.ValueOf(f.Args.Remaining)})
		err := ret[0].Interface()
		if err != nil {
			f.err = err.(error)
			return f.err
		}
	}

	if f.Args != nil && f.Args.execute != nil {
		return f.Args.RunCommand(all)
	}

	return nil
}

// AddChoice to list. Makes the choice lowercase and trims leading and trailing spaces.
func (f *Flag) AddChoice(c string) {
	f.Choices = append(f.Choices, strings.ToLower(strings.TrimSpace(c)))
}

// SetChoices to a new list.
func (f *Flag) SetChoices(list []string) {
	f.Choices = list
}
