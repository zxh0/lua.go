package state

import "runtime"
import "strings"
import "luago/luanum"
import . "luago/api"

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_close
// lua-5.3.4/src/lstate.c#lua_close
func (self *luaState) Close() {
	// todo
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_atpanic
// lua-5.3.4/src/lapi.c#lua_atpanic
func (self *luaState) AtPanic(panicf GoFunction) GoFunction {
	oldPanicf := self.panicf
	self.panicf = panicf
	return oldPanicf
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_version
// lua-5.3.4/src/lapi.c#lua_version
func (self *luaState) Version() float64 {
	return LUA_VERSION_NUM
}

// [-1, +0, v]
// http://www.lua.org/manual/5.3/manual.html#lua_error
func (self *luaState) Error() int {
	panic("todo!")
}

// [-0, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_gc
func (self *luaState) GC(what, data int) int {
	runtime.GC()
	return 0
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_stringtonumber
func (self *luaState) StringToNumber(s string) bool {
	if n, ok := luanum.ParseInteger(s, 10); ok {
		self.PushInteger(n)
		return true
	}
	if n, ok := luanum.ParseFloat(s); ok {
		self.PushNumber(n)
		return true
	}
	return false
}

// [-1, +(2|0), e]
// http://www.lua.org/manual/5.3/manual.html#lua_next
func (self *luaState) Next(idx int) bool {
	t := self.stack.get(idx)
	if tbl, ok := t.(*luaTable); ok {
		key := self.stack.pop()
		nextKey, nextVal := tbl.next(key)
		if nextKey != nil {
			self.stack.push(nextKey)
			self.stack.push(nextVal)
			return true
		} else {
			return false
		}
	}
	panic("not table!")
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_len
func (self *luaState) Len(idx int) {
	val := self.stack.get(idx)
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
