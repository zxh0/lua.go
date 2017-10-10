package state

import . "luago/api"

// [-2, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_settable
func (self *luaState) SetTable(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()
	self.setTable(t, k, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setfield
func (self *luaState) SetField(idx int, k string) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, k, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_seti
func (self *luaState) SetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v, false)
}

// [-2, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawset
func (self *luaState) RawSet(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()
	self.setTable(t, k, v, true)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawseti
func (self *luaState) RawSetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v, true)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawsetp
func (self *luaState) RawSetP(idx int, p UserData) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, p, v, true)
}

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
	self.setTable(t, name, v, false)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_setmetatable
func (self *luaState) SetMetatable(idx int) {
	val := self.stack.get(idx)
	mtVal := self.stack.pop()

	if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, self)
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

// t[k]=v
func (self *luaState) setTable(t, k, v luaValue, raw bool) {
	if tbl, ok := t.(*luaTable); ok {
		if raw || !tbl.hasMetafield("__newindex") || tbl.get(k) != nil {
			tbl.put(k, v)
			return
		}
	} else if raw {
		panic("not table!")
	}

	if mf := getMetafield(t, "__newindex", self); mf != nil {
		switch x := mf.(type) {
		case *luaTable:
			self.setTable(x, k, v, true)
			return
		case *closure, GoFunction:
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
