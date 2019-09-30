package state

import (
	"fmt"

	. "github.com/zxh0/lua.go/api"
)

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawlen
func (state *luaState) RawLen(idx int) uint {
	val := state.stack.get(idx)
	switch x := val.(type) {
	case string:
		return uint(len(x))
	case *luaTable:
		return uint(x.len())
	default:
		return 0
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_typename
// lua-5.3.4/src/lapi.c#lua_typename()
func (state *luaState) TypeName(tp LuaType) string {
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
		panic("unreachable!")
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_type
// lua-5.3.4/src/lapi.c#lua_type()
func (state *luaState) Type(idx int) LuaType {
	if state.stack.isValid(idx) {
		val := state.stack.get(idx)
		return typeOf(val)
	} else {
		return LUA_TNONE
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnone
// lua-5.3.4/src/lua.h#lua_isnone()
func (state *luaState) IsNone(idx int) bool {
	return state.Type(idx) == LUA_TNONE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnil
// lua-5.3.4/src/lua.h#lua_isnil()
func (state *luaState) IsNil(idx int) bool {
	return state.Type(idx) == LUA_TNIL
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnoneornil
// lua-5.3.4/src/lua.h#lua_isnoneornil()
func (state *luaState) IsNoneOrNil(idx int) bool {
	return state.Type(idx) <= 0
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isboolean
// lua-5.3.4/src/lua.h#lua_isboolean()
func (state *luaState) IsBoolean(idx int) bool {
	return state.Type(idx) == LUA_TBOOLEAN
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_istable
// lua-5.3.4/src/lua.h#lua_istable()
func (state *luaState) IsTable(idx int) bool {
	return state.Type(idx) == LUA_TTABLE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isfunction
// lua-5.3.4/src/lua.h#lua_isfunction()
func (state *luaState) IsFunction(idx int) bool {
	return state.Type(idx) == LUA_TFUNCTION
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isthread
// lua-5.3.4/src/lua.h#lua_isthread()
func (state *luaState) IsThread(idx int) bool {
	return state.Type(idx) == LUA_TTHREAD
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isuserdata
// http://www.lua.org/manual/5.3/manual.html#lua_islightuserdata
// lua-5.3.4/src/lapi.c#lua_isuserdata()
func (state *luaState) IsUserData(idx int) bool {
	return state.Type(idx) == LUA_TUSERDATA
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isstring
// lua-5.3.4/src/lapi.c#lua_isstring()
func (state *luaState) IsString(idx int) bool {
	t := state.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnumber
// lua-5.3.4/src/lapi.c#lua_isnumber()
func (state *luaState) IsNumber(idx int) bool {
	_, ok := state.ToNumberX(idx)
	return ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isinteger
// lua-5.3.4/src/lapi.c#lua_isinteger()
func (state *luaState) IsInteger(idx int) bool {
	val := state.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_iscfunction
// lua-5.3.4/src/lapi.c#lua_iscfunction()
func (state *luaState) IsGoFunction(idx int) bool {
	val := state.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc != nil
	}
	return false
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_toboolean
func (state *luaState) ToBoolean(idx int) bool {
	val := state.stack.get(idx)
	return convertToBoolean(val)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tointeger
func (state *luaState) ToInteger(idx int) int64 {
	i, _ := state.ToIntegerX(idx)
	return i
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tointegerx
func (state *luaState) ToIntegerX(idx int) (int64, bool) {
	val := state.stack.get(idx)
	return convertToInteger(val)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tonumber
func (state *luaState) ToNumber(idx int) float64 {
	n, _ := state.ToNumberX(idx)
	return n
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tonumberx
func (state *luaState) ToNumberX(idx int) (float64, bool) {
	val := state.stack.get(idx)
	return convertToFloat(val)
}

// [-0, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_tostring
// http://www.lua.org/manual/5.3/manual.html#lua_tolstring
// lua-5.3.4/src/lua.h#lua_tostring()
// lua-5.3.4/src/lapi.c#lua_tolstring()
func (state *luaState) ToString(idx int) string {
	s, _ := state.ToStringX(idx)
	return s
}

func (state *luaState) ToStringX(idx int) (string, bool) {
	val := state.stack.get(idx)

	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x) // todo
		state.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tocfunction
func (state *luaState) ToGoFunction(idx int) GoFunction {
	val := state.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc
	}
	return nil
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tothread
func (state *luaState) ToThread(idx int) LuaState {
	val := state.stack.get(idx)
	if val != nil {
		if ls, ok := val.(*luaState); ok {
			return ls
		}
	}
	return nil
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_touserdata
func (state *luaState) ToUserData(idx int) UserData {
	val := state.stack.get(idx)
	if val != nil {
		if ud, ok := val.(*userdata); ok {
			return ud.data
		}
	}
	return nil
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_topointer
func (state *luaState) ToPointer(idx int) interface{} {
	val := state.stack.get(idx)
	if val == nil || typeOf(val) < LUA_TTABLE {
		return nil
	} else {
		return val
	}
}
