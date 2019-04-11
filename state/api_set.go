package state

import . "github.com/zxh0/lua.go/api"

// [-2, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_settable
func (self *luaState) SetTable(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()
	self.setTable(t, k, v, 1)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setfield
func (self *luaState) SetField(idx int, k string) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, k, v, 1)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_seti
func (self *luaState) SetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v, 1)
}

// [-2, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawset
func (self *luaState) RawSet(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()
	self.setTable(t, k, v, 0)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawseti
func (self *luaState) RawSetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v, 0)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawsetp
func (self *luaState) RawSetP(idx int, p UserData) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, p, v, 0)
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
	self.setTable(t, name, v, 1)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_setmetatable
func (self *luaState) SetMetatable(idx int) {
	val := self.stack.get(idx)
	mtVal := self.stack.pop()

	if mtVal == nil {
		setMetatable(val, nil, self)
	} else if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, self)
	} else {
		panic("table expected!") // todo
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
func (self *luaState) setTable(t, k, v luaValue, mtLv int) {
	if mtLv > MAXTAGLOOP {
		panic("'__newindex' chain too long; possible loop")
	}

	if tbl, ok := t.(*luaTable); ok {
		if mtLv == 0 || tbl.get(k) != nil || !tbl.hasMetafield("__newindex") {
			tbl.put(k, v)
			return
		}
	}

	if mtLv > 0 {
		if mf := getMetafield(t, "__newindex", self); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				self.setTable(x, k, v, mtLv+1)
				return
			case *closure:
				self.stack.push(mf)
				self.stack.push(t)
				self.stack.push(k)
				self.stack.push(v)
				self.Call(3, 0)
				return
			}
		}
	}

	panic("not a table!") // todo
}
