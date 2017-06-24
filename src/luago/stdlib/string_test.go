package stdlib

import "testing"
import "assert"

func TestSubStr(t *testing.T) {
	assert.StringEqual(t, subStr("1234567890", 0, 99), "1234567890")
	assert.StringEqual(t, subStr("1234567890", 0, 10), "1234567890")
	assert.StringEqual(t, subStr("1234567890", 1, -1), "1234567890")
	assert.StringEqual(t, subStr("1234567890", 2, -3), "2345678")
	assert.StringEqual(t, subStr("1234567890", 4, 4), "4")
	assert.StringEqual(t, subStr("1234567890", 5, 3), "")
}

func TestFind(t *testing.T) {
	assert.IntEqual(t, find("1234512345", "", 1, true), 1)
	assert.IntEqual(t, find("1234512345", "234", 99, true), -1)
	assert.IntEqual(t, find("1234512345", "234", 0, true), 2)
	assert.IntEqual(t, find("1234512345", "234", 1, true), 2)
	assert.IntEqual(t, find("1234512345", "234", 2, true), 2)
	assert.IntEqual(t, find("1234512345", "234", 3, true), 7)
	assert.IntEqual(t, find("1234512345", "234", -1, true), -1)
	assert.IntEqual(t, find("1234512345", "234", -4, true), 7)
}

func TestParseFmtStr(t *testing.T) {
	assert.StringsEqual(t, parseFmtStr(""), []string{""})
	assert.StringsEqual(t, parseFmtStr("abc"), []string{"abc"})
	assert.StringsEqual(t, parseFmtStr("%q%s"), []string{"%q", "%s"})
	assert.StringsEqual(t, parseFmtStr("a%qb%sc"), []string{"a", "%q", "b", "%s", "c"})
	assert.StringsEqual(t, parseFmtStr("%%%d"), []string{"%%", "%d"})
	assert.StringsEqual(t, parseFmtStr("-%.20s.20s"), []string{"-", "%.20s", ".20s"})
}
