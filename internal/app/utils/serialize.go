package utils

import (
	"encoding/json"
	"io"
)

// ToJSON serializes the given interface into a string based JSON format
func ToJson(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// FromJSON deserializes the object from JSON string
// given in the io.Reader to the given interface
func FromJson(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}
