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

	a := self.stack.get(idx1)
	b := self.stack.get(idx2)
	return self.eq(a, b, true)
}

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_compare
func (self *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
	a := self.stack.get(idx1)
	b := self.stack.get(idx2)
	switch op {
	case LUA_OPEQ:
		return self.eq(a, b, false)
	case LUA_OPLT:
		return self.lt(a, b)
	case LUA_OPLE:
		return self.le(a, b)
	default:
		panic("invalid compare op!")
	}
}

func (self *luaState) eq(a, b luaValue, raw bool) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	case string:
		y, ok := b.(string)
		return ok && x == y
	case GoFunction:
		// todo: funcs are uncomparable!
		if y, ok := b.(GoFunction); ok {
			return fmt.Sprintf("%p", x) == fmt.Sprintf("%p", y)
		} else {
			return false
		}
	case *luaTable:
		if raw {
			return a == b
		}
		if y, ok := b.(*luaTable); ok {
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
		return a == b
	}
}

func (self *luaState) lt(a, b luaValue) bool {
	switch x := a.(type) {
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x < y
		case int64:
			return x < float64(y)
		default:
			return false
		}
	case string:
		y, ok := b.(string)
		return ok && x < y
	default:
		if result, ok := callMetamethod(a, b, "__lt", self); ok {
			return convertToBoolean(result)
		} else {
			panic("todo: __lt!")
		}
	}
}

func (self *luaState) le(a, b luaValue) bool {
	switch x := a.(type) {
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x <= y
		case int64:
			return x <= float64(y)
		default:
			return false
		}
	case string:
		y, ok := b.(string)
		return ok && x <= y
	default:
		if result, ok := callMetamethod(a, b, "__le", self); ok {
			return convertToBoolean(result)
		} else if result, ok := callMetamethod(b, a, "__lt", self); ok {
			return !convertToBoolean(result)
		} else {
			panic("todo: __le!")
		}
	}
}
