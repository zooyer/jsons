package jsons

import (
	"bytes"
	"encoding/json"
)

type (
	Raw    json.RawMessage
	Bool   bool
	Number json.Number
	String string
	Array  []interface{}
	Object map[string]interface{}
	Value  struct {
		value interface{}
	}
)

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalIdent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte) (val Value, err error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err = decoder.Decode(&val.value); err != nil {
		return
	}
	return
}
