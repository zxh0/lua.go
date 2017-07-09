package state

import "strings"
import . "luago/api"

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_len
func (self *luaState) Len(index int) {
	val := self.stack.get(index)
	if result, ok := self.callMetaOp1(val, "__len"); ok {
		self.stack.push(result)
		return
	}

	switch x := val.(type) {
	case string:
		length := int64(len(x))
		self.stack.push(length)
	case *luaTable:
		length := int64(x.len())
		self.stack.push(length)
	default:
		panic("todo: len!")
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawlen
func (self *luaState) RawLen(index int) uint {
	val := self.stack.get(index)
	switch x := val.(type) {
	case string:
		return uint(len(x))
	case *luaTable:
		return uint(x.len())
	default:
		return 0
	}
}

// [-n, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_concat
func (self *luaState) Concat(n int) {
	if n == 0 {
		self.stack.push("")
	} else if n == 1 {
		// do nothing
	} else if n > 1 {
		a := make([]string, n)
		for i := 0; i < n; i++ {
			a[n-1-i] = popString(self)
		}
		s := strings.Join(a, "")
		self.stack.push(s)
	} else {
		panic("todo!")
	}
}

func popString(ls *luaState) string {
	s := ""

	switch ls.Type(-1) {
	case LUA_TNIL:
		s = "nil"
	case LUA_TBOOLEAN:
		if ls.ToBoolean(-1) {
			s = "true"
		} else {
			s = "false"
		}
	case LUA_TSTRING, LUA_TNUMBER:
		s, _ = ls.ToString(-1)
	default:
		panic("todo popString()!")
	}

	ls.Pop(1)
	return s
}
