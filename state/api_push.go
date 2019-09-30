package state

import (
	"fmt"

	. "github.com/zxh0/lua.go/api"
)

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushnil
func (state *luaState) PushNil() {
	state.stack.push(nil)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushboolean
func (state *luaState) PushBoolean(b bool) {
	state.stack.push(b)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushinteger
func (state *luaState) PushInteger(n int64) {
	state.stack.push(n)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushnumber
func (state *luaState) PushNumber(n float64) {
	state.stack.push(n)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_pushstring
func (state *luaState) PushString(s string) {
	state.stack.push(s)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_pushfstring
func (state *luaState) PushFString(fmtStr string, a ...interface{}) string {
	str := fmt.Sprintf(fmtStr, a...)
	state.stack.push(str)
	return str
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushcfunction
func (state *luaState) PushGoFunction(f GoFunction) {
	state.stack.push(newGoClosure(f, 0))
}

// [-n, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_pushcclosure
func (state *luaState) PushGoClosure(f GoFunction, n int) {
	closure := newGoClosure(f, n)
	for i := n; i > 0; i-- {
		val := state.stack.pop()
		closure.upvals[i-1] = &upvalue{&val}
	}
	state.stack.push(closure)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushlightuserdata
func (state *luaState) PushLightUserData(d UserData) {
	ud := &userdata{data: d}
	state.stack.push(ud)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushthread
func (state *luaState) PushThread() bool {
	state.stack.push(state)
	return state.isMainThread()
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushglobaltable
func (state *luaState) PushGlobalTable() {
	global := state.registry.get(LUA_RIDX_GLOBALS)
	state.stack.push(global)
}
