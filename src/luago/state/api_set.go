package state

import . "luago/api"

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_register
func (self *luaState) Register(name string, f GoFunction) {
	self.PushGoFunction(f)
	self.SetGlobal(name)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setglobal
func (self *luaState) SetGlobal(name string) {
	global := self.registry.get(LUA_RIDX_GLOBALS).(*luaTable)
	val := self.stack.pop()
	global.put(name, val)
}

// [-2, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_settable
func (self *luaState) SetTable(index int) {
	t := self.stack.get(index)
	v := self.stack.pop()
	k := self.stack.pop()
	self._setTable(t, k, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setfield
func (self *luaState) SetField(index int, k string) {
	t := self.stack.get(index)
	v := self.stack.pop()
	self._setTable(t, k, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_seti
func (self *luaState) SetI(index int, n int64) {
	t := self.stack.get(index)
	v := self.stack.pop()
	self._setTable(t, n, v, false)
}

// [-2, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawset
func (self *luaState) RawSet(index int) {
	t := self.stack.get(index)
	v := self.stack.pop()
	k := self.stack.pop()
	self._setTable(t, k, v, true)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawseti
func (self *luaState) RawSetI(index int, i int64) {
	t := self.stack.get(index)
	v := self.stack.pop()
	self._setTable(t, i, v, true)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawsetp
func (self *luaState) RawSetP(index int, p UserData) {
	t := self.stack.get(index)
	v := self.stack.pop()
	self._setTable(t, p, v, true)
}

// t[k]=v
func (self *luaState) _setTable(t, k, v luaValue, raw bool) {
	if tbl, ok := t.(*luaTable); ok {
		if raw || !tbl.hasMetaField("__newindex") || tbl.get(k) != nil {
			tbl.put(k, v)
			return
		}
	} else if raw {
		panic("not table!")
	}

	if mf := self.getMetaField(t, "__newindex"); mf != nil {
		switch x := mf.(type) {
		case *luaTable:
			self._setTable(x, k, v, true)
			return
		case *luaClosure, GoFunction:
			self.stack.push(mf)
			self.stack.push(t)
			self.stack.push(k)
			self.stack.push(v)
			self.Call(3, 0)
			return
		}
	}

	panic("not table!")
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_setmetatable
func (self *luaState) SetMetaTable(index int) {
	val := self.stack.get(index)
	mtVal := self.stack.pop()

	if mt, ok := mtVal.(*luaTable); ok {
		self.setMetaTable(val, mt)
	} else {
		panic("not table: " + valToString(mtVal)) // todo
	}
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_setuservalue
func (self *luaState) SetUserValue(index int) {
	// val := self.stack.pop()
	// ud := self.stack.get(index)
	panic("todo!")
}
