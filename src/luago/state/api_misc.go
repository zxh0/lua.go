package state

import "runtime"
import "luago/number"
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
	err := newLuaTable(0, 1)
	err.put("_ERR", self.stack.pop())
	panic(err)
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
	if n, ok := number.ParseInteger(s); ok {
		self.PushInteger(n)
		return true
	}
	if n, ok := number.ParseFloat(s); ok {
		self.PushNumber(n)
		return true
	}
	return false
}

// [-1, +(2|0), e]
// http://www.lua.org/manual/5.3/manual.html#lua_next
func (self *luaState) Next(idx int) bool {
	val := self.stack.get(idx)
	if t, ok := val.(*luaTable); ok {
		key := self.stack.pop()
		if nextKey, nextVal := t.next(key); nextKey != nil {
			self.stack.push(nextKey)
			self.stack.push(nextVal)
			return true
		} else {
			return false
		}
	}
	panic("not a table!")
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_len
func (self *luaState) Len(idx int) {
	val := self.stack.get(idx)
	if s, ok := val.(string); ok {
		self.stack.push(int64(len(s)))
	} else if result, ok := callMetamethod(val, val, "__len", self); ok {
		self.stack.push(result)
	} else if t, ok := val.(*luaTable); ok {
		self.stack.push(int64(t.len()))
	} else {
		typeName := self.TypeName(typeOf(val))
		panic("attempt to get length of a " + typeName + " value")
	}
}

// [-n, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_concat
func (self *luaState) Concat(n int) {
	if n == 0 {
		self.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if self.IsString(-1) && self.IsString(-2) {
				s2, _ := self.ToString(-1)
				s1, _ := self.ToString(-2)
				self.stack.pop()
				self.stack.pop()
				self.stack.push(s1 + s2)
				continue
			}

			b := self.stack.pop()
			a := self.stack.pop()
			if result, ok := callMetamethod(a, b, "__concat", self); ok {
				self.stack.push(result)
				continue
			}

			var typeName string
			if _, ok := convertToFloat(a); !ok {
				typeName = self.TypeName(typeOf(a))
			} else {
				typeName = self.TypeName(typeOf(b))
			}
			panic("attempt to concatenate a " + typeName + " value")
		}
	}
	// n == 1, do nothing
}
