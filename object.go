package jsons

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func (o Object) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func (o *Object) Scan(v interface{}) error {
	if v == nil {
		*o = nil
		return nil
	}
	switch val := v.(type) {
	case []byte:
		return json.Unmarshal(val, o)
	case string:
		return json.Unmarshal([]byte(val), o)
	}

	return errors.New("invalid scan object source")
}

func (Object) GormDBDataType(*gorm.DB, *schema.Field) string {
	return "json"
}

func (o Object) Get(keys ...interface{}) (val Value) {
	if o == nil {
		return
	}
	switch len(keys) {
	case 0:
		return value(o)
	case 1:
		if key, ok := keys[0].(string); ok {
			if _, exists := o[key]; exists {
				return value(o[key])
			}
		}
	default:
		var v = value(o)
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

func (o Object) Set(keys ...interface{}) {
	if length := len(keys); length > 1 {
		key := keys[length-2]
		val := value(keys[length-1])
		switch value := o.Get(keys[:length-2]...).value.(type) {
		case Value:
			value.Set(key, val)
		case Array:
			value.Set(key, val)
		case Object:
			if key, ok := key.(string); ok && value != nil {
				value[key] = val
			}
		}
	}
}

func (o Object) Int(keys ...interface{}) int64 {
	return o.Get(keys...).Int()
}

func (o Object) Float(keys ...interface{}) float64 {
	return o.Get(keys...).Float()
}

func (o Object) Number(keys ...interface{}) Number {
	return o.Get(keys...).Number()
}

func (o Object) Bool(keys ...interface{}) bool {
	return o.Get(keys...).Bool()
}

func (o Object) String(keys ...interface{}) string {
	return o.Get(keys...).String()
}

func (o Object) Object(keys ...interface{}) Object {
	return o.Get(keys...).Object()
}

func (o Object) Array(keys ...interface{}) Array {
	return o.Get(keys...).Array()
}

func (o Object) Interface(keys ...interface{}) interface{} {
	return o.Get(keys...).Interface()
}

func (o Object) Len(keys ...interface{}) int {
	switch len(keys) {
	case 0:
		return len(o)
	default:
		return o.Get(keys...).Len()
	}
}

func (o Object) Keys(keys ...interface{}) []string {
	var object = o.Object(keys...)
	var key = make([]string, 0, len(object))
	for k, _ := range object {
		key = append(key, k)
	}
	return key
}

func (o Object) Exist(keys ...interface{}) bool {
	if o == nil {
		return false
	}
	switch len(keys) {
	case 0:
		return o != nil
	case 1:
		if key, ok := keys[0].(string); ok {
			_, exists := o[key]
			return exists
		}
	default:
		var end = len(keys) - 1
		return o.Object(keys[:end]...).Exist(keys[end])
	}
	return false
}

func (o Object) Delete(keys ...interface{}) {
	switch len(keys) {
	case 0:
		return
	case 1:
		if key, ok := keys[0].(string); ok {
			delete(o, key)
		}
	default:
		var end = len(keys) - 1
		o.Object(keys[:end]...).Delete(keys[end])
	}
}

func (o Object) Clone(keys ...interface{}) Value {
	switch len(keys) {
	case 0:
		var object = make(Object)
		data, _ := json.Marshal(o)
		_ = json.Unmarshal(data, &object)
		return value(object)
	default:
		return o.Get(keys...).Clone()
	}
}

func (o Object) Range(fn func(key string, value Value) (continued bool)) bool {
	for key, val := range o {
		if !fn(key, value(val)) {
			return false
		}
	}
	return true
}

func (o Object) JSON(keys ...interface{}) []byte {
	return o.Get(keys...).JSON()
}

func (o Object) JSONValue(keys ...interface{}) Value {
	return o.Get(keys...)
}

func (o Object) JSONString(keys ...interface{}) string {
	return o.Get(keys...).JSONString()
}
