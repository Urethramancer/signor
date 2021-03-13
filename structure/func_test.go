package structure_test

import (
	"testing"

	"github.com/Urethramancer/signor/structure"
)

func TestFunc(t *testing.T) {
	pkg, err := structure.NewPackage("func_test.go")
	if err != nil {
		t.Errorf("Couldn't load myself: %s", err.Error())
		t.FailNow()
	}

	t.Logf("\n%s\n", pkg.String())
	t.Logf("Loaded package %s from %s", pkg.Name, pkg.Filename)
	t.Logf("%s", pkg.Name)
}
