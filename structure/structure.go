package structure

import (
	"strings"

	"github.com/Urethramancer/signor/stringer"
)

// Structure holds a struct and its preceding comment.
type Structure struct {
	Name    string
	Comment string
	Fields  []*Field
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
func (st *Structure) String() (string, error) {
	b := stringer.New()
	if st.Comment != "" {
		_, err := b.WriteStrings(st.Comment, "\n")
		if err != nil {
			return "", err
		}
	}
	_, err := b.WriteStrings("type ", st.Name, " struct {\n")
	if err != nil {
		return "", err
	}

	for _, f := range st.Fields {
		b.WriteString("\t")
		x, err := f.String()
		if err != nil {
			return "", err
		}

		_, err = b.WriteStrings(x, "\n")
		if err != nil {
			return "", err
		}
	}
	b.WriteString("}\n")
	return b.String(), nil
}

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
