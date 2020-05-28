// Package paths unifies typical configuration paths for different operating systems.
package paths

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Urethramancer/signor/files"
)

var basepath string
var configpath string
var configname string

func init() {
	setBasePath()
}

// BasePath is the basepath to create program-specific folders and dotfiles under.
func BasePath() string {
	return basepath
}

// ConfigPath is the basepath for the currentøy running app, generated via SetConfigPath().
// If called before SetConfigPath(), the path will be built from the name of the running program.
func ConfigPath() string {
	if configpath == "" {
		return SetConfigPath(os.Args[0])
	}

	return configpath
}

// SetConfigPath builds the full path for configuration files and returns it.
func SetConfigPath(program string) string {
	program = filepath.Base(program)
	program = strings.ReplaceAll(program, " ", "")
	configpath = filepath.Join(basepath, program)
	if !files.DirExists(configpath) {
		err := os.MkdirAll(configpath, 0700)
		if err != nil {
			return basepath
		}
	}

	return configpath
}

// SetJSONConfig sets the file name for the program's main configuration file, ending in .json.
func SetJSONConfig() string {
	configname = filepath.Join(configpath, "config.json")
	return configname
}

// SetINIConfig sets the file name for the program's main configuration file, ending in .ini.
func SetINIConfig() string {
	configname = filepath.Join(configpath, "config.ini")
	return configname
}

// ConfigName returns the full path to the current program's configuration file.
// Generate the name with SetJSONNConfig() or SetINIConfig().
func ConfigName() string {
	return configname
}
