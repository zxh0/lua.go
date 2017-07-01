package state

import . "luago/api"

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnone
func (self *luaState) IsNone(index int) bool {
	absIdx := self.stack.absIndex(index)
	return absIdx == 0
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnoneornil
func (self *luaState) IsNoneOrNil(index int) bool {
	return self.IsNone(index) || self.IsNil(index)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnil
func (self *luaState) IsNil(index int) bool {
	val := self.stack.get(index)
	return val == nil
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isboolean
func (self *luaState) IsBoolean(index int) bool {
	val := self.stack.get(index)
	return typeOf(val) == LUA_TBOOLEAN
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isinteger
func (self *luaState) IsInteger(index int) bool {
	val := self.stack.get(index)
	return fullTypeOf(val) == LUA_TNUMINT
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnumber
func (self *luaState) IsNumber(index int) bool {
	_, ok := self.ToNumberX(index)
	return ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isstring
func (self *luaState) IsString(index int) bool {
	val := self.stack.get(index)
	t := typeOf(val)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_istable
func (self *luaState) IsTable(index int) bool {
	val := self.stack.get(index)
	return typeOf(val) == LUA_TTABLE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isfunction
func (self *luaState) IsFunction(index int) bool {
	val := self.stack.get(index)
	return typeOf(val) == LUA_TFUNCTION
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_iscfunction
func (self *luaState) IsGoFunction(index int) bool {
	val := self.stack.get(index)
	t := fullTypeOf(val)
	return t == LUA_TLGF || t == LUA_TGCL
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isuserdata
// http://www.lua.org/manual/5.3/manual.html#lua_islightuserdata
func (self *luaState) IsUserData(index int) bool {
	val := self.stack.get(index)
	return typeOf(val) == LUA_TUSERDATA
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isthread
func (self *luaState) IsThread(index int) bool {
	val := self.stack.get(index)
	return typeOf(val) == LUA_TTHREAD
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isyieldable
func (self *luaState) IsYieldable() bool {
	panic("todo!")
}
