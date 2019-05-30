package files

import (
	"encoding/json"
	"io/ioutil"
)

// LoadJSON and unmarshal structure.
func LoadJSON(fn string, out interface{}) error {
	f, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	return json.Unmarshal(f, out)
}

// SaveJSON after marshalling neatly.
func SaveJSON(path string, data interface{}) error {
	var b []byte
	var err error
	b, err = json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	return WriteFile(path, b)
}
