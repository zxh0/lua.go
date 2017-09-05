package state

import "fmt"
import . "luago/api"

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawequal
func (self *luaState) RawEqual(idx1, idx2 int) bool {
	if self.stack.absIndex(idx1) == 0 ||
		self.stack.absIndex(idx2) == 0 {
		return false
	}

	val1 := self.stack.get(idx1)
	val2 := self.stack.get(idx2)
	return self.eq(val1, val2, true)
}

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_compare
func (self *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
	val1 := self.stack.get(idx1)
	val2 := self.stack.get(idx2)
	switch op {
	case LUA_OPEQ:
		return self.eq(val1, val2, false)
	case LUA_OPLT:
		return self.lt(val1, val2)
	case LUA_OPLE:
		return self.le(val1, val2)
	default:
		panic("invalid compare op!")
	}
}

func (self *luaState) eq(val1, val2 luaValue, raw bool) bool {
	switch x := val1.(type) {
	case nil:
		return val2 == nil
	case bool:
		y, ok := val2.(bool)
		return ok && x == y
	case int64:
		switch y := val2.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := val2.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	case string:
		y, ok := val2.(string)
		return ok && x == y
	case GoFunction:
		// todo: funcs are uncomparable!
		if y, ok := val2.(GoFunction); ok {
			return fmt.Sprintf("%p", x) == fmt.Sprintf("%p", y)
		} else {
			return false
		}
	case *luaTable:
		if raw {
			return val1 == val2
		}
		if y, ok := val2.(*luaTable); ok {
			if x == y {
				return true
			} else if result, ok := callMetamethod(x, y, "__eq", self); ok {
				return convertToBoolean(result)
			} else {
				return false
			}
		} else {
			return false
		}
	default:
		return val1 == val2
	}
}

func (self *luaState) lt(val1, val2 luaValue) bool {
	switch x := val1.(type) {
	case int64:
		switch y := val2.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		default:
			return false
		}
	case float64:
		switch y := val2.(type) {
		case float64:
			return x < y
		case int64:
			return x < float64(y)
		default:
			return false
		}
	case string:
		y, ok := val2.(string)
		return ok && x < y
	default:
		if result, ok := callMetamethod(val1, val2, "__lt", self); ok {
			return convertToBoolean(result)
		} else {
			panic("todo: __lt!")
		}
	}
}

func (self *luaState) le(val1, val2 luaValue) bool {
	switch x := val1.(type) {
	case int64:
		switch y := val2.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		default:
			return false
		}
	case float64:
		switch y := val2.(type) {
		case float64:
			return x <= y
		case int64:
			return x <= float64(y)
		default:
			return false
		}
	case string:
		y, ok := val2.(string)
		return ok && x <= y
	default:
		if result, ok := callMetamethod(val1, val2, "__le", self); ok {
			return convertToBoolean(result)
		} else if result, ok := callMetamethod(val2, val1, "__lt", self); ok {
			return !convertToBoolean(result)
		} else {
			panic("todo: __le!")
		}
	}
}
