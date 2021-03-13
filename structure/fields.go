package structure

import (
	"strings"
	"text/scanner"
)

// parseFields parses all the fields of a structure.
func (pkg *Package) parseFields(st *Structure) {
	pkg.tok = pkg.Scan()
	if pkg.tok == scanner.EOF {
		return
	}

	for pkg.tok = pkg.Scan(); pkg.tok != scanner.EOF && pkg.TokenText() != "}"; pkg.tok = pkg.Scan() {
		if pkg.TokenText() != "\n" {
			f := pkg.parseField()
			if f != nil {
				st.Fields = append(st.Fields, f)
			}
		}
	}
}

// parseField parses an individual field, skipping blank lines.
func (pkg *Package) parseField() *Field {
	f := &Field{}
	name := pkg.TokenText()
	if pkg.tok == scanner.Comment {
		f.Name = name
		f.IsComment = true
		return f
	}
	// Probably going to be an embedded struct here
	if pkg.Peek() == '.' {
		for pkg.tok = pkg.Scan(); pkg.tok != scanner.EOF && pkg.TokenText() == "." && pkg.TokenText() != "\n"; pkg.tok = pkg.Scan() {
			pkg.tok = pkg.Scan()
			if pkg.tok == scanner.EOF {
				// Malformed struct
				return nil
			}
			name += "." + pkg.TokenText()
		}
	}

	// Definitely an embed, so we're done
	if pkg.TokenText() == "\n" {
		pkg.tok = pkg.Scan()
		f.Value = name
		return f
	}

	// Not an embed, so what we have so far is the field name
	f.Name = name

	for pkg.tok = pkg.Scan(); pkg.tok != scanner.EOF; pkg.tok = pkg.Scan() {
		switch pkg.TokenText() {
		case "\n":
			n := strings.Index(f.Value, "`")
			if n > -1 {
				t := f.Value[n:len(f.Value)]
				t = strings.ReplaceAll(t, "`", "")
				f.parseTags(t)
				f.Value = f.Value[0:n]
			}
			return f

		case "*":
			if f.IsArray || f.IsMap {
				f.IsPointerValue = true
			} else {
				f.IsPointer = true
			}

		case "map":
			f.IsMap = true

		case "[":
			if f.IsMap {
				pkg.tok = pkg.Scan()
				if pkg.tok == scanner.EOF {
					// Well, that's not right
					return nil
				}

				f.Key = pkg.TokenText()

			} else {
				f.IsArray = true
			}

		case "]":
			// Just eat the closing bracket

		case ".":
			pkg.tok = pkg.Scan()
			if pkg.tok == scanner.EOF {
				return nil
			}
			f.Value += "." + pkg.TokenText()

		default:
			f.Value += pkg.TokenText()
		}
	}

	// We won't get here very often, if at all
	return nil
}
