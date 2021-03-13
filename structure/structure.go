package structure

import (
	"strings"
	"text/scanner"

	"github.com/Urethramancer/signor/stringer"
)

// Structure holds a struct and its preceding comment.
type Structure struct {
	Name    string
	Comment string
	Fields  []*Field
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

// NewStructure simply returns a Structure struct with the specified name.
func NewStructure(name, comment string) *Structure {
	return &Structure{
		Name:    name,
		Comment: comment,
	}
}

// MakeTags for this structure. Unexported fields will be skipped.
func (st *Structure) MakeTags(json, omitempty bool) {
	for _, f := range st.Fields {
		f.MakeTags(json, omitempty)
	}
}

// String representation of the struct and contents (somewhat pretty-printed).
func (st *Structure) String() string {
	b := stringer.New()
	if st.Comment != "" {
		b.WriteStrings(st.Comment, "\n")
	}
	b.WriteStrings("type ", st.Name, " struct {\n")

	for _, f := range st.Fields {
		b.WriteString("\t")
		b.WriteStrings(f.String(), "\n")
	}
	b.WriteString("}\n")
	return b.String()
}

// ProtoString is a protobuf representation of the structure.
func (st *Structure) ProtoString() string {
	var b strings.Builder
	if st.Comment != "" {
		b.WriteString(st.Comment)
		b.WriteString("\n")
	}
	b.WriteString("message ")
	b.WriteString(st.Name)
	b.WriteString(" {\n")
	count := 1
	for _, f := range st.Fields {
		s := f.ProtoString(count)
		if strings.HasPrefix(s, "//") || strings.HasPrefix(s, "/*") {
			b.WriteString("\t")
			b.WriteString(s)
			b.WriteString("\n")
			continue
		}
		if s != "" {
			b.WriteString("\t")
			b.WriteString(s)
			b.WriteString("\n")
			count++
		}
	}
	b.WriteString("}\n")
	return b.String()
}
