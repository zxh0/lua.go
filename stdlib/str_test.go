package stdlib

import "testing"
import "github.com/stretchr/testify/assert"

// func TestSubStr(t *testing.T) {
// 	assert.Equal(t, subStr("1234567890", 0, 99), "1234567890")
// 	assert.Equal(t, subStr("1234567890", 0, 10), "1234567890")
// 	assert.Equal(t, subStr("1234567890", 1, -1), "1234567890")
// 	assert.Equal(t, subStr("1234567890", 2, -3), "2345678")
// 	assert.Equal(t, subStr("1234567890", 4, 4), "4")
// 	assert.Equal(t, subStr("1234567890", 5, 3), "")
// }

// func TestFind(t *testing.T) {
// 	assert.Equal(t, find("1234512345", "", 1, true), 1)
// 	assert.Equal(t, find("1234512345", "234", 99, true), -1)
// 	assert.Equal(t, find("1234512345", "234", 0, true), 2)
// 	assert.Equal(t, find("1234512345", "234", 1, true), 2)
// 	assert.Equal(t, find("1234512345", "234", 2, true), 2)
// 	assert.Equal(t, find("1234512345", "234", 3, true), 7)
// 	assert.Equal(t, find("1234512345", "234", -1, true), -1)
// 	assert.Equal(t, find("1234512345", "234", -4, true), 7)
// }

func TestParseFmtStr(t *testing.T) {
	assert.Equal(t, parseFmtStr(""), []string{""})
	assert.Equal(t, parseFmtStr("abc"), []string{"abc"})
	assert.Equal(t, parseFmtStr("%q%s"), []string{"%q", "%s"})
	assert.Equal(t, parseFmtStr("a%qb%sc"), []string{"a", "%q", "b", "%s", "c"})
	assert.Equal(t, parseFmtStr("%%%d"), []string{"%%", "%d"})
	assert.Equal(t, parseFmtStr("-%.20s.20s"), []string{"-", "%.20s", ".20s"})
}
