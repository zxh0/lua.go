package state

import . "luago/api"
import "luago/luanum"

type luaValue interface{}

/* typeOf */

func typeOf(val luaValue) LuaType {
	switch val.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case int64, float64:
		return LUA_TNUMBER
	case string:
		return LUA_TSTRING
	case *luaTable:
		return LUA_TTABLE
	case *luaState:
		return LUA_TTHREAD
	case *userData:
		return LUA_TUSERDATA
	case *luaClosure, *goClosure, GoFunction:
		return LUA_TFUNCTION
	default:
		panic("unkonwn type!")
	}
}

/* convert */

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

/* metatable */

func getMetatable(val luaValue, ls *luaState) *luaTable {
	switch x := val.(type) {
	case nil:
		return ls.mtOfNil
	case bool:
		return ls.mtOfBool
	case int64, float64:
		return ls.mtOfNumber
	case string:
		return ls.mtOfString
	case *luaClosure, *goClosure, GoFunction:
		return ls.mtOfFunc
	case *luaTable:
		return x.metatable
	case *userData:
		return x.metatable
	default: // todo
		return nil
	}
}

func setMetatable(val luaValue, mt *luaTable, ls *luaState) {
	switch x := val.(type) {
	case nil:
		ls.mtOfNil = mt
	case bool:
		ls.mtOfBool = mt
	case int64, float64:
		ls.mtOfNumber = mt
	case string:
		ls.mtOfString = mt
	case *luaClosure, *goClosure, GoFunction:
		ls.mtOfFunc = mt
	case *luaTable:
		x.metatable = mt
	case *userData:
		x.metatable = mt
	default:
		// todo
	}
}

func getMetafield(val luaValue, fieldName string, ls *luaState) luaValue {
	if mt := getMetatable(val, ls); mt != nil {
		return mt.get(fieldName)
	} else {
		return nil
	}
}

func callMetamethod(a, b luaValue, mmName string, ls *luaState) (luaValue, bool) {
	var mm luaValue
	if mm = getMetafield(a, mmName, ls); mm == nil {
		if mm = getMetafield(b, mmName, ls); mm == nil {
			return nil, false
		}
	}

	ls.stack.check(4)
	ls.stack.push(mm)
	ls.stack.push(a)
	ls.stack.push(b)
	ls.Call(2, 1)
	return ls.stack.pop(), true
}
