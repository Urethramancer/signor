// +build linux dragonfly freebsd netbsd openbsd solaris

package paths

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/Urethramancer/signor/env"
	"github.com/Urethramancer/signor/files"
)

// Linux, BSD etc. path handling for CLI tools.

func setBasePath() {
	u, err := user.Current()
	if err != nil {
		return
	}

	dir := filepath.Join(u.HomeDir, ".config")
	basepath = env.Get("XDG_CONFIG_HOME", dir)
	if !files.DirExists(basepath) {
		_ = os.MkdirAll(basepath, 0700)
	}
}
