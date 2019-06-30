package structure

import (
	"io/ioutil"
	"sort"
	"strings"
	"text/scanner"

	"github.com/Urethramancer/signor/stringer"
)

/*
 *	Package holds the package scanner and other internal data,
 *	plus the exported structures for structs and imports.
 */
type Package struct {
	scanner.Scanner
	tok      rune
	Name     string
	contents string

	InternalImports []string
	ExternalImports []string

	Structs []*Structure
}

func NewPackage(filenames ...string) (*Package, error) {
	pkg := &Package{
		Structs: make([]*Structure, 0),
	}

	for _, fn := range filenames {
		var err error
		src, err := ioutil.ReadFile(fn)
		if err != nil {
			return nil, err
		}

		pkg.contents += "\n" + string(src)
	}

	pkg.Init(strings.NewReader(pkg.contents))
	pkg.Filename = filenames[0]
	pkg.Whitespace ^= 1 << '\n'
	pkg.Mode ^= scanner.SkipComments
	pkg.parse()
	return pkg, nil
}

// parse scans the file for structures
func (pkg *Package) parse() {
	var comment string
	for pkg.tok = pkg.Scan(); pkg.tok != scanner.EOF; pkg.tok = pkg.Scan() {
		if pkg.tok == scanner.Comment {
			comment = pkg.TokenText()
			continue
		}

		switch pkg.TokenText() {
		case "package":
			pkg.parsePackage()
			comment = ""

		case "import":
			pkg.parseImports()
			comment = ""

		case "type":
			pkg.parseType(comment)
			comment = ""

		case "func":
			pkg.skipFunc()
			comment = ""

		case "\n":

		default:
		}
	}
	pkg.sortImports()
}

func (pkg *Package) parsePackage() {
	pkg.tok = pkg.Scan()
	if pkg.tok == scanner.EOF {
		return
	}

	if pkg.Name == "" {
		pkg.Name = pkg.TokenText()
	}
}

func (pkg *Package) parseImports() {
	pkg.tok = pkg.Scan()
	if pkg.tok == scanner.EOF {
		return
	}

	if pkg.TokenText() == "(" {
		pkg.tok = pkg.Scan()
		for pkg.TokenText() != ")" && pkg.tok != scanner.EOF {
			if pkg.TokenText() != "\n" {
				pkg.addImport(pkg.TokenText())
			}
			pkg.tok = pkg.Scan()
		}
	} else {
		pkg.addImport(pkg.TokenText())
		pkg.tok = pkg.Scan()
	}
}

func (pkg *Package) addImport(i string) {
	if strings.Contains(i, ".") && strings.Contains(i, "/") {
		pkg.ExternalImports = append(pkg.ExternalImports, pkg.TokenText())
	} else {
		pkg.InternalImports = append(pkg.InternalImports, pkg.TokenText())
	}
}

func (pkg *Package) sortImports() {
	sort.Strings(pkg.InternalImports)
	pkg.InternalImports = removeDuplicateStrings(pkg.InternalImports)
	sort.Strings(pkg.ExternalImports)
	pkg.ExternalImports = removeDuplicateStrings(pkg.ExternalImports)
}

// parseType mainly looks for structures and their fields.
func (pkg *Package) parseType(comment string) {
	pkg.tok = pkg.Scan()
	if pkg.tok == scanner.EOF {
		return
	}

	name := pkg.TokenText()
	if name == "(" {
		pkg.parseTypes()
		return
	}

	pkg.tok = pkg.Scan()
	// If we drop out at this point, the type is incomplete.
	if pkg.tok == scanner.EOF {
		return
	}

	if pkg.TokenText() == "interface" {
		pkg.skipInterface()
		return
	}

	if pkg.TokenText() == "struct" {
		st := NewStructure(name, comment)
		pkg.Structs = append(pkg.Structs, st)
		pkg.parseFields(st)
	}
}

// parseTypes is used to parse multiple types in parentheses.
// TODO: Actually parse these rather than skipping.
func (pkg *Package) parseTypes() {
	parens := 1
	for pkg.tok = pkg.Scan(); pkg.tok != scanner.EOF && parens > 0; pkg.tok = pkg.Scan() {
		switch pkg.TokenText() {
		case "(":
			parens++

		case ")":
			parens--
		}
	}
}

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

// skipFunc parses past the end of a function since we don't care about those.
func (pkg *Package) skipFunc() {
	for pkg.tok = pkg.Scan(); pkg.tok != scanner.EOF && pkg.TokenText() != "{"; pkg.tok = pkg.Scan() {
	}
	braces := 1
	for pkg.tok = pkg.Scan(); pkg.tok != scanner.EOF && braces > 0; pkg.tok = pkg.Scan() {
		switch pkg.TokenText() {
		case "{":
			braces++

		case "}":
			braces--
		}
	}
}

// skipInterface parses past the end of an interface.
func (pkg *Package) skipInterface() {
	for pkg.tok = pkg.Scan(); pkg.tok != scanner.EOF && pkg.TokenText() != "{"; pkg.tok = pkg.Scan() {
	}
	braces := 1
	for pkg.tok = pkg.Scan(); pkg.tok != scanner.EOF && braces > 0; pkg.tok = pkg.Scan() {
		switch pkg.TokenText() {
		case "{":
			braces++

		case "}":
			braces--
		}
	}
}

// MakeTags for all structures. Unexported fields will be skipped.
func (pkg *Package) MakeTags(json, omitempty bool) {
	for _, st := range pkg.Structs {
		st.MakeTags(json, omitempty)
	}
}

func (pkg *Package) String() string {
	b := stringer.New()
	b.WriteStrings("package ", pkg.Name, "\n\n")

	b.WriteString("import (\n")
	if len(pkg.InternalImports) > 0 {
		for _, inc := range pkg.InternalImports {
			b.WriteI("\t", "\"", inc, "\"", "\n")
		}
	}

	if len(pkg.ExternalImports) > 0 {
		b.WriteString("\n")
		for _, inc := range pkg.ExternalImports {
			b.WriteStrings("\t", inc, "\n")
		}
	}
	b.WriteString(")\n\n")

	for _, st := range pkg.Structs {
		b.WriteString(st.String())
		b.WriteString("\n")
	}
	return b.String()
}

func (pkg *Package) ProtoString() string {
	var b strings.Builder
	b.WriteString("syntax = ")
	b.WriteRune('"')
	b.WriteString("proto3")
	b.WriteRune('"')
	b.WriteString(";\npackage ")
	b.WriteString(pkg.Name)
	b.WriteString(";\n\n")
	for _, st := range pkg.Structs {
		b.WriteString(st.ProtoString())
		b.WriteString("\n")
	}
	return b.String()
}

func (pkg *Package) MergeExternalImports() []string {
	list := []string{}

	m := make(map[string]bool)
	for _, i := range pkg.ExternalImports {
		a := strings.Split(i, "/")
		i = strings.Join(a[:3], "/")
		m[i] = true
	}
	for k := range m {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}
