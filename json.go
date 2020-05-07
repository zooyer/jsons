package jsons

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type (
	Raw    json.RawMessage
	Bool   bool
	Number json.Number
	String string
	Array  []interface{}
	Object map[string]interface{}
	Value  struct {
		Val interface{}
	}
)

func New(v interface{}) Value {
	var val Value
	switch v := v.(type) {
	case Value:
		return v
	case Raw, Bool, Number, String, Array, Object:
		val.Val = v
	case bool:
		val.Val = Bool(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
		val.Val = Number(fmt.Sprint(v))
	case string:
		val.Val = String(v)
	case []interface{}:
		val.Val = Array(v)
	case map[string]interface{}:
		val.Val = Object(v)
	case nil:
		val.Val = Raw(nil)
	default:
		if raw, err := json.Marshal(v); err == nil {
			if err = json.Unmarshal(raw, &val.Val); err == nil {
				return New(val.Val)
			}
			val.Val = Raw(raw)
		} else {
			val.Val = Raw(nil)
		}
	}

	return val
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
	}

	return errors.New("invalid scan object source")
}

func (v Value) Value() (driver.Value, error) {
	return json.Marshal(v.Val)
}

func (v *Value) Scan(val interface{}) error {
	if v == nil {
		v.Val = nil
		return nil
	}

	switch val := val.(type) {
	case []byte:
		return v.UnmarshalJSON(val)
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
	return json.Marshal(v.Val)
}

func (v *Value) UnmarshalJSON(data []byte) error {
	var err error
	if err = json.Unmarshal(data, &v.Val); err != nil {
		return err
	}
	v.Val = New(v.Val).Val
	return nil
}

// ************* json value function ***************

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

func (v Value) JSON() []byte {
	data, _ := v.MarshalJSON()
	if len(data) == 0 {
		return []byte("null")
	}
	return data
}

func (v Value) JSONString() string {
	return string(v.JSON())
}

func (v Value) Raw() Raw {
	if r, ok := v.Val.(Raw); ok {
		return r
	}
	if raw, err := v.MarshalJSON(); err == nil {
		return raw
	}
	return nil
}

func (v Value) Bool() Bool {
	if b, ok := v.Val.(Bool); ok {
		return b
	}
	return false
}

func (v Value) Number() Number {
	if n, ok := v.Val.(Number); ok {
		return n
	}
	return "0"
}

func (v Value) String() String {
	if s, ok := v.Val.(String); ok {
		return s
	}
	return ""
}

func (v Value) Array() Array {
	if a, ok := v.Val.(Array); ok {
		return a
	}
	return nil
}

func (v Value) Object() Object {
	if o, ok := v.Val.(Object); ok {
		return o
	}
	return nil
}

func (v Value) IsNull() bool {
	if v.IsRaw() {
		return v.Raw().IsNull()
	}

	return false
}

func (v Value) IsRaw() bool {
	_, ok := v.Val.(Raw)
	return ok
}

func (v Value) IsBool() bool {
	_, ok := v.Val.(Bool)
	return ok
}

func (v Value) IsNumber() bool {
	_, ok := v.Val.(Number)
	return ok
}

func (v Value) IsString() bool {
	_, ok := v.Val.(String)
	return ok
}

func (v Value) IsArray() bool {
	_, ok := v.Val.(Array)
	return ok
}

func (v Value) IsObject() bool {
	_, ok := v.Val.(Object)
	return ok
}

func (v Value) ToRaw() []byte {
	return v.Raw()
}

func (v Value) ToBool() bool {
	return bool(v.Bool())
}

func (v Value) ToFloat() float64 {
	f64, _ := v.Number().Float64()
	return f64
}

func (v Value) ToInt() int64 {
	i64, _ := v.Number().Int64()
	return i64
}

func (v Value) ToString() string {
	return string(v.String())
}

func (v Value) ToArray() []interface{} {
	return v.Array()
}

func (v Value) ToObject() map[string]interface{} {
	return v.Object()
}

func (v Value) Len() int {
	switch v.Val.(type) {
	case Array:
		return v.Array().Len()
	case Object:
		return v.Object().Len()
	case String:
		return len(v.String())
	}
	return 0
}

func (v Value) Cap() int {
	if v.IsArray() {
		return v.Array().Cap()
	}
	return 0
}

func (v Value) Slice(begin, end int) Array {
	return v.Array().Slice(begin, end)
}

func (v Value) Index(index int) Value {
	return v.Array().Index(index)
}

func (v Value) Range(fn func(index int, value Value) (continued bool)) bool {
	return v.Array().Range(fn)
}

func (v Value) Exist(key string) bool {
	return v.Object().Exist(key)
}

func (v Value) Get(key string) Value {
	return v.Object().Get(key)
}

func (v Value) Map(fn func(key string, value Value) (continued bool)) bool {
	return v.Object().Map(fn)
}

func (v Value) Delete(key string) {
	v.Object().Delete(key)
}

// ************ json raw function ************

func (r Raw) IsNull() bool {
	if r == nil || string(r) == "null" {
		return true
	}

	return false
}

// ************ json number function ************

func (n Number) Float64() (float64, error) {
	return json.Number(n).Float64()
}

func (n Number) Int64() (int64, error) {
	return json.Number(n).Int64()
}

func (n Number) String() string {
	return json.Number(n).String()
}

// ************ json array function *************

func (a Array) Len() int {
	return len(a)
}

func (a Array) Cap() int {
	return cap(a)
}

func (a Array) Index(index int) Value {
	if index < len(a) {
		return New(a[index])
	}
	return New(nil)
}

func (a Array) Slice(begin, end int) Array {
	return a[begin:end]
}

func (a Array) Range(fn func(index int, value Value) (continued bool)) bool {
	for index, value := range a {
		if !fn(index, New(value)) {
			return false
		}
	}
	return true
}

// ************ json object function ************

func (o Object) Len() int {
	return len(o)
}

func (o Object) Get(key string) Value {
	if len(o) > 0 {
		if val, exist := o[key]; exist {
			return New(val)
		}
	}
	return New(nil)
}

func (o Object) Delete(key string) {
	delete(o, key)
}

func (o Object) Exist(key string) bool {
	if len(o) > 0 {
		_, exist := o[key]
		return exist
	}
	return false
}

func (o Object) Map(fn func(key string, value Value) (continued bool)) bool {
	for key, value := range o {
		if !fn(key, New(value)) {
			return false
		}
	}
	return true
}
