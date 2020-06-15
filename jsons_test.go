package jsons

import (
	"testing"
)

func TestNew(t *testing.T) {
	assert := func(expression bool, msg ...interface{}) {
		if !expression {
			t.Fatal(msg...)
		}
	}

	var val = New(nil)
	switch val.Val.(type) {
	case Raw, Bool, Number, String, Array, Object, Value:
		data, err := val.MarshalJSON()
		assert(err == nil, err)
		assert(string(data) == "null", string(data))
	default:
		t.Fatal("error type")
	}

	assert(New(nil).IsRaw())
	assert(New(false).IsBool())
	assert(New(0).IsNumber())
	assert(New("").IsString())
	assert(New([]int{}).IsArray())
	assert(New(map[string]int{}).IsObject())
	assert(New(struct{}{}).IsObject())

	assert(New(nil).JSONString() == "null")
	assert(New(false).JSONString() == "false")
	assert(New(0).JSONString() == "0")
	assert(New("").JSONString() == `""`)
	assert(New([]int{}).JSONString() == "[]")
	assert(New(map[string]int{}).JSONString() == "{}")
	assert(New(struct{}{}).JSONString() == "{}")
	assert(New((*int)(nil)).JSONString() == "null")
	assert(New(([]int)(nil)).JSONString() == "null")
	assert(New((map[string]int)(nil)).JSONString() == "null")
	assert(New((*struct{})(nil)).JSONString() == "null")
	assert(New(map[int]int{1: 2}).JSONString() == `{"1":2}`)

	val = New(nil)
	assert(val.Object().Get("not exist key").Raw() == nil)
	assert(val.Object().Get("not exist key").Bool() == false)
	assert(val.Object().Get("not exist key").Number() == "0")
	assert(val.Object().Get("not exist key").String() == "")
	assert(val.Object().Get("not exist key").Array() == nil)
	assert(val.Object().Get("not exist key").Object() == nil)

	assert(val.Object().Get("not exist key").Array().Index(99999).Raw() == nil)
	assert(val.Object().Get("not exist key").Object().Get("not exist key").Raw() == nil)

	assert(val.Array().Index(99999).Raw() == nil)
	assert(val.Array().Index(99999).Bool() == false)
	assert(val.Array().Index(99999).Number() == "0")
	assert(val.Array().Index(99999).String() == "")
	assert(val.Array().Index(99999).Array() == nil)
	assert(val.Array().Index(99999).Object() == nil)

	assert(val.Array().Index(99999).Array().Index(99999).Raw() == nil)
	assert(val.Array().Index(99999).Object().Get("not exist key").Raw() == nil)

	assert(val.Raw() == nil)
	assert(val.Bool() == false)
	assert(val.Number() == "0")
	assert(val.String() == "")
	assert(val.Array() == nil)
	assert(val.Object() == nil)

	val = New(Object{
		"a": "aa",
		"b": "bb",
		"c": Object{
			"d": Object{
				"e": "ee",
			},
		},
	})
	assert(val.Exist("c", "d", "e"), "e must exist")
	assert(!val.Exist("c", "d", "e", "f"), "f must not exist")
	val.Delete("c", "d", "e")
	assert(!val.Exist("c", "d", "e"), "e must not exist")
	val.Get("c", "d").Set("e", "eee")
	assert(val.Exist("c", "d", "e"), "e must exist")
	assert(val.Get("c", "d", "e").ToString() == "eee", "must eee")
}
