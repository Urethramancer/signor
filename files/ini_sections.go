package files

// INISection holds one or more fields.
type INISection struct {
	Fields map[string]*INIField
	// Order fields were loaded or added in.
	Order []string
}
