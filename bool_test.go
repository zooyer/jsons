package jsons

import (
	"testing"

	"database/sql/driver"
	"github.com/tj/assert"
)

func TestBool_JSON(t *testing.T) {
	var b Bool
	assert.Equal(t, b.JSON(), []byte("false"))
	b = true
	assert.Equal(t, b.JSON(), []byte("true"))
	b = false
	assert.Equal(t, b.JSON(), []byte("false"))
}

func TestBool_JSONString(t *testing.T) {
	var b Bool
	assert.Equal(t, b.JSONString(), "false")
	b = true
	assert.Equal(t, b.JSONString(), "true")
	b = false
	assert.Equal(t, b.JSONString(), "false")
}

func TestBool_JSONValue(t *testing.T) {
	var b Bool
	assert.Equal(t, b.JSONValue(), value(false))
	b = true
	assert.Equal(t, b.JSONValue(), value(true))
	b = false
	assert.Equal(t, b.JSONValue(), value(false))
}

func TestBool_Scan(t *testing.T) {
	var b Bool
	var err error

	if err = b.Scan("true"); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, b.Bool(), true)

	if err = b.Scan("false"); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, b.Bool(), false)

	if err = b.Scan([]byte("true")); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, b.Bool(), true)

	if err = b.Scan([]byte("false")); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, b.Bool(), false)
}

func TestBool_Value(t *testing.T) {
	var b Bool
	var err error
	var value driver.Value

	b = true
	if value, err = b.Value(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, value, []byte("true"))

	b = false
	if value, err = b.Value(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, value, []byte("false"))
}
