package jsons

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

func value(v interface{}) Value {
	var val Value
	switch v := v.(type) {
	case Value:
		return v
	case Bool, Number, String, Array, Object:
		val.value = v
	case bool:
		val.value = Bool(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
		var num Number
		data, _ := json.Marshal(v)
		_ = json.Unmarshal(data, &num)
		val.value = num
	case json.Number:
		val.value = Number(v)
	case string:
		val.value = String(v)
	case []interface{}:
		if v != nil {
			val.value = Array(v)
		}
	case map[string]interface{}:
		if v != nil {
			val.value = Object(v)
		}
	case nil:
		val.value = nil
	default:
		// value of copy, change will not affect the original value
		vv := reflect.ValueOf(v)
		switch vv.Kind() {
		case reflect.Slice, reflect.Array:
			if !vv.IsNil() {
				var array = make(Array, vv.Len())
				for i := 0; i < vv.Len(); i++ {
					array[i] = vv.Index(i).Interface()
				}
				val.value = array
			}
		case reflect.Map:
			if !vv.IsNil() {
				var object = make(Object)
				var keys = vv.MapKeys()
				for _, k := range keys {
					key := fmt.Sprint(k.Interface())
					val := vv.MapIndex(k).Interface()
					object[key] = val
				}
				val.value = object
			}
		case reflect.Ptr:
			if !vv.IsNil() {
				val.value = value(vv.Elem().Interface())
			}
		case reflect.Struct:
			var object = make(Object)
			var length = vv.NumField()
			for i := 0; i < length; i++ {
				key := vv.Type().Field(i).Name
				val := vv.Field(i).Interface()
				object[key] = val
			}
			val.value = object
		}
	}

	return val
}

func (v Value) Value() (driver.Value, error) {
	return json.Marshal(v.value)
}

func (v *Value) Scan(val interface{}) error {
	if val == nil {
		v.value = nil
		return nil
	}

	switch val := val.(type) {
	case []byte:
		return v.UnmarshalJSON(val)
	case string:
		return v.UnmarshalJSON([]byte(val))
	}

	return errors.New("invalid scan json source")
}

func (v Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *Value) UnmarshalJSON(data []byte) error {
	var err error
	var val interface{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err = decoder.Decode(&val); err != nil {
		return err
	}
	v.value = value(val)
	return nil
}

func (v Value) Get(keys ...interface{}) (val Value) {
	if len(keys) == 0 {
		return v
	}
	switch value := v.value.(type) {
	case Value:
		return value.Get(keys...)
	case Array:
		return v.Array().Get(keys...)
	case Object:
		return v.Object().Get(keys...)
	case []interface{}:
		return Array(value).Get(keys...)
	case map[string]interface{}:
		return Object(value).Get(keys...)
	}
	return
}

func (v Value) Set(keys ...interface{}) {
	if len(keys) == 0 {
		return
	}
	switch value := v.value.(type) {
	case Value:
		value.Set(keys...)
	case Array:
		value.Set(keys...)
	case Object:
		value.Set(keys...)
	case []interface{}:
		Array(value).Set(keys...)
	case map[string]interface{}:
		Object(value).Set(keys...)
	}
}

func (v Value) Int(keys ...interface{}) int64 {
	i, _ := v.Number(keys...).Int64()
	return i
}

func (v Value) Uint(keys ...interface{}) uint64 {
	u, _ := v.Number(keys...).Uint64()
	return u
}

func (v Value) Float(keys ...interface{}) float64 {
	f, _ := v.Number(keys...).Float64()
	return f
}

func (v Value) Number(keys ...interface{}) Number {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.Number()
	case Number:
		return value
	case json.Number:
		return Number(value)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		data, _ := json.Marshal(value)
		return Number(data)
	}
	return ""
}

func (v Value) Bool(keys ...interface{}) bool {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.Bool()
	case Bool:
		return bool(value)
	case bool:
		return value
	}
	return false
}

func (v Value) String(keys ...interface{}) string {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.String()
	case String:
		return string(value)
	case string:
		return value
	case []byte:
		return string(value)
	}
	return ""
}

func (v Value) Array(keys ...interface{}) Array {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.Array()
	case Array:
		return value
	case []interface{}:
		return value
	}
	return nil
}

