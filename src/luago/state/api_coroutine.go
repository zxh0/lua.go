package state

import . "luago/api"

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_newthread
func (self *luaState) NewThread() LuaState {
	return NewLuaState() // todo
}

// http://www.lua.org/manual/5.3/manual.html#lua_status
func (self *luaState) Status() int {
	panic("todo!")
}

// [-?, +?, –]
// http://www.lua.org/manual/5.3/manual.html#lua_resume
func (self *luaState) Resume(from LuaState, nArgs int) {
	panic("todo!")
}

// [-?, +?, e]
// http://www.lua.org/manual/5.3/manual.html#lua_yield
func (self *luaState) Yield(nResults int) int {
	panic("todo!")
}

// [-?, +?, e]
// http://www.lua.org/manual/5.3/manual.html#lua_yieldk
func (self *luaState) YieldK() {
	panic("todo!")
}
