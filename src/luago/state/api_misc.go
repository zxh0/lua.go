package state

import "runtime"
import "luago/binchunk"
import "luago/compiler"
import "luago/luanum"
import . "luago/api"

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_atpanic
func (self *luaState) AtPanic(panicf LuaGoFunction) LuaGoFunction {
	oldPanicf := self.panicf
	self.panicf = panicf
	return oldPanicf
}

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_register
func (self *luaState) Register(name string, f LuaGoFunction) {
	self.PushGoFunction(f)
	self.SetGlobal(name)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_getglobal
func (self *luaState) GetGlobal(name string) LuaType {
	global := self.registry.get(LUA_RIDX_GLOBALS).(*luaTable)
	val := global.get(name)
	self.stack.push(val)
	return typeOf(val)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setglobal
func (self *luaState) SetGlobal(name string) {
	global := self.registry.get(LUA_RIDX_GLOBALS).(*luaTable)
	val := self.stack.pop()
	global.put(name, val)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_getuservalue
func (self *luaState) GetUserValue(index int) LuaType {
	panic("todo!")
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_setuservalue
func (self *luaState) SetUserValue(index int) {
	// val := self.stack.pop()
	// ud := self.stack.get(index)
	panic("todo!")
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_close
func (self *luaState) Close() {
	//panic("todo!")
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_dump
func (self *luaState) Dump() {
	panic("todo!")
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
// http://www.lua.org/manual/5.3/manual.html#lua_load
func (self *luaState) Load(chunk []byte, chunkName, mode string) LuaThreadStatus {
	var proto *binchunk.Prototype
	if binchunk.IsBinaryChunk(chunk) {
		proto = binchunk.Undump(chunk)
	} else {
		proto = compiler.Compile(chunkName, string(chunk))
	}

	cl := newLuaClosure(proto)
	if len(proto.Upvalues) > 0 { // todo
		env := self.registry.get(LUA_RIDX_GLOBALS)
		cl.upvals[0] = &env
	}
	self.stack.push(cl)
	return LUA_OK
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

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_type
func (self *luaState) Type(index int) LuaType {
	if absIdx := self.stack.absIndex(index); absIdx > 0 {
		val := self.stack.get(index)
		return typeOf(val)
	} else {
		return LUA_TNONE
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_typename
func (self *luaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TUSERDATA:
		return "userdata"
	case LUA_TTHREAD:
		return "thread"
	default:
		panic("unknown data type!")
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_version
func (self *luaState) Version() float64 {
	panic("todo!")
}
