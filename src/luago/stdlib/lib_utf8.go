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
	// s := ls.CheckString(1)
	// n := ls.CheckInteger(2)
	// //i = (n >= 0) ? 1 : len(s) + 1;
	// i := ls.OptInteger(3, 1)
	// ls.ArgCheck(1 <= i && i <= len(s)+1, 3,
	// 	"position out of range")
	panic("todo: utfByteOffset!")
}

// utf8.codepoint (s [, i [, j]])
// http://www.lua.org/manual/5.3/manual.html#pdf-utf8.codepoint
// lua-5.3.4/src/lutf8lib.c#codepoint()
func utfCodePoint(ls LuaState) int {
	s := ls.CheckString(1)
	i := ls.OptInteger(2, 1)
	j := ls.OptInteger(3, i)
	ls.ArgCheck(i >= 1, 2, "out of range")
	ls.ArgCheck(int(j) <= len(s), 3, "out of range")
	if i > j {
		return 0 /* empty interval; return no values */
	}

	codePoints := decodeUtf8(subStr(s, int(i), int(j)))
	ls.CheckStackL(len(codePoints), "string slice too long")
	for _, cp := range codePoints {
		ls.PushInteger(int64(cp))
	}
	return len(codePoints)
}

func decodeUtf8(str string) []rune {
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

	ls.PushString(encodeUtf8(codePoints))
	return 1
}

func encodeUtf8(codePoints []rune) string {
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
	i := int(ls.OptInteger(2, 1))
	j := int(ls.OptInteger(3, -1))
	ls.ArgCheck(1 <= i && i <= len(s)+1, 2,
		"initial position out of string")
	ls.ArgCheck(j < len(s)+1, 3,
		"final position out of string")

	s1 := subStr(s, i, j)

	// var n int64 = 0
	// for len(s1) > 0 {
	// 	r, size := utf8.DecodeRuneInString(s1)
	// 	if r == utf8.RuneError { /* conversion error? */
	// 		ls.PushNil()                 /* return nil ... */
	// 		ls.PushInteger(int64(i + 1)) /* ... and current position */
	// 		return 2
	// 	}

	// 	s1 = s1[size:]
	// 	i += size
	// 	n += 1
	// }
	// ls.PushInteger(n)

	n := utf8.RuneCountInString(s1)
	ls.PushInteger(int64(n))
	return 1
}

// utf8.codes (s)
// http://www.lua.org/manual/5.3/manual.html#pdf-utf8.codes
func utfIterCodes(ls LuaState) int {
	panic("todo: utfIterCodes!")
}
