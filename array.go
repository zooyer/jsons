package jsons

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"sort"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func (a Array) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Array) Scan(v interface{}) error {
	if v == nil {
		*a = nil
		return nil
	}
	switch val := v.(type) {
	case []byte:
		return json.Unmarshal(val, a)
	case string:
		return json.Unmarshal([]byte(val), a)
	}

	return errors.New("invalid scan array source")
}

func (Array) GormDBDataType(*gorm.DB, *schema.Field) string {
	return "json"
}

func (a Array) Len(keys ...interface{}) int {
	switch len(keys) {
	case 0:
		return len(a)
	default:
		return a.Get(keys...).Len()
	}
}

func (a Array) Cap(keys ...interface{}) int {
	switch len(keys) {
	case 0:
		return cap(a)
	default:
		return a.Get(keys...).Cap()
	}
}

func (a Array) Raw(keys ...interface{}) Raw {
	return a.JSONValue(keys...).Raw()
}

func (a Array) Set(keys ...interface{}) {
	if length := len(keys); length > 1 {
		key := keys[length-2]
		val := value(keys[length-1])
		switch value := a.Get(keys[:length-2]...).value.(type) {
		case Value:
			value.Set(key, val)
		case Object:
			value.Set(key, val)
		case Array:
			if index, ok := key.(int); ok && index < len(value) {
				value[index] = val
			}
		}
	}
}

func (a Array) Get(keys ...interface{}) (val Value) {
	if a == nil {
		return
	}
	switch len(keys) {
	case 0:
		return value(a)
	case 1:
		if key, ok := keys[0].(int); ok && key < len(a) {
			return value(a[key])
		}
	default:
		var v = value(a)
		for _, k := range keys {
			switch key := k.(type) {
			case int:
				v = v.Array().Get(key)
			case string:
				v = v.Object().Get(key)
			}
		}
		return v
	}
	return
}

func (a Array) Reverse(keys ...interface{}) Array {
	switch len(keys) {
	case 0:
		for b, e := 0, len(a)-1; b < e; b, e = b+1, e-1 {
			a[b], a[e] = a[e], a[b]
		}
		return a
	default:
		return a.Get(keys...).Reverse()
	}
}

func (a Array) Clone(keys ...interface{}) Value {
	switch len(keys) {
	case 0:
		var array = make(Array, len(a), cap(a))
		data, _ := json.Marshal(a)
		_ = json.Unmarshal(data, &array)
		return value(array)
	default:
		return a.Get(keys...).Clone()
	}
}

func (a Array) Range(fn func(index int, value Value) (continued bool)) bool {
	for index, val := range a {
		if !fn(index, value(val)) {
			return false
		}
	}
	return true
}

func (a Array) Sort(less func(i, j int) bool) Array {
	sort.Slice(a, less)
	return a
}

func (a Array) Index(value interface{}) int {
	value = original(value)
	for i, v := range a {
		if original(v) == value {
			return i
		}
	}
	return -1
}

func (a Array) Contains(value interface{}) bool {
	return a.Index(value) >= 0
}

func (a Array) Slice(begin, end int) Array {
	return a[begin:end]
}

func (a Array) Append(arr Array) Array {
	return append(a, arr...)
}

func (a Array) JSON(keys ...interface{}) []byte {
	return a.Get(keys...).JSON()
}

func (a Array) JSONValue(keys ...interface{}) Value {
	return a.Get(keys...)
}

func (a Array) JSONString(keys ...interface{}) string {
	return a.Get(keys...).JSONString()
}
