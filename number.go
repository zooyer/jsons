package jsons

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func (n Number) Value() (driver.Value, error) {
	return json.Marshal(n)
}

func (n *Number) Scan(v interface{}) error {
	if v == nil {
		*n = ""
		return nil
	}
	switch val := v.(type) {
	case []byte:
		return json.Unmarshal(val, n)
	case string:
		return json.Unmarshal([]byte(val), n)
	}

	return errors.New("invalid scan number source")
}

func (Number) GormDBDataType(*gorm.DB, *schema.Field) string {
	return "json"
}

func (n Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(json.Number(n))
}

func (n *Number) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*json.Number)(n))
}

func (n Number) Raw() Raw {
	return n.JSONValue().Raw()
}

func (n Number) Float64() (float64, error) {
	return json.Number(n).Float64()
}

func (n Number) Int64() (int64, error) {
	return json.Number(n).Int64()
}

func (n Number) Uint64() (uint64, error) {
	return strconv.ParseUint(string(n), 10, 64)
}

func (n Number) String() string {
	return json.Number(n).String()
}

func (n Number) JSON() []byte {
	return n.JSONValue().JSON()
}

func (n Number) JSONValue() Value {
	return value(n)
}

func (n Number) JSONString() string {
	return n.JSONValue().JSONString()
}
