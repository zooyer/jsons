package jsons

import (
	"github.com/tj/assert"
	"testing"
)

func TestMysqlPath(t *testing.T) {
	assert.Equal(t, MysqlPath(), "")
	assert.Equal(t, MysqlPath(true), `'$'`)
	assert.Equal(t, MysqlPath(false), `'$'`)
	assert.Equal(t, MysqlPath(1.1), `'$'`)
	assert.Equal(t, MysqlPath('a', 'b'), `'$'`)
	assert.Equal(t, MysqlPath("a", 'b', 1, 1.2), `'$."a".[1]'`)
	assert.Equal(t, MysqlPath(1, 2, 3, 4, 5), `'$.[1].[2].[3].[4].[5]'`)
	assert.Equal(t, MysqlPath("a", "b", "c", "d"), `'$."a"."b"."c"."d"'`)
	assert.Equal(t, MysqlPath("a", 1, "b", 2), `'$."a".[1]."b".[2]'`)
	assert.Equal(t, MysqlPath(1, "a", 2, "b"), `'$.[1]."a".[2]."b"'`)
}
