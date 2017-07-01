package stdlib

import "fmt"
import "strings"
import . "luago/lua"

var strLib = map[string]LuaGoFunction{
	"len":      strLen,
	"rep":      strRep,
	"reverse":  strReverse,
	"lower":    strLower,
	"upper":    strUpper,
	"sub":      strSub,
	"char":     strChar,
	"byte":     strByte,
	"format":   strFormat,
	"packsize": strPackSize,
	"pack":     strPack,
	"unpack":   strUnpack,
	"dump":     strDump,
	"find":     strFind,
	"match":    strMatch,
	"gsub":     strGsub,
	"gmatch":   strGmatch,
}

func OpenStringLib(ls LuaState) int {
	ls.NewLib(strLib)
	createMetaTable(ls)
	return 1
}

func createMetaTable(ls LuaState) {
	ls.CreateTable(0, 1)       /* table to be metatable for strings */
	ls.PushString("dummy")     /* dummy string */
	ls.PushValue(-2)           /* copy table */
	ls.SetMetaTable(-2)        /* set table as metatable for strings */
	ls.Pop(1)                  /* pop dummy string */
	ls.PushValue(-2)           /* get string library */
	ls.SetField(-2, "__index") /* metatable.__index = string */
	ls.Pop(1)                  /* pop metatable */
}

/* Basic String Functions */

// string.len (s)
// http://www.lua.org/manual/5.3/manual.html#pdf-string.len
func strLen(ls LuaState) int {
	s := ls.CheckString(1)
	ls.PushInteger(int64(len(s)))
	return 1
}

// string.rep (s, n [, sep])
// http://www.lua.org/manual/5.3/manual.html#pdf-string.rep
// lua-5.3.4/src/lstrlib.c#str_rep()
func strRep(ls LuaState) int {
	s := ls.CheckString(1)
	n := ls.CheckInteger(2)
	sep := ls.OptString(3, "")

	if n <= 0 {
		ls.PushString("")
	} else if n == 1 {
		ls.PushString(s)
	} else {
		ls.PushString(_rep(int(n), s, sep))
	}

	return 1
}

func _rep(n int, s, sep string) string {
	a := make([]string, n)
	for i := 0; i < n; i++ {
		a[i] = s
	}
	return strings.Join(a, sep)
}

// string.reverse (s)
// http://www.lua.org/manual/5.3/manual.html#pdf-string.reverse
func strReverse(ls LuaState) int {
	s := ls.CheckString(1)

	strLen := len(s)
	if strLen < 2 {
		ls.PushString(s)
	} else {
		a := make([]byte, strLen)
		for i := 0; i < strLen; i++ {
			a[i] = s[strLen-1-i]
		}
		ls.PushString(string(a))
	}

	return 1
}

// string.lower (s)
// http://www.lua.org/manual/5.3/manual.html#pdf-string.lower
func strLower(ls LuaState) int {
	s := ls.CheckString(1)
	ls.PushString(strings.ToLower(s))
	return 1
}

// string.upper (s)
// http://www.lua.org/manual/5.3/manual.html#pdf-string.upper
func strUpper(ls LuaState) int {
	s := ls.CheckString(1)
	ls.PushString(strings.ToUpper(s))
	return 1
}

// string.sub (s, i [, j])
// http://www.lua.org/manual/5.3/manual.html#pdf-string.sub
func strSub(ls LuaState) int {
	s := ls.CheckString(1)
	i := int(ls.ToInteger(2))
	j := int(ls.OptInteger(3, -1))

	sub := subStr(s, i, j)

	ls.PushString(sub)
	return 1
}

// string.char (···)
// http://www.lua.org/manual/5.3/manual.html#pdf-string.char
func strChar(ls LuaState) int {
	nArgs := ls.GetTop()

	s := make([]byte, nArgs)
	for i := 1; i <= nArgs; i++ {
		s[i-1] = byte(ls.ToInteger(i))
	}

	ls.PushString(string(s))
	return 1
}

// string.byte (s [, i [, j]])
// http://www.lua.org/manual/5.3/manual.html#pdf-string.byte
func strByte(ls LuaState) int {
	s := ls.CheckString(1)
	i := int(ls.OptInteger(2, 1))
	j := int(ls.OptInteger(3, int64(i)))

	subStr := subStr(s, i, j)
	if subStr == "" {
		ls.PushNil()
		return 1
	}

	for i := 0; i < len(subStr); i++ {
		ls.PushInteger(int64(subStr[i]))
	}
	return len(subStr)
}

