package jsons

import (
	"testing"
)

func Test(t *testing.T) {
	assert := func(expression bool, msg ...interface{}) {
		if !expression {
			t.Fatal(msg...)
		}
	}

	var val = value(nil)
	data, err := val.MarshalJSON()
	assert(err == nil, err)

	switch val.value.(type) {
	case nil:
		assert(string(data) == "null", string(data))
	case Bool, Number, String, Array, Object, Value:
	default:
		t.Fatal("error type")
	}

	val = value(Object{
		"b": Object{
			"c": Object{
				"d": value("d"),
			},
		},
		"c": Object{
			"d": "d",
		},
	})

	val.Set("a", "b")
	assert(val.String("a") == "b")
	val.Delete("a")
	assert(!val.Exist("a"))

	val.Object("c").Set("d", "dd")
	assert(val.String("c", "d") == "dd")
	val.Object("c")["d"] = "ddd"
	assert(val.String("c", "d") == "ddd")
	val.Delete("c", "d")
	assert(!val.Exist("c", "d"))
	val.Delete("c")
	assert(!val.Exist("c"))

	assert(val.String("b", "c", "d") == "d")
	val.Get("b").Get("c").Set("d", "dd")
	assert(val.String("b", "c", "d") == "dd")
	val.Get("b", "c").Set("d", "ddd")
	assert(val.String("b", "c", "d") == "ddd")
	val.Object("b", "c")["d"] = "4d"
	assert(val.String("b", "c", "d") == "4d")

	assert(value(false).IsBool())
	assert(value(0).IsNumber())
	assert(value("").IsString())
	assert(value([]int{}).IsArray())
	assert(value(map[string]int{}).IsObject())
	assert(value(struct{}{}).IsObject())

	assert(value(nil).JSONString() == "null")
	assert(value(false).JSONString() == "false")
	assert(value(0).JSONString() == "0")
	assert(value("").JSONString() == `""`)
	assert(value([]int{}).JSONString() == "[]")
	assert(value(map[string]int{}).JSONString() == "{}")
	assert(value(struct{}{}).JSONString() == "{}")
	assert(value((*int)(nil)).JSONString() == "null")
	assert(value(([]int)(nil)).JSONString() == "null")
	assert(value((map[string]int)(nil)).JSONString() == "null")
	assert(value((*struct{})(nil)).JSONString() == "null")
	assert(value(map[int]int{1: 2}).JSONString() == `{"1":2}`)

	val = value(nil)
	assert(val.Object().Get("not exist key").Bool() == false)
	assert(val.Object().Get("not exist key").Number() == "")
	assert(val.Object().Get("not exist key").String() == "")
	assert(val.Object().Get("not exist key").Array() == nil)
	assert(val.Object().Get("not exist key").Object() == nil)

	assert(val.Array().Get(99999).Bool() == false)
	assert(val.Array().Get(99999).Number() == "")
	assert(val.Array().Get(99999).String() == "")
	assert(val.Array().Get(99999).Array() == nil)
	assert(val.Array().Get(99999).Object() == nil)

	assert(val.Bool() == false)
	assert(val.Number() == "")
	assert(val.String() == "")
	assert(val.Array() == nil)
	assert(val.Object() == nil)

	val = value(Object{
		"a": "aa",
		"b": "bb",
		"c": Object{
			"d": Object{
				"e": "ee",
			},
			"e": Array{
				1, "hello", "world",
			},
		},
	})
	assert(val.Int("c", "e", 0) == 1)
	val.Set("c", "e", 0, value(2))
	assert(val.Int("c", "e", 0) == 2)
	assert(val.Exist("c", "d", "e"), "e must exist")
	assert(!val.Exist("c", "d", "e", "f"), "f must not exist")
	val.Delete("c", "d", "e")
	assert(!val.Exist("c", "d", "e"), "e must not exist")
	val.Get("c", "d").Set("e", "eee")
	assert(val.Exist("c", "d", "e"), "e must exist")
}
