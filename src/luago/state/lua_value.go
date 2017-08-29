package state

import . "luago/api"
import "luago/luanum"

type luaValue interface{}

func typeOf(val luaValue) LuaType {
	return fullTypeOf(val) & 0x0F
}

func fullTypeOf(val luaValue) LuaType {
	switch val.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case int64:
		return LUA_TNUMINT
	case float64:
		return LUA_TNUMFLT
	case string:
		return LUA_TSTRING
	case *luaTable:
		return LUA_TTABLE
	case *luaState:
		return LUA_TTHREAD
	case *userData:
		return LUA_TUSERDATA
	case *luaClosure:
		return LUA_TLCL
	case *goClosure:
		return LUA_TGCL
	case GoFunction:
		return LUA_TLGF
	default:
		panic("unkonwn type!")
	}
}

func convertToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

// http://www.lua.org/manual/5.3/manual.html#3.4.3
func convertToNumber(val luaValue) (float64, bool) {
	switch x := val.(type) {
	case int64:
		return float64(x), true
	case float64:
		return x, true
	case string:
		return luanum.ParseFloat(x)
	default:
		return 0, false
	}
}

// http://www.lua.org/manual/5.3/manual.html#3.4.3
func convertToInteger(val luaValue) (int64, bool) {
	switch x := val.(type) {
	case int64:
		return x, true
	case float64:
		return luanum.FloatToInteger(x)
	case string:
		return _stringToInteger(x)
	default:
		return 0, false
	}
}

func _stringToInteger(s string) (int64, bool) {
	if i, ok := luanum.ParseInteger(s, 10); ok {
		return i, true
	}
	if f, ok := luanum.ParseFloat(s); ok {
		return luanum.FloatToInteger(f)
	}
	return 0, false
}
