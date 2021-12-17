package jsons

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

func (s String) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *String) Scan(v interface{}) error {
	if v == nil {
		*s = ""
		return nil
	}
	switch val := v.(type) {
	case []byte:
		return json.Unmarshal(val, s)
	case string:
		return json.Unmarshal([]byte(val), s)
	}

	return errors.New("invalid scan string source")
}

func (s String) JSON() []byte {
	return s.JSONValue().JSON()
}

func (s String) JSONValue() Value {
	return value(s)
}

func (s String) JSONString() string {
	return s.JSONValue().JSONString()
}
