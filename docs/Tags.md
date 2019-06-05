# Supported tags in option structures
To use the flags package the user provides a structure with fields and tags to guide the input parser.

Internally, the `Args` structure will be filled with option fields of type `Flag`, with each having many possible fields which are used in parsing command line arguments and tool commands. When parsing is done, the user's structure will have its flags set and any commands will have been run.

## Option order
Any options and flags before the first tool command applies to the main option structure. Any after apply to that tool command, and if there are further tool commands they get the following options as deep as you want to go.

## Array options
When an options structure contains a string slice as a possible argument, the command line parser will fill it from a comma-separated argument. A default argument can also be supplied, with comma-separated values.

## Map options
Similar to the array options, but each comma-separated element is a key=value pair.

## Boolean flags
Specify inside the `opt` tag. These will set a corresponding boolean in the `Args` structure.

- `required`: this option must be specified. Sets `Required`.

## Short option
A `short` tag is a single symbol specified with a single hyphen (dash) in front of it. Multiple boolean flags may be combined in a dash string, and one option which takes an argument may appear among them. Behaviour when combining multiple non-boolean options will most likely not be what you want.

## Long option
A `long` tag is a keyword specified with two hyphens (double-dash) in front of it.

## Choices
The `choices` tag can contain a comma-separated list of allowed inputs. Goes well with the `default` tag.

## Placeholders
The `placeholder` tag provides a keyword to show in the usage output instead of the string for the input type. Recommended for most non-boolean options.

## Commands
Tool commands are keyword arguments which are specified as its own standalone argument. Only one top-level command can be called in one invocation, but they can be nested with their own sub-commands, each having its own structure with options and even deeper commands. A good example of this is Git.

### Aliases
Tool commands can have alternative names if the base name is too long to type all the time. Use the `aliases` tag to specify a comma-separated list of alternatives.

### Environment variables
You also have the option to use the `env` tag to specify an environment variable to read the value from. This may be useful in script-friendly programs.

## Examples

### Required option
This option is required, and has a placeholder FILE.

```go
type Options struct {
	Config	string	`opt:"required" short:"C" long:"config" help:"A configuration file." placeholder:"FILE"`
}
```

### Default option
An alternative to bailing out is to set a reasonable default.

```go
type Options struct {
	Config	string	`short:"c" help:"A configuration file." placeholder:"FILE" default:"config.json"`
}
```

Default options are also useful for choices. The default will be chosen if the option isn't specified at all, or if the supplied argument isn't one of the valid choices.

```go
type Options struct {
	Colour	string	`short:"C" long:"colour" help:"A configuration file." choices:"red,green,blue" default:"blue"`
}
```

### Option with environment variable
```go
type Options struct {
	Config	string	`short:"c" help:"A configuration file." placeholder:"FILE" default:"config.json" key:"APP_CONFIG"`
}
```

### Arrays
```go
type Options struct {
	Files	[]string	`short:"f" help:"One or more files." placeholder:"FILE..."`
	Numbers	[]int		`short:"n" help:"A bunch of numbers." placeholder:"N..." default:"1,2,3"`
}
```

```sh
$ app -f one,two,three
```

### Maps
```go
type Options struct {
	Users	map[string]string	`short:"u" long:"users" help:"Users and e-mail." placeholder:"USER=EMAIL..."`
}

```sh
$ app --users admin=admin@localhost,mailer=mailer@localhost
```

### Commands
```go
type Options struct {
	Info	InfoCmd	`command:"info" help:"Show information on a subject." aliases:"i,help"`
}

// InfoCmd options.
type InfoCmd struct {
	Subject string `opt:"required" help:"A subject to show information about."`
}

// Run the top-level command.
func (h *InfoCmd) Run(args []string) error {
	// This doesn't run if Subject is missing.
	// args will contain any additional flags the user supplied.
	return nil
}
```
