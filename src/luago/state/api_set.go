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
	t := self.registry.get(LUA_RIDX_GLOBALS)
	v := self.stack.pop()
	self._setTable(t, name, v, false)
}

// [-2, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_settable
func (self *luaState) SetTable(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()
	self._setTable(t, k, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setfield
func (self *luaState) SetField(idx int, k string) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self._setTable(t, k, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_seti
func (self *luaState) SetI(idx int, n int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self._setTable(t, n, v, false)
}

// [-2, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawset
func (self *luaState) RawSet(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()
	self._setTable(t, k, v, true)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawseti
func (self *luaState) RawSetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self._setTable(t, i, v, true)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawsetp
func (self *luaState) RawSetP(idx int, p UserData) {
	t := self.stack.get(idx)
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
func (self *luaState) SetMetatable(idx int) {
	val := self.stack.get(idx)
	mtVal := self.stack.pop()

	if mt, ok := mtVal.(*luaTable); ok {
		self.setMetatable(val, mt)
	} else {
		panic("not table: " + valToString(mtVal)) // todo
	}
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_setuservalue
func (self *luaState) SetUserValue(idx int) {
	// val := self.stack.pop()
	// ud := self.stack.get(idx)
	panic("todo!")
}
