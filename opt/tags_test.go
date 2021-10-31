package opt_test

import (
	"testing"

	"github.com/Urethramancer/signor/opt"
)

func TestTags(t *testing.T) {
	var o struct {
		opt.DefaultHelp
		One      bool     `short:"1" help:"An exampled flag." group:"Group A"`
		Two      bool     `short:"2" help:"An exampled flag." group:"Group A"`
		Three    bool     `short:"3" help:"An exampled flag." group:"Group B"`
		Four     bool     `short:"4" help:"An exampled flag." group:"Group B"`
		Five     bool     `short:"5" help:"An exampled flag." group:"Group C"`
		Six      bool     `short:"6" help:"An exampled flag." group:"Group C"`
		CmdOne   struct{} `command:"one" help:"Command one." group:"CG1"`
		CmdTwo   struct{} `command:"two" help:"Command two." group:"CG1"`
		CmdThree struct{} `command:"three" help:"Command three." group:"CG2"`
		CmdFour  struct{} `command:"four" help:"Command four." group:"CG2"`
	}

	a := opt.Parse(&o)
	a.Usage()
}
