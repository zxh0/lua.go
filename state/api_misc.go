package state

import (
	"runtime"

	. "github.com/zxh0/lua.go/api"
	"github.com/zxh0/lua.go/number"
)

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_close
// lua-5.3.4/src/lstate.c#lua_close
func (state *luaState) Close() {
	// todo
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_atpanic
// lua-5.3.4/src/lapi.c#lua_atpanic
func (state *luaState) AtPanic(panicf GoFunction) GoFunction {
	oldPanicf := state.panicf
	state.panicf = panicf
	return oldPanicf
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_version
// lua-5.3.4/src/lapi.c#lua_version
func (state *luaState) Version() float64 {
	return LUA_VERSION_NUM
}

// [-1, +0, v]
// http://www.lua.org/manual/5.3/manual.html#lua_error
func (state *luaState) Error() int {
	err := newLuaTable(0, 1)
	err.put("_ERR", state.stack.pop())
	panic(err)
}

// [-0, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_gc
func (state *luaState) GC(what, data int) int {
	runtime.GC()
	return 0
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_stringtonumber
func (state *luaState) StringToNumber(s string) bool {
	if n, ok := number.ParseInteger(s); ok {
		state.PushInteger(n)
		return true
	}
	if n, ok := number.ParseFloat(s); ok {
		state.PushNumber(n)
		return true
	}
	return false
}

// [-1, +(2|0), e]
// http://www.lua.org/manual/5.3/manual.html#lua_next
func (state *luaState) Next(idx int) bool {
	val := state.stack.get(idx)
	if t, ok := val.(*luaTable); ok {
		key := state.stack.pop()
		if nextKey := t.nextKey(key); nextKey != nil {
			state.stack.push(nextKey)
			state.stack.push(t.get(nextKey))
			return true
		}
		return false
	}
	panic("table expected!")
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_len
func (state *luaState) Len(idx int) {
	val := state.stack.get(idx)
	if s, ok := val.(string); ok {
		state.stack.push(int64(len(s)))
	} else if result, ok := callMetamethod(val, val, "__len", state); ok {
		state.stack.push(result)
	} else if t, ok := val.(*luaTable); ok {
		state.stack.push(int64(t.len()))
	} else {
		typeName := state.TypeName(typeOf(val))
		panic("attempt to get length of a " + typeName + " value")
	}
}

// [-n, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_concat
func (state *luaState) Concat(n int) {
	if n == 0 {
		state.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if state.IsString(-1) && state.IsString(-2) {
				s2 := state.ToString(-1)
				s1 := state.ToString(-2)
				state.stack.pop()
				state.stack.pop()
				state.stack.push(s1 + s2)
				continue
			}

			b := state.stack.pop()
			a := state.stack.pop()
			if result, ok := callMetamethod(a, b, "__concat", state); ok {
				state.stack.push(result)
				continue
			}

			var typeName string
			if _, ok := convertToFloat(a); !ok {
				typeName = state.TypeName(typeOf(a))
			} else {
				typeName = state.TypeName(typeOf(b))
			}
			panic("attempt to concatenate a " + typeName + " value")
		}
	}
	// n == 1, do nothing
}
