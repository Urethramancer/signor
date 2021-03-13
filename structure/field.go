package structure

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/Urethramancer/signor/stringer"
)

// Field holds a structure's field,
type Field struct {
	// Name of the field. This can be empty if using composition.
	Name string
	// Key name of the field, if it's a map.
	Key string
	// Value name of the field is the only required string.
	Value string
	// Tags are stored for potential manipulation later
	Tags Tags
	// IsArray is mutually exclusive with Key having contents and IsMap being true.
	IsArray bool
	// IsMap means Key is set.
	IsMap bool
	// IsPointer refers to the whole field being a pointer.
	IsPointer bool
	// IsPointerValue will be true if a map or array contains pointers.
	IsPointerValue bool
	// IsComment is a comment for the next field.
	IsComment bool
}

// Tags for different input/output formats.
type Tags struct {
	// JSON tags will be output from this.
	JSON struct {
		// Name will be the field name in lowercase if unspecified.
		Name string
		// Omitempty keyword will be output if this is true.
		Omitempty bool
	}
}

func (f *Field) parseTags(tags string) {
	if f.IsComment {
		return
	}

	tags = strings.ReplaceAll(tags, "\"", "")
	if tags == "" {
		return
	}

	j := strings.SplitN(tags, ":", 2)
	for i := 0; i < len(j); {
		switch j[i] {
		case "json":
			tags := strings.Split(j[1], ",")
			for i, t := range tags {
				if i == 0 {
					f.Tags.JSON.Name = t
				} else {
					if t == "omitempty" {
						f.Tags.JSON.Omitempty = true
					}
				}
			}
			i++
		default:
			i++
		}
	}
}

// MakeTags for specified output formats. Currently only JSON is supported.
func (f *Field) MakeTags(json, omitempty bool) {
	if f.IsComment || f.Name == "" || !unicode.IsUpper(rune(f.Name[0])) {
		return
	}

	if json {
		f.Tags.JSON.Name = strings.ToLower(f.Name)
		f.Tags.JSON.Omitempty = omitempty
	}
}

// String representation of the field (compact; one tab as separator).
func (f *Field) String() string {
	var b stringer.Stringer

	if f.Name == "" {
		b.WriteString(f.Value)
		return b.String()
	}

	b.WriteStrings(f.Name, "\t")
	if f.IsPointer {
		b.WriteString("*")
	}

	if f.IsArray {
		b.WriteString("[]")
		if f.IsPointerValue {
			b.WriteString("*")
		}

		b.WriteString(f.Value)
	} else if f.IsMap {
		b.WriteStrings("map[", f.Key, "]")
		if f.IsPointerValue {
			b.WriteString("*")
		}

		b.WriteString(f.Value)
	} else {
		b.WriteString(f.Value)
	}

	if f.Tags.JSON.Name != "" {
		b.WriteStrings("\t", "`json:\"", f.Tags.JSON.Name)
		if f.Tags.JSON.Omitempty {
			b.WriteString(",omitempty")
		}

		b.WriteString("\"`")
	}

	return b.String()
}

// ProtoString is a protocol buffers representation of the field.
// The input is the order it has in the structure.
func (f *Field) ProtoString(count int) string {
	// Composite structure; we don't really have an equivalent in protobufs
	if f.Name == "" {
		return ""
	}

	if f.IsComment {
		return f.Name
	}

	if f.IsMap {
		return fmt.Sprintf("map<%s, %s> %s = %d;", f.Key, f.Value, f.Name, count)
	}

	if f.IsArray && f.Value == "byte" {
		return fmt.Sprintf("bytes %s = %d;", f.Name, count)
	}

	t := protoType(f.Value)
	if t != "" {
		return fmt.Sprintf("%s %s = %d;", protoType(f.Value), f.Name, count)
	}

	return ""
}
