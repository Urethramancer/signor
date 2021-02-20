package files

import "fmt"

// INIField contains a variable and its data.
type INIField struct {
	// Value will be stripped of surrounding whitespace when loaded.
	Value string
	// Type lets the user choose which Get* method to use when loading unknown files.
	Type   byte
	boolV  bool
	intV   int
	floatV float64
}

// GetBool returns a field as a bool.
func (f *INIField) GetBool(key string) bool {
	return f.boolV
}

// SetBool sets a field to a bool.
func (f *INIField) SetBool(key string, value bool) {
	f.boolV = value
	f.Type = INIBool
	f.Value = fmt.Sprintf("%t", value)
}

// GetInt returns a field as an int.
func (f *INIField) GetInt(key string) int {
	return f.intV
}

// SetInt sets a field as an int.
func (f *INIField) SetInt(key string, value int) {
	f.intV = value
	f.Type = INIInt
	f.Value = fmt.Sprintf("%d", value)
}

// GetFloat returns a field as a float64.
func (f *INIField) GetFloat(key string) float64 {
	return f.floatV
}

// SetFloat sets a field to a float64.
func (f *INIField) SetFloat(key string, value float64) {
	f.floatV = value
	f.Type = INIFloat
	f.Value = fmt.Sprintf("%f", value)
}

// SetString sets a field as a string.
func (f *INIField) SetString(key, value string) {
	f.Value = value
	f.Type = INIString
}
