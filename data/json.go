package data

import (
	"encoding/json"
	"io"
)

// Serialize interface into string json
func ToJson(i interface{}, w io.Writer) error {
	encoder := json.NewEncoder(w)

	return encoder.Encode(i)
}

// Deserialize Json to the given interface
func FromJson(i interface{}, r io.Reader) error {
	decoder := json.NewDecoder(r)

	return decoder.Decode(i)
}
