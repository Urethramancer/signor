package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// INI file base structure.
type INI struct {
	// Sections with settings.
	Sections map[string]*INISection
	// Order sections were loaded or added in.
	Order []string
}

// INISection holds one or more fields.
type INISection struct {
	Fields map[string]*INIField
	// Order fields were loaded or added in.
	Order []string
}

// INIField is a variable and its data.
type INIField struct {
	// Value will be stripped of surrounding whitespace when loaded.
	Value string
	// Type lets the user choose which Get* method to use when loading unknown files.
	Type   byte
	boolV  bool
	intV   int
	floatV float64
}

const (
	INIBool = iota
	INIInt
	INIFloat
	INIString
)

// LoadINI from file and take a guess at the types of each value.
func LoadINI(filename string) (*INI, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	ini := INI{
		Sections: make(map[string]*INISection),
	}
	r := bufio.NewReader(f)
	loop := true
	for loop {
		l, err := r.ReadString('\n')
		if err != nil {
			loop = false
		} else {
			l = l[:len(l)-1]
			// This automatically skips comments, and really anything else
			// unknown that isn't after the first section header.
			if strings.HasPrefix(l, "[") && strings.HasSuffix(l, "]") {
				s := INISection{
					Fields: make(map[string]*INIField),
				}
				name := l[1 : len(l)-1]
				s.parse(r)
				ini.Sections[name] = &s
				ini.Order = append(ini.Order, name)
			}
		}
	}
	return &ini, err
}

// parse section properties until a new section or end of file.
func (s *INISection) parse(r *bufio.Reader) {
	loop := true
	for loop {
		next, err := r.Peek(2)
		// EOF
		if err != nil {
			return
		}

		// Skip blank lines
		if next[0] == '\n' {
			return
		}

		// New section, so this one's done
		if next[0] == '[' || next[1] == '[' {
			return
		}

		p, err := r.ReadString('\n')
		if err != nil {
			return
		}

		// Skip comments
		if strings.HasPrefix(p, "#") || strings.HasPrefix(p, ";") {
			continue
		}

		a := strings.SplitN(p, "=", 2)
		if a == nil || len(a) != 2 {
			return
		}

		a[0] = strings.TrimSpace(a[0])
		a[1] = strings.TrimSpace(a[1])
		switch a[1] {
		case "yes", "true", "no", "false":
			s.AddBool(a[0], boolValue(a[1]))
			return
		}

		// TODO: Figure out ints and floats.
		s.AddString(a[0], a[1])
	}
}

// boolValue from common strings.
func boolValue(s string) bool {
	switch s {
	case "yes", "true":
		return true
	}

	return false
}

// AddBool adds a new bool field to the section.
func (s *INISection) AddBool(key string, value bool) {
	f := INIField{}
	f.SetBool(key, value)
	s.Fields[key] = &f
	s.Order = append(s.Order, key)
}

// AddBool adds a new int field to the section.
func (s *INISection) AddInt(key string, value int) {
	f := INIField{}
	f.SetInt(key, value)
	s.Fields[key] = &f
	s.Order = append(s.Order, key)
}

// AddBool adds a new float64 field to the section.
func (s *INISection) AddFloat(key string, value float64) {
	f := INIField{}
	f.SetFloat(key, value)
	s.Fields[key] = &f
	s.Order = append(s.Order, key)
}

// AddBool adds a new string field to the section.
func (s *INISection) AddString(key string, value string) {
	f := INIField{}
	f.SetString(key, value)
	s.Fields[key] = &f
	s.Order = append(s.Order, key)
}

// SaveINI outputs the INI to a file.
// If tabbed is true, the fields will be saved with a tab character prepended.
func (f *INIField) SaveINI(filename string, tabbed bool) error {
	return nil
}

// GetBool returns the field as a bool.
func (f *INIField) GetBool(key string) bool {
	return f.boolV
}

// SetBool sets the field as a bool.
func (f *INIField) SetBool(key string, value bool) {
	f.boolV = value
	f.Type = INIBool
	f.Value = fmt.Sprintf("%t", value)
}

// GetInt returns the field as an int.
func (f *INIField) GetInt(key string) int {
	return f.intV
}

// SetInt sets the field as an int.
func (f *INIField) SetInt(key string, value int) {
	f.intV = value
	f.Type = INIInt
	f.Value = fmt.Sprintf("%d", value)
}

// GetFloat returns the field as a float64.
func (f *INIField) GetFloat(key string) float64 {
	return f.floatV
}

// SetFloat sets the field as a float64.
func (f *INIField) SetFloat(key string, value float64) {
	f.floatV = value
	f.Type = INIFloat
	f.Value = fmt.Sprintf("%f", value)
}

// SetString sets the field as a string.
func (f *INIField) SetString(key, value string) {
	f.Value = value
	f.Type = INIString
}