// string.format (formatstring, ···)
// http://www.lua.org/manual/5.3/manual.html#pdf-string.format
func strFormat(ls LuaState) int {
	fmtStr := ls.CheckString(1)
	if len(fmtStr) <= 1 || strings.IndexByte(fmtStr, '%') < 0 {
		ls.PushString(fmtStr)
		return 1
	}

	argIdx := 1
	parsedFmt := parseFmtStr(fmtStr)
	formatted := make([]string, 0, len(parsedFmt))
	for i := 0; i < len(parsedFmt); i++ {
		tagOrStr := parsedFmt[i]
		if tagOrStr[0] == '%' {
			specifier := tagOrStr[len(tagOrStr)-1]
			if specifier != '%' {
				argIdx += 1
			}

			formatted = append(formatted, _fmtArg(specifier, tagOrStr, ls, argIdx))
		} else {
			formatted = append(formatted, tagOrStr)
		}
	}

	ls.PushString(strings.Join(formatted, ""))
	return 1
}

func _fmtArg(specifier byte, tag string, ls LuaState, argIdx int) string {
	switch specifier {
	case 'c': // character
		return string([]byte{byte(ls.ToInteger(argIdx))})
	case 'i':
		tag = tag[:len(tag)-1] + "d" // %i -> %d
		return fmt.Sprintf(tag, ls.ToInteger(argIdx))
	case 'd', 'o': // integer, octal
		return fmt.Sprintf(tag, ls.ToInteger(argIdx))
	case 'u': // unsigned integer
		tag = tag[:len(tag)-1] + "d" // %u -> %d
		return fmt.Sprintf(tag, uint(ls.ToInteger(argIdx)))
	case 'x', 'X': // hex integer
		return fmt.Sprintf(tag, uint(ls.ToInteger(argIdx)))
	case 'f': // float
		return fmt.Sprintf(tag, ls.ToNumber(argIdx))
	case 's': // string
		return fmt.Sprintf(tag, castToString(ls, argIdx))
	case 'q': // double quoted string
		return quote(ls.CheckString(argIdx))
	case '%':
		return "%"
	default:
		panic("todo! tag=" + tag)
	}
}

// string.packsize (fmt)
// http://www.lua.org/manual/5.3/manual.html#pdf-string.packsize
func strPackSize(ls LuaState) int {
	fmt := ls.CheckString(1)
	if fmt == "j" {
		ls.PushInteger(8) // todo
	} else {
		panic("strPackSize!")
	}
	return 1
}

// string.pack (fmt, v1, v2, ···)
// http://www.lua.org/manual/5.3/manual.html#pdf-string.pack
func strPack(ls LuaState) int {
	panic("strPack!")
}

// string.unpack (fmt, s [, pos])
// http://www.lua.org/manual/5.3/manual.html#pdf-string.unpack
func strUnpack(ls LuaState) int {
	panic("strUnpack!")
}

// string.dump (function [, strip])
// http://www.lua.org/manual/5.3/manual.html#pdf-string.dump
func strDump(ls LuaState) int {
	panic("strDump!")
}

/* Pattern-Matching Functions */

// string.find (s, pattern [, init [, plain]])
// http://www.lua.org/manual/5.3/manual.html#pdf-string.find
func strFind(ls LuaState) int {
	s := ls.CheckString(1)
	pattern := ls.CheckString(2)
	init := int(ls.OptInteger(3, 1))
	plain := getOptionalBoolArg(ls, 4, false)

	start, end := find(s, pattern, init, plain)

	if start < 0 {
		ls.PushNil()
		return 1
	}
	ls.PushInteger(int64(start))
	ls.PushInteger(int64(end))
	return 2
}

// string.match (s, pattern [, init])
// http://www.lua.org/manual/5.3/manual.html#pdf-string.match
func strMatch(ls LuaState) int {
	s := ls.CheckString(1)
	pattern := ls.CheckString(2)
	init := int(ls.OptInteger(3, 1))

	captures := match(s, pattern, init)

	if captures == nil {
		ls.PushNil()
		return 1
	} else {
		for i := 0; i < len(captures); i += 2 {
			capture := s[captures[i]:captures[i+1]]
			ls.PushString(capture)
		}
		return len(captures) / 2
	}
}

// string.gsub (s, pattern, repl [, n])
// http://www.lua.org/manual/5.3/manual.html#pdf-string.gsub
func strGsub(ls LuaState) int {
	s := ls.CheckString(1)
	pattern := ls.CheckString(2)
	repl := ls.CheckString(3) // todo
	n := int(ls.OptInteger(4, -1))

	newStr, nMatches := gsub(s, pattern, repl, n)
	ls.PushString(newStr)
	ls.PushInteger(int64(nMatches))
	return 2
}

// string.gmatch (s, pattern)
// http://www.lua.org/manual/5.3/manual.html#pdf-string.gmatch
func strGmatch(ls LuaState) int {
	s := ls.CheckString(1)
	pattern := ls.CheckString(2)

	gmatchAux := func(ls LuaState) int {
		captures := match(s, pattern, 1)
		if captures != nil {
			for i := 0; i < len(captures); i += 2 {
				capture := s[captures[i]:captures[i+1]]
				ls.PushString(capture)
			}
			s = s[captures[len(captures)-1]:]
			return len(captures) / 2
		} else {
			return 0
		}
	}

	ls.PushGoFunction(gmatchAux)
	return 1
}
