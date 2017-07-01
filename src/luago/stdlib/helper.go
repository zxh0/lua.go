package stdlib

import "fmt"
import . "luago/api"

// todo: remove?
func getOptionalBoolArg(ls LuaState, idx int, defaultVal bool) bool {
	if ls.GetTop() >= idx {
		return ls.ToBoolean(idx)
	}
	return defaultVal
}

// todo
func castToString(ls LuaState, idx int) string {
	switch ls.Type(idx) {
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		if ls.ToBoolean(idx) {
			return "true"
		} else {
			return "false"
		}
	case LUA_TSTRING, LUA_TNUMBER:
		return ls.CheckString(idx)
	case LUA_TTABLE:
		return fmt.Sprintf("table: %p", ls.ToPointer(idx))
	case LUA_TFUNCTION:
		return fmt.Sprintf("function: %p", ls.ToPointer(idx))
	}

	// TODO
	return ""
}