func (v Value) Object(keys ...interface{}) Object {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.Object()
	case Object:
		return value
	case map[string]interface{}:
		return value
	}
	return nil
}

func (v Value) Interface(keys ...interface{}) interface{} {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.Interface()
	default:
		return value
	}
}

func (v Value) IsNull(keys ...interface{}) bool {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.IsNull()
	case Number:
		return value == ""
	case Array:
		return value == nil
	case Object:
		return value == nil
	case nil:
		return true
	case []byte:
		return len(value) == 0 || string(value) == "null"
	case []interface{}:
		return value == nil
	case map[string]interface{}:
		return value == nil
	default:
		return value == nil
	}
}

func (v Value) IsBool(keys ...interface{}) bool {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.IsBool()
	case Bool, bool:
		return true
	}
	return false
}

func (v Value) IsNumber(keys ...interface{}) bool {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.IsNumber()
	case Number, json.Number, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	}
	return false
}

func (v Value) IsString(keys ...interface{}) bool {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.IsString()
	case String, string:
		return true
	}
	return false
}

func (v Value) IsArray(keys ...interface{}) bool {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.IsArray()
	case Array, []interface{}:
		return true
	}
	return false
}

func (v Value) IsObject(keys ...interface{}) bool {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.IsObject()
	case Object, map[string]interface{}:
		return true
	}
	return false
}

func (v Value) Type(keys ...interface{}) string {
	var value = v.Get(keys...)
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

func (v Value) Len(keys ...interface{}) int {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.Len()
	case Array:
		return value.Len()
	case Object:
		return value.Len()
	case String:
		return len(value)
	case []byte:
		return len(value)
	case string:
		return len(value)
	case []interface{}:
		return len(value)
	case map[string]interface{}:
		return len(value)
	}
	return 0
}

func (v Value) Cap(keys ...interface{}) int {
	switch value := v.Get(keys...).value.(type) {
	case Value:
		return value.Cap()
	case Array:
		return value.Cap()
	case []interface{}:
		return cap(value)
	}
	return 0
}

func (v Value) Range(fn func(key interface{}, value Value) (continued bool)) bool {
	switch {
	case v.IsArray():
		return v.Array().Range(func(index int, value Value) (continued bool) {
			return fn(index, value)
		})
	case v.IsObject():
		return v.Object().Range(func(key string, value Value) (continued bool) {
			return fn(key, value)
		})
	}
	return false
}

func (v Value) Slice(begin, end int) Array {
	return v.Array().Slice(begin, end)
}

func (v Value) Index(value interface{}) int {
	return v.Array().Index(value)
}

func (v Value) Append(arr Array) Array {
	return v.Array().Append(arr)
}

func (v Value) Contains(value interface{}) bool {
	return v.Array().Contains(value)
}

func (v Value) Reverse(keys ...interface{}) Array {
	return v.Array(keys...).Reverse()
}

func (v Value) Sort(less func(i, j int) bool) Array {
	return v.Array().Sort(less)
}

func (v Value) Keys(keys ...interface{}) []string {
	return v.Object(keys...).Keys()
}

func (v Value) Exist(keys ...interface{}) bool {
	if len(keys) > 0 {
		var end = len(keys) - 1
		return v.Object(keys[:end]...).Exist(keys[end])
	}
	return false
}

func (v Value) Delete(keys ...interface{}) {
	if len(keys) > 0 {
		var end = len(keys) - 1
		v.Object(keys[:end]...).Delete(keys[end])
	}
}

func (v Value) Clone(keys ...interface{}) Value {
	val := v.Get(keys...)
	switch {
	case val.IsArray():
		return value(val.Array().Clone())
	case val.IsObject():
		return value(val.Object().Clone())
	}
	return v
}

func (v Value) JSON(keys ...interface{}) []byte {
	data, _ := v.Get(keys...).MarshalJSON()
	if len(data) == 0 {
		return []byte("null")
	}
	return data
}

func (v Value) JSONString(keys ...interface{}) string {
	return string(v.JSON(keys...))
}

func (v Value) Marshal(obj interface{}) error {
	data, err := v.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

func (v *Value) Unmarshal(obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return v.UnmarshalJSON(data)
}
