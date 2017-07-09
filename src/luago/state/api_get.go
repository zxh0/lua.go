package state

import . "luago/api"

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_getglobal
func (self *luaState) GetGlobal(name string) LuaType {
	global := self.registry.get(LUA_RIDX_GLOBALS).(*luaTable)
	val := global.get(name)
	self.stack.push(val)
	return typeOf(val)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_createtable
func (self *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	self.stack.push(t)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_newtable
func (self *luaState) NewTable() {
	self.CreateTable(0, 0)
}

// [-1, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_gettable
func (self *luaState) GetTable(index int) LuaType {
	t := self.stack.get(index)
	k := self.stack.pop()
	return self._getTable(t, k, false)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_getfield
func (self *luaState) GetField(index int, k string) LuaType {
	t := self.stack.get(index)
	return self._getTable(t, k, false)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_geti
func (self *luaState) GetI(index int, i int64) LuaType {
	t := self.stack.get(index)
	return self._getTable(t, i, false)
}

// [-1, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawget
func (self *luaState) RawGet(index int) LuaType {
	t := self.stack.get(index)
	k := self.stack.pop()
	return self._getTable(t, k, true)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawgeti
func (self *luaState) RawGetI(index int, n int64) LuaType {
	t := self.stack.get(index)
	return self._getTable(t, n, true)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawgetp
func (self *luaState) RawGetP(index int, p UserData) LuaType {
	t := self.stack.get(index)
	return self._getTable(t, p, true)
}

// push(t[k])
func (self *luaState) _getTable(t, k luaValue, raw bool) LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		if v != nil || raw || !tbl.hasMetaField("__index") {
			self.stack.push(v)
			return typeOf(v)
		}
	} else if raw {
		panic("not table!")
	}

	if mf := self.getMetaField(t, "__index"); mf != nil {
		switch x := mf.(type) {
		case *luaTable:
			return self._getTable(x, k, true)
		case *luaClosure, GoFunction:
			self.stack.push(mf)
			self.stack.push(t)
			self.stack.push(k)
			self.Call(2, 1)
			v := self.stack.get(-1)
			return typeOf(v)
		}
	}

	panic("not table!") // todo
}

// [-0, +(0|1), –]
// http://www.lua.org/manual/5.3/manual.html#lua_getmetatable
func (self *luaState) GetMetaTable(index int) bool {
	val := self.stack.get(index)

	if mt := self.getMetaTable(val); mt != nil {
		self.stack.push(mt)
		return true
	} else {
		return false
	}
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_getuservalue
func (self *luaState) GetUserValue(index int) LuaType {
	panic("todo!")
}
