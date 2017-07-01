package state

import . "luago/api"

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

// [-1, +(2|0), e]
// http://www.lua.org/manual/5.3/manual.html#lua_next
func (self *luaState) Next(index int) bool {
	t := self.stack.get(index)
	if tbl, ok := t.(*luaTable); ok {
		key := self.stack.pop()
		nextKey, nextVal := tbl.next(key)
		if nextKey != nil {
			self.stack.push(nextKey)
			self.stack.push(nextVal)
			return true
		} else {
			return false
		}
	}
	panic("not table!")
}

// [-0, +(0|1), –]
// http://www.lua.org/manual/5.3/manual.html#lua_getmetatable
func (self *luaState) GetMetaTable(index int) bool {
	val := self.stack.get(index)

	if mt := getMetaTable(val); mt != nil {
		self.stack.push(mt)
		return true
	} else {
		return false
	}
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_setmetatable
func (self *luaState) SetMetaTable(index int) {
	val := self.stack.get(index)
	mtVal := self.stack.pop()

	if mt, ok := mtVal.(*luaTable); ok {
		setMetaTable(val, mt)
	} else {
		panic("not table: " + valToString(mtVal)) // todo
	}
}

// [-1, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_gettable
func (self *luaState) GetTable(index int) LuaType {
	t := self.stack.get(index)
	k := self.stack.pop()
	return self._getTable(t, k, false)
}

// [-2, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_settable
func (self *luaState) SetTable(index int) {
	t := self.stack.get(index)
	v := self.stack.pop()
	k := self.stack.pop()
	self._setTable(t, k, v, false)
}

// [-1, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawget
func (self *luaState) RawGet(index int) LuaType {
	t := self.stack.get(index)
	k := self.stack.pop()
	return self._getTable(t, k, true)
}

// [-2, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawset
func (self *luaState) RawSet(index int) {
	t := self.stack.get(index)
	v := self.stack.pop()
	k := self.stack.pop()
	self._setTable(t, k, v, true)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_geti
func (self *luaState) GetI(index int, i int64) LuaType {
	t := self.stack.get(index)
	return self._getTable(t, i, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_seti
func (self *luaState) SetI(index int, n int64) {
	t := self.stack.get(index)
	v := self.stack.pop()
	self._setTable(t, n, v, false)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawgeti
func (self *luaState) RawGetI(index int, n int64) LuaType {
	t := self.stack.get(index)
	return self._getTable(t, n, true)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawseti
func (self *luaState) RawSetI(index int, i int64) {
	t := self.stack.get(index)
	v := self.stack.pop()
	self._setTable(t, i, v, true)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_getfield
func (self *luaState) GetField(index int, k string) LuaType {
	t := self.stack.get(index)
	return self._getTable(t, k, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setfield
func (self *luaState) SetField(index int, k string) {
	t := self.stack.get(index)
	v := self.stack.pop()
	self._setTable(t, k, v, false)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawgetp
func (self *luaState) RawGetP(index int, p LuaUserData) LuaType {
	t := self.stack.get(index)
	return self._getTable(t, p, true)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawsetp
func (self *luaState) RawSetP(index int, p LuaUserData) {
	t := self.stack.get(index)
	v := self.stack.pop()
	self._setTable(t, p, v, true)
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

	if mf := getMetaField(t, "__index"); mf != nil {
		switch x := mf.(type) {
		case *luaTable:
			return self._getTable(x, k, true)
		case *luaClosure, LuaGoFunction:
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

	if mf := getMetaField(t, "__newindex"); mf != nil {
		switch x := mf.(type) {
		case *luaTable:
			self._setTable(x, k, v, true)
			return
		case *luaClosure, LuaGoFunction:
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
