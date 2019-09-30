package state

import (
	. "github.com/zxh0/lua.go/api"
)

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_newtable
func (state *luaState) NewTable() {
	state.CreateTable(0, 0)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_createtable
func (state *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	state.stack.push(t)
}

// [-1, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_gettable
func (state *luaState) GetTable(idx int) LuaType {
	t := state.stack.get(idx)
	k := state.stack.pop()
	return state.getTable(t, k, 1)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_getfield
func (state *luaState) GetField(idx int, k string) LuaType {
	t := state.stack.get(idx)
	return state.getTable(t, k, 1)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_geti
func (state *luaState) GetI(idx int, i int64) LuaType {
	t := state.stack.get(idx)
	return state.getTable(t, i, 1)
}

// [-1, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawget
func (state *luaState) RawGet(idx int) LuaType {
	t := state.stack.get(idx)
	k := state.stack.pop()
	return state.getTable(t, k, 0)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawgeti
func (state *luaState) RawGetI(idx int, i int64) LuaType {
	t := state.stack.get(idx)
	return state.getTable(t, i, 0)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawgetp
func (state *luaState) RawGetP(idx int, p UserData) LuaType {
	t := state.stack.get(idx)
	return state.getTable(t, p, 0)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_getglobal
func (state *luaState) GetGlobal(name string) LuaType {
	t := state.registry.get(LUA_RIDX_GLOBALS)
	return state.getTable(t, name, 1)
}

// [-0, +(0|1), –]
// http://www.lua.org/manual/5.3/manual.html#lua_getmetatable
func (state *luaState) GetMetatable(idx int) bool {
	val := state.stack.get(idx)

	if mt := getMetatable(val, state); mt != nil {
		state.stack.push(mt)
		return true
	} else {
		return false
	}
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_getuservalue
func (state *luaState) GetUserValue(idx int) LuaType {
	panic("todo!")
}

// push(t[k])
func (state *luaState) getTable(t, k luaValue, mtLv int) LuaType {
	if mtLv > MAXTAGLOOP {
		panic("'__index' chain too long; possible loop")
	}

	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		if mtLv == 0 || v != nil || !tbl.hasMetafield("__index") {
			state.stack.push(v)
			return typeOf(v)
		}
	}

	if mtLv > 0 {
		if mf := getMetafield(t, "__index", state); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				return state.getTable(x, k, mtLv+1)
			case *closure:
				state.stack.push(mf)
				state.stack.push(t)
				state.stack.push(k)
				state.Call(2, 1)
				v := state.stack.get(-1)
				return typeOf(v)
			}
		}
	}

	panic("not a table!") // todo
}
