package structure

import (
	"os"
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
	Funcs   []string
}

func NewPackage(filenames ...string) (*Package, error) {
	pkg := &Package{
		Structs: make([]*Structure, 0),
	}

	for _, fn := range filenames {
		var err error
		src, err := os.ReadFile(fn)
		if err != nil {
			return nil, err
		}

		pkg.contents += "\n" + string(src)
	}

	pkg.Init(strings.NewReader(pkg.contents))
	pkg.Filename = filenames[0]
	pkg.Whitespace ^= 1 << '\n'
	pkg.Mode ^= scanner.SkipComments
	pkg.Parse()
	return pkg, nil
}

// parse scans the file for structures
func (pkg *Package) Parse() {
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
			pkg.parseFunc()
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
	pkg.InternalImports = stringer.RemoveDuplicateStrings(pkg.InternalImports)
	sort.Strings(pkg.ExternalImports)
	pkg.ExternalImports = stringer.RemoveDuplicateStrings(pkg.ExternalImports)
}

func (pkg *Package) String() (string, error) {
	b := stringer.New()
	_, err := b.WriteStrings("package ", pkg.Name, "\n\n", "import (\n")
	if err != nil {
		return "", err
	}

	if len(pkg.InternalImports) > 0 {
		for _, inc := range pkg.InternalImports {
			_, err := b.WriteI("\t", "\"", inc, "\"", "\n")
			if err != nil {
				return "", err
			}
		}
	}

	if len(pkg.ExternalImports) > 0 {
		b.WriteString("\n")
		for _, inc := range pkg.ExternalImports {
			_, err := b.WriteStrings("\t", inc, "\n")
			if err != nil {
				return "", err
			}
		}
	}
	_, err = b.WriteString(")\n\n")
	if err != nil {
		return "", err
	}

	for _, st := range pkg.Structs {
		s, err := st.String()
		if err != nil {
			return "", err
		}

		_, err = b.WriteString(s)
		if err != nil {
			return "", err
		}

		_, err = b.WriteString("\n")
		if err != nil {
			return "", err
		}
	}

	for _, f := range pkg.Funcs {
		b.WriteStrings(f, "\n")
	}

	return b.String(), nil
}

// ProtoString generates protocol buffer output.
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
