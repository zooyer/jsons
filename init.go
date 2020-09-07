package jsons

import (
	"database/sql/driver"
	"encoding/json"
)

var (
	_ driver.Valuer = Raw(nil)
	_ driver.Valuer = Bool(false)
	_ driver.Valuer = Number("0")
	_ driver.Valuer = String("")
	_ driver.Valuer = Array(nil)
	_ driver.Valuer = Object(nil)
	_ driver.Valuer = Value{}

	_ driver.Valuer = (*Raw)(nil)
	_ driver.Valuer = (*Bool)(nil)
	_ driver.Valuer = (*Number)(nil)
	_ driver.Valuer = (*String)(nil)
	_ driver.Valuer = (*Array)(nil)
	_ driver.Valuer = (*Object)(nil)
	_ driver.Valuer = (*Value)(nil)
)

var (
	_ json.Marshaler = Raw(nil)
	_ json.Marshaler = Number("0")
	_ json.Marshaler = Value{}

	_ json.Unmarshaler = (*Raw)(nil)
	_ json.Unmarshaler = (*Number)(nil)
	_ json.Unmarshaler = (*Value)(nil)
)
