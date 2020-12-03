package jsons

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tj/assert"
	"sort"
	"testing"
)

func TestRaw(t *testing.T) {
	var raw = Raw("")
	assert.Equal(t, raw.IsValid(), false)
	assert.Equal(t, Raw("null").IsValid(), true)

	assert.Equal(t, Raw(`-123`).IsNumber(), true)
	assert.Equal(t, Raw(`0`).IsNumber(), true)
	assert.Equal(t, Raw(`123`).IsNumber(), true)
	assert.Equal(t, Raw(`   123`).IsNumber(), true)

	assert.Equal(t, Raw(`""`).IsString(), true)
	assert.Equal(t, Raw(`"hello,world"`).IsString(), true)
	assert.Equal(t, Raw(`   "hello,world"`).IsString(), true)

	assert.Equal(t, Raw(`true`).IsBool(), true)
	assert.Equal(t, Raw(`false`).IsBool(), true)
	assert.Equal(t, Raw(`   true`).IsBool(), true)

	assert.Equal(t, Raw(`[]`).IsArray(), true)
	assert.Equal(t, Raw(`[1,2]`).IsArray(), true)
	assert.Equal(t, Raw(`    ["hello", "world"]`).IsArray(), true)

	assert.Equal(t, Raw(`{}`).IsObject(), true)
	assert.Equal(t, Raw(`{"name": "seven"}`).IsObject(), true)
	assert.Equal(t, Raw(`    {"age": 1024}`).IsObject(), true)

	raw = Raw(`{"key": ["hello", [1024, {"name": "seven", "age": 10240}]]}`)
	assert.Equal(t, raw.String("key", 1, 1, "name"), "seven")
	assert.Equal(t, raw.Int("key", 1, 1, "age"), int64(10240))
	assert.Equal(t, raw.Get("key", 1, 1, "name"), Raw(`"seven"`))
	assert.Equal(t, raw.Get(), raw)
	assert.Equal(t, raw.Interface("key", 1, 1, "name"), "seven")
	assert.Equal(t, raw.Len("key"), 2)
	keys := raw.Keys("key", 1, 1)
	sort.Strings(keys)
	keys2 := []string{"name", "age"}
	sort.Strings(keys2)
	assert.Equal(t, keys, keys2)

	assert.Equal(t, raw.Exist("key", 1, 0), true)
	assert.Equal(t, raw.Exist("key", 1, 1), true)
	assert.Equal(t, raw.Exist("key", 1, 2), false)
	assert.Equal(t, raw.Exist("key", 1, 1, "name"), true)

	raw = Raw(`null`)
	assert.Equal(t, raw.IsValid(), true)
	assert.Equal(t, raw.IsNull(), true)
	raw = Raw(`0`)
	assert.Equal(t, raw.IsValid(), true)
	assert.Equal(t, raw.IsNull(), false)
	raw = Raw(`-1`)
	diff := cmp.Diff(raw.IsValid(), true)
	if diff != "" {
		t.Fatal(diff)
	}
	diff = cmp.Diff(raw.IsNull(), false)
	if diff != "" {
		t.Fatal(diff)
	}
}
