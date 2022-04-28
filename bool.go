package jsons

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm/schema"

	"gorm.io/gorm"
)

func (b Bool) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *Bool) Scan(v interface{}) error {
	if v == nil {
		*b = false
		return nil
	}
	switch val := v.(type) {
	case []byte:
		return json.Unmarshal(val, b)
	case string:
		return json.Unmarshal([]byte(val), b)
	}

	return errors.New("invalid scan bool source")
}

func (Bool) GormDBDataType(*gorm.DB, *schema.Field) string {
	return "json"
}

func (b Bool) Bool() bool {
	return bool(b)
}

func (b Bool) JSON() []byte {
	return b.JSONValue().JSON()
}

func (b Bool) JSONValue() Value {
	return value(b)
}

func (b Bool) JSONString() string {
	return b.JSONValue().JSONString()
}
