package state

import (
	. "github.com/zxh0/lua.go/api"
)

// [-2, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_settable
func (state *luaState) SetTable(idx int) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	k := state.stack.pop()
	state.setTable(t, k, v, 1)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setfield
func (state *luaState) SetField(idx int, k string) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, k, v, 1)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_seti
func (state *luaState) SetI(idx int, i int64) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, i, v, 1)
}

// [-2, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawset
func (state *luaState) RawSet(idx int) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	k := state.stack.pop()
	state.setTable(t, k, v, 0)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawseti
func (state *luaState) RawSetI(idx int, i int64) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, i, v, 0)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawsetp
func (state *luaState) RawSetP(idx int, p UserData) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, p, v, 0)
}

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_register
func (state *luaState) Register(name string, f GoFunction) {
	state.PushGoFunction(f)
	state.SetGlobal(name)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setglobal
func (state *luaState) SetGlobal(name string) {
	t := state.registry.get(LUA_RIDX_GLOBALS)
	v := state.stack.pop()
	state.setTable(t, name, v, 1)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_setmetatable
func (state *luaState) SetMetatable(idx int) {
	val := state.stack.get(idx)
	mtVal := state.stack.pop()

	if mtVal == nil {
		setMetatable(val, nil, state)
	} else if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, state)
	} else {
		panic("table expected!") // todo
	}
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_setuservalue
func (state *luaState) SetUserValue(idx int) {
	// val := state.stack.pop()
	// ud := state.stack.get(idx)
	panic("todo!")
}

// t[k]=v
func (state *luaState) setTable(t, k, v luaValue, mtLv int) {
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
		if mf := getMetafield(t, "__newindex", state); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				state.setTable(x, k, v, mtLv+1)
				return
			case *closure:
				state.stack.push(mf)
				state.stack.push(t)
				state.stack.push(k)
				state.stack.push(v)
				state.Call(3, 0)
				return
			}
		}
	}

	panic("not a table!") // todo
}
