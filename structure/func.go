package structure

import (
	"strings"
	"text/scanner"

	"github.com/Urethramancer/signor/stringer"
)

// parseFunc parses and stores functions.
func (pkg *Package) parseFunc() {
	b := stringer.New()
	b.WriteString("func ")
	braces := 1
	nl := false
	for pkg.tok = pkg.Scan(); pkg.tok != scanner.EOF && braces > 0; pkg.tok = pkg.Scan() {
		if nl {
			if pkg.TokenText() == "}" {
				b.WriteStrings(strings.Repeat("\t", braces-2))
			} else {
				b.WriteString(strings.Repeat("\t", braces-1))
			}
			nl = false
		}

		switch pkg.TokenText() {
		case "\n":
			b.WriteStrings("\n")
			nl = true

		case "{":
			braces++
			b.WriteString(" {")

		case "}":
			braces--
			b.WriteString("}")

		case ":", "!", "*":
			b.WriteStrings(" ", pkg.TokenText())

		case "if", "for", ",", "=":
			b.WriteStrings(pkg.TokenText(), " ")

		default:
			b.WriteString(pkg.TokenText())
		}

	}

	pkg.Funcs = append(pkg.Funcs, b.String())
}
