package state

import (
	. "github.com/zxh0/lua.go/api"
)

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_gettop
// lua-5.3.4/src/lapi.c#lua_gettop()
func (state *luaState) GetTop() int {
	return state.stack.top
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_absindex
// lua-5.3.4/src/lapi.c#lua_absindex()
func (state *luaState) AbsIndex(idx int) int {
	return state.stack.absIndex(idx)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_upvalueindex
func (state *luaState) UpvalueIndex(idx int) int {
	return LUA_REGISTRYINDEX - idx
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_checkstack
// lua-5.3.4/src/lapi.c#lua_checkstack()
func (state *luaState) CheckStack(n int) bool {
	state.stack.check(n)
	return true // never fails
}

// [-n, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pop
// lua-5.3.4/src/lua.h#lua_pop()
func (state *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		state.stack.pop()
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_copy
// lua-5.3.4/src/lapi.c#lua_copy()
func (state *luaState) Copy(fromIdx, toIdx int) {
	val := state.stack.get(fromIdx)
	state.stack.set(toIdx, val)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushvalue
// lua-5.3.4/src/lapi.c#lua_pushvalue()
func (state *luaState) PushValue(idx int) {
	val := state.stack.get(idx)
	state.stack.push(val)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_replace
// lua-5.3.4/src/lua.h#lua_replace()
func (state *luaState) Replace(idx int) {
	val := state.stack.pop()
	state.stack.set(idx, val)
}

// [-1, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_insert
// lua-5.3.4/src/lua.h#lua_insert()
func (state *luaState) Insert(idx int) {
	state.Rotate(idx, 1)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_remove
// lua-5.3.4/src/lua.h#lua_remove()
func (state *luaState) Remove(idx int) {
	state.Rotate(idx, -1)
	state.Pop(1)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rotate
// lua-5.3.4/src/lapi.c#lua_rotate()
func (state *luaState) Rotate(idx, n int) {
	t := state.stack.top - 1           /* end of stack segment being rotated */
	p := state.stack.absIndex(idx) - 1 /* start of segment */
	var m int                          /* end of prefix */
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	state.stack.reverse(p, m)   /* reverse the prefix with length 'n' */
	state.stack.reverse(m+1, t) /* reverse the suffix */
	state.stack.reverse(p, t)   /* reverse the entire segment */
}

// [-?, +?, –]
// http://www.lua.org/manual/5.3/manual.html#lua_settop
// lua-5.3.4/src/lapi.c#lua_settop()
func (state *luaState) SetTop(idx int) {
	newTop := state.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}

	n := state.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			state.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			state.stack.push(nil)
		}
	}
}

// [-?, +?, –]
// http://www.lua.org/manual/5.3/manual.html#lua_xmove
// lua-5.3.4/src/lapi.c#lua_xmove()
func (state *luaState) XMove(to LuaState, n int) {
	vals := state.stack.popN(n)
	to.(*luaState).stack.pushN(vals, n)
}
