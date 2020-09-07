package jsons

import (
	"fmt"
	"strconv"
)

func MysqlPath(keys ...interface{}) string {
	var path string
	if len(keys) > 0 {
		path += "'$"
		for _, k := range keys {
			switch k := k.(type) {
			case int:
				path += fmt.Sprintf(".[%d]", k)
			case string:
				path += fmt.Sprintf(".%s", strconv.Quote(k))
			}
		}
		path += "'"
	}
	return path
}
