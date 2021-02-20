package files

import (
	"bufio"
	"os"
	"strings"

	"github.com/Urethramancer/signor/stringer"
)

// INI file base structure.
type INI struct {
	// Sections with settings.
	Sections map[string]*INISection
	// Order sections were loaded or added in.
	Order []string
}

const (
	// INIBool type
	INIBool = iota
	// INIInt type
	INIInt
	// INIFloat type
	INIFloat
	// INIString type
	INIString
)

// NewINI returns an empty INI structure.
func NewINI() *INI {
	return &INI{
		Sections: make(map[string]*INISection),
	}
}

// LoadINI from file and take a guess at the types of each value.
func LoadINI(filename string) (*INI, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	ini := NewINI()
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
				name := l[1 : len(l)-1]
				s := ini.AddSection(name)
				s.parse(r)
			}
		}
	}
	return ini, err
}

// Save outputs the INI to a file.
// If tabbed is true, the fields will be saved with a tab character prepended.
func (ini *INI) Save(filename string, tabbed bool) error {
	b := stringer.New()
	count := 0
	for _, secname := range ini.Order {
		if count > 0 {
			b.WriteString("\n")
		}
		count++
		b.WriteStrings("[", secname, "]\n")
		for _, key := range ini.Sections[secname].Order {
			f := ini.Sections[secname].Fields[key]
			if tabbed {
				b.WriteStrings("\t", key, "=", f.Value, "\n")
			} else {
				b.WriteStrings(key, "=", f.Value, "\n")
			}
		}
	}
	return WriteFile(filename, []byte(b.String()))
}

// AddSection to INI structure.
func (ini *INI) AddSection(name string) *INISection {
	sec := &INISection{
		Fields: make(map[string]*INIField),
	}
	ini.Sections[name] = sec
	ini.Order = append(ini.Order, name)
	return sec
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

// GetBool returns a field as a bool.
func (s *INISection) GetBool(key string) bool {
	v, ok := s.Fields[key]
	if !ok {
		return false
	}

	return v.boolV
}

// AddBool adds a new bool field to the section.
func (s *INISection) AddBool(key string, value bool) {
	f := INIField{}
	f.SetBool(key, value)
	s.Fields[key] = &f
	s.Order = append(s.Order, key)
}

// GetInt returns a field as an int.
func (s *INISection) GetInt(key string) int {
	v, ok := s.Fields[key]
	if !ok {
		return 0
	}

	return v.intV
}

// AddInt adds a new int field to the section.
func (s *INISection) AddInt(key string, value int) {
	f := INIField{}
	f.SetInt(key, value)
	s.Fields[key] = &f
	s.Order = append(s.Order, key)
}

// GetFloat returns a field as a float64.
func (s *INISection) GetFloat(key string) float64 {
	v, ok := s.Fields[key]
	if !ok {
		return 0.0
	}

	return v.floatV
}

// AddFloat adds a new float64 field to the section.
func (s *INISection) AddFloat(key string, value float64) {
	f := INIField{}
	f.SetFloat(key, value)
	s.Fields[key] = &f
	s.Order = append(s.Order, key)
}

// GetString returns a field as a string.
func (s *INISection) GetString(key string) string {
	v, ok := s.Fields[key]
	if !ok {
		return ""
	}

	return v.Value
}

// AddString adds a new string field to the section.
func (s *INISection) AddString(key string, value string) {
	f := INIField{}
	f.SetString(key, value)
	s.Fields[key] = &f
	s.Order = append(s.Order, key)
}
