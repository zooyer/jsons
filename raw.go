package jsons

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func (r Raw) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *Raw) Scan(v interface{}) error {
	if v == nil {
		*r = []byte("null")
		return nil
	}

	switch val := v.(type) {
	case []byte:
		return json.Unmarshal(val, r)
	case string:
		return json.Unmarshal([]byte(val), r)
	}

	return errors.New("invalid scan raw source")
}

func (Raw) GormDBDataType(*gorm.DB, *schema.Field) string {
	return "json"
}

func (r Raw) MarshalJSON() ([]byte, error) {
	return json.RawMessage(r).MarshalJSON()
}

func (r *Raw) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(r).UnmarshalJSON(data)
}

func (r Raw) decode(v interface{}) error {
	var buf = bytes.NewReader(r)
	decoder := json.NewDecoder(buf)
	return decoder.Decode(v)
}

func (r Raw) isNull() bool {
	if string(bytes.Trim(r, " \t\r\n")) == "null" {
		return true
	}
	return false
}

func (r Raw) isBool() (b bool, err error) {
	return b, r.decode(&b)
}

func (r Raw) isNumber() (num Number, err error) {
	return num, r.decode(&num)
}

func (r Raw) isString() (str string, err error) {
	return str, r.decode(&str)
}

func (r Raw) isArray() (arr []Raw, err error) {
	return arr, r.decode(&arr)
}

func (r Raw) isObject() (obj map[string]Raw, err error) {
	return obj, r.decode(&obj)
}

func (r Raw) IsValid() bool {
	return json.Valid(r)
}

func (r Raw) IsNull() bool {
	return r.isNull()
}

func (r Raw) IsBool() bool {
	_, err := r.isBool()
	return err == nil
}

func (r Raw) IsNumber() bool {
	_, err := r.isNumber()
	return err == nil
}

func (r Raw) IsString() bool {
	_, err := r.isString()
	return err == nil
}

func (r Raw) IsArray() bool {
	_, err := r.isArray()
	return err == nil
}

func (r Raw) IsObject() bool {
	_, err := r.isObject()
	return err == nil
}

func (r Raw) Get(keys ...interface{}) Raw {
	var val = r
	for _, key := range keys {
		switch idx := key.(type) {
		case int:
			val = val.Array()[idx]
		case string:
			val = val.Object()[idx]
		}
	}
	return val
}

func (r Raw) Int(keys ...interface{}) int64 {
	i, _ := r.Get(keys...).Number().Int64()
	return i
}

func (r Raw) Float(keys ...interface{}) float64 {
	f, _ := r.Get(keys...).Number().Float64()
	return f
}

func (r Raw) Number(keys ...interface{}) Number {
	n, _ := r.Get(keys...).isNumber()
	return n
}

func (r Raw) Bool(keys ...interface{}) bool {
	b, _ := r.Get(keys...).isBool()
	return b
}

func (r Raw) String(keys ...interface{}) string {
	s, _ := r.Get(keys...).isString()
	return s
}

func (r Raw) Array(keys ...interface{}) []Raw {
	a, _ := r.Get(keys...).isArray()
	return a
}

func (r Raw) Object(keys ...interface{}) map[string]Raw {
	o, _ := r.Get(keys...).isObject()
	return o
}

func (r Raw) Interface(keys ...interface{}) interface{} {
	var i interface{}
	_ = r.Get(keys...).decode(&i)
	return i
}

func (r Raw) Type(keys ...interface{}) string {
	var value = r.Get(keys...)
	switch {
	case value.IsObject():
		return "object"
	case value.IsArray():
		return "array"
	case value.IsString():
		return "string"
	case value.IsNumber():
		return "number"
	case value.IsBool():
		return "bool"
	case value.IsNull():
		return "null"
	}
	return "undefined"
}

func (r Raw) Len(keys ...interface{}) int {
	value := r.Get(keys...)
	switch {
	case value.IsArray():
		return len(value.Array())
	case value.IsObject():
		return len(value.Object())
	case value.IsString():
		return len(value.String())
	}
	return 0
}

func (r Raw) Keys(keys ...interface{}) []string {
	var key []string
	for k := range r.Object(keys...) {
		key = append(key, k)
	}
	return key
}

func (r Raw) Exist(keys ...interface{}) bool {
	if len(keys) < 1 {
		return false
	}
	var key = keys[len(keys)-1]
	var val = r.Get(keys[:len(keys)-1]...)
	switch idx := key.(type) {
	case int:
		if idx < val.Len() {
			return true
		}
	case string:
		if obj := val.Object(); obj != nil {
			_, exists := obj[idx]
			return exists
		}
	}
	return false
}

func (r Raw) JSON(keys ...interface{}) []byte {
	data, _ := json.Marshal(r.Get(keys...))
	return data
}

func (r Raw) JSONValue(keys ...interface{}) Value {
	var val Value
	_ = json.Unmarshal(r.JSON(keys...), &val)
	return val
}

func (r Raw) JSONString(keys ...interface{}) string {
	return string(r.JSON(keys...))
}
