package jsons

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/tj/assert"
	"testing"
)

func TestRaw(t *testing.T) {
	var raw = Raw("")
	assert.Equal(t, raw.IsValid(), false)
	assert.Equal(t, raw.IsNull(), false)
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
	var i int
	if err := json.Unmarshal([]byte(`     123`), &i); err != nil {
		t.Fatal(err)
	}
	t.Log(i)
	assert.Equal(t, Raw(`123`).IsNumber(), true)
	assert.Equal(t, Raw(` 123`).IsNumber(), true)
}
