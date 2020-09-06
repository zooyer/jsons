package jsons

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
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
				val.value = value(vv.Elem())
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

// ******** implement sql interface **********

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

// ****** implement json marshal/unmarshal *******

func (r Raw) MarshalJSON() ([]byte, error) {
	return json.RawMessage(r).MarshalJSON()
}

func (r *Raw) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(r).UnmarshalJSON(data)
}

func (n Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(json.Number(n))
}

func (n *Number) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*json.Number)(n))
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

// ************ json raw function ************

func (r Raw) IsNull() bool {
	if r == nil || string(r) == "null" {
		return true
	}
	return false
}

func (r Raw) Clone() Raw {
	if r == nil {
		return nil
	}
	return append(Raw{}, r...)
}

// ************ json number function ************

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

// ************* json value function ***************

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
	return "0"
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

// ************ json array function *************

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
	for i, v := range a {
		if v == value {
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

// ************ json object function ************

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

func (o Object) Clone(keys ...interface{}) Object {
	var object = make(Object)
	data, _ := json.Marshal(o.Get(keys...))
	_ = json.Unmarshal(data, &object)
	return object
}

func (o Object) Range(fn func(key string, value Value) (continued bool)) bool {
	for key, val := range o {
		if !fn(key, value(val)) {
			return false
		}
	}
	return true
}
