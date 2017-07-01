package state

import . "luago/api"

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushnil
func (self *luaState) PushNil() {
	self.stack.push(nil)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushboolean
func (self *luaState) PushBoolean(b bool) {
	self.stack.push(b)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushinteger
func (self *luaState) PushInteger(n int64) {
	self.stack.push(n)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushnumber
func (self *luaState) PushNumber(n float64) {
	self.stack.push(n)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_pushstring
func (self *luaState) PushString(s string) {
	self.stack.push(s)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushcfunction
func (self *luaState) PushGoFunction(f LuaGoFunction) {
	self.stack.push(f)
}

// [-n, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_pushcclosure
func (self *luaState) PushGoClosure(fn LuaGoFunction, n int) {
	if n == 0 {
		self.stack.push(fn)
	} else { // closure
		vals := self.stack.popN(n)
		closure := &goClosure{fn, vals}
		self.stack.push(closure)
	}
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushlightuserdata
func (self *luaState) PushUserData(d LuaUserData) {
	ud := &userData{data: d}
	self.stack.push(ud)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_pushfstring
func (self *luaState) PushFString(fmt string) {
	panic("todo!")
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_pushvfstring
func (self *luaState) PushVFString() {
	panic("todo!")
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushglobaltable
func (self *luaState) PushGlobalTable() {
	global := self.registry.get(LUA_RIDX_GLOBALS)
	self.stack.push(global)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushthread
func (self *luaState) PushThread() bool {
	panic("todo!")
}
