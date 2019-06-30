package main

var jsonHeader = `import (
	"encoding/json"
	"io/ioutil"
)
`

var jsonLoader = `
// New $STRUCT$ structure.
func New() *$STRUCT$ {
	cfg := &$STRUCT${}
	// TODO: Fill this with defaults!
	return cfg
}

// Load a $STRUCT$ structure.
func Load(filename string) (*$STRUCT$, error) {
	var out $STRUCT$
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// Save the $STRUCT$ structure.
func (c *$STRUCT$) Save(filename string) error {
	var out []byte
	out, err := json.MarshalIndent(c, "\t", "")
	if err != nil {
		return err
	}

	ioutil.WriteFile(filename, out, 0600)
	return nil
}
`
