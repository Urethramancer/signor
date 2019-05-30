package files

import (
	"os"
)

// Exists checks for the existence of a file or directory.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks for the existence of a directory.
// It will return false if a file with that name exists.
func DirExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	if stat.IsDir() {
		return true
	}

	return false
}

// FileExists checks for the existence of a file.
// It will return false if a directory with that name exists.
func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	if stat.IsDir() {
		return false
	}

	return true
}

// WriteFile saves a file with private permissions.
func WriteFile(path string, data []byte) error {
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	defer NewCloser(out).Close(true)
	_, err = out.Write(data)
	return err
}

// EnsureDirExists by creating it.
func EnsureDirExists(path string) error {
	if !DirExists(path) {
		return os.MkdirAll(path, 0700)
	}

	return nil
}
