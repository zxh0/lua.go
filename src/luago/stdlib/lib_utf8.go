package stdlib

import "unicode/utf8"
import . "luago/api"

/* pattern to match a single UTF-8 character */
const UTF8PATT = "[\x00-\x7F\xC2-\xF4][\x80-\xBF]*"

const MAX_UNICODE = 0x10FFFF

var utf8Lib = map[string]GoFunction{
	"offset":    utfByteOffset,
	"codepoint": utfCodePoint,
	"char":      utfChar,
	"len":       utfLen,
	"codes":     utfIterCodes,
	/* placeholders */
	"charpattern": nil,
}

func OpenUTF8Lib(ls LuaState) int {
	ls.NewLib(utf8Lib)
	ls.PushString(UTF8PATT)
	ls.SetField(-2, "charpattern")
	return 1
}

// utf8.offset (s, n [, i])
// http://www.lua.org/manual/5.3/manual.html#pdf-utf8.offset
func utfByteOffset(ls LuaState) int {
	s := ls.CheckString(1)
	sLen := int64(len(s))
	n := ls.CheckInteger(2)
	i := int64(1)
	if n < 0 {
		i = sLen + 1
	}
	i = posRelat(ls.OptInteger(3, i))
	ls.ArgCheck(1 <= i && i <= sLen+1, 3,
		"position out of range")

	// todo
	panic("todo: utfByteOffset!")
}

// utf8.codepoint (s [, i [, j]])
// http://www.lua.org/manual/5.3/manual.html#pdf-utf8.codepoint
// lua-5.3.4/src/lutf8lib.c#codepoint()
func utfCodePoint(ls LuaState) int {
	s := ls.CheckString(1)
	sLen := len(s)
	i := posRelat(ls.OptInteger(2, 1), sLen)
	j := posRelat(ls.OptInteger(3, int64(i)), sLen)

	ls.ArgCheck(i >= 1, 2, "out of range")
	ls.ArgCheck(int(j) <= sLen, 3, "out of range")
	if i > j {
		return 0 /* empty interval; return no values */
	}
	//if (pose - posi >= INT_MAX)  /* (lua_Integer -> int) overflow? */
	//	return luaL_error(L, "string slice too long");

	codePoints := _decodeUtf8(s[i-1 : j])
	ls.CheckStack2(len(codePoints), "string slice too long")
	for _, cp := range codePoints {
		ls.PushInteger(int64(cp))
	}
	return len(codePoints)
}

func _decodeUtf8(str string) []rune {
	codePoints := make([]rune, 0, len(str))

	for len(str) > 0 {
		r, size := utf8.DecodeRuneInString(str)
		codePoints = append(codePoints, r)

		str = str[size:]
	}

	return codePoints
}

// utf8.char (···)
// http://www.lua.org/manual/5.3/manual.html#pdf-utf8.char
// lua-5.3.4/src/lutf8lib.c#utfchar()
func utfChar(ls LuaState) int {
	n := ls.GetTop() /* number of arguments */
	codePoints := make([]rune, n)

	for i := 1; i <= n; i++ {
		cp := ls.CheckInteger(i)
		ls.ArgCheck(0 <= cp && cp <= MAX_UNICODE, i, "value out of range")
		codePoints[i-1] = rune(cp)
	}

	ls.PushString(_encodeUtf8(codePoints))
	return 1
}

func _encodeUtf8(codePoints []rune) string {
	buf := make([]byte, 6)
	str := make([]byte, 0, len(codePoints))

	for _, cp := range codePoints {
		n := utf8.EncodeRune(buf, cp)
		str = append(str, buf[0:n]...)
	}

	return string(str)
}

// utf8.len (s [, i [, j]])
// http://www.lua.org/manual/5.3/manual.html#pdf-utf8.len
// lua-5.3.4/src/lutf8lib.c#utflen()
func utfLen(ls LuaState) int {
	s := ls.CheckString(1)
	sLen := len(s)
	i := posRelat(ls.OptInteger(2, 1), sLen)
	j := posRelat(ls.OptInteger(3, -1), sLen)
	ls.ArgCheck(1 <= i && i <= sLen+1, 2,
		"initial position out of string")
	ls.ArgCheck(j <= sLen, 3,
		"final position out of string")

	if i > j {
		ls.PushInteger(0)
	} else {
		n := utf8.RuneCountInString(s[i-1 : j])
		ls.PushInteger(int64(n))
	}

	return 1
}

// utf8.codes (s)
// http://www.lua.org/manual/5.3/manual.html#pdf-utf8.codes
func utfIterCodes(ls LuaState) int {
	panic("todo: utfIterCodes!")
}
