package paths_test

import (
	"os"
	"testing"

	"github.com/Urethramancer/signor/files"
	"github.com/Urethramancer/signor/paths"
)

func TestPaths(t *testing.T) {
	t.Logf("BasePath: %s", paths.BasePath())
	if !files.DirExists(paths.BasePath()) {
		t.Errorf("%s does not exist!", paths.BasePath())
		t.FailNow()
	}

	paths.SetConfigPath("testapp")
	t.Logf("ConfigPath: %s", paths.ConfigPath())
	if !files.DirExists(paths.ConfigPath()) {
		t.Errorf("%s does not exist!", paths.ConfigPath())
		t.FailNow()
	}

	err := os.Remove(paths.ConfigPath())
	if err != nil {
		t.Errorf("Couldn't remove ConfigPath: %s", err.Error())
		t.FailNow()
	}

	paths.SetJSONConfig()
	t.Logf("JSON file: %s", paths.ConfigName())
	paths.SetINIConfig()
	t.Logf("INI file: %s", paths.ConfigName())
}
