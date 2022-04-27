package jsons

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
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

func (Bool) GormDataType(gorm.Dialect) string {
	return "json"
}
