// +build darwin

package paths

import (
	"os/user"
	"path/filepath"
)

// macOS path handling for CLI tools.

func setBasePath() {
	u, err := user.Current()
	if err != nil {
		return
	}

	basepath = filepath.Join(u.HomeDir, "Library", "Application Support")
}
