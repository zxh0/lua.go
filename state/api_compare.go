package state

import (
	. "github.com/zxh0/lua.go/api"
)

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawequal
func (state *luaState) RawEqual(idx1, idx2 int) bool {
	if !state.stack.isValid(idx1) || !state.stack.isValid(idx2) {
		return false
	}

	a := state.stack.get(idx1)
	b := state.stack.get(idx2)
	return state.eq(a, b, true)
}

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_compare
func (state *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
	if !state.stack.isValid(idx1) || !state.stack.isValid(idx2) {
		return false
	}

	a := state.stack.get(idx1)
	b := state.stack.get(idx2)
	switch op {
	case LUA_OPEQ:
		return state.eq(a, b, false)
	case LUA_OPLT:
		return state.lt(a, b)
	case LUA_OPLE:
		return state.le(a, b)
	default:
		panic("invalid compare op!")
	}
}

func (state *luaState) eq(a, b luaValue, raw bool) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
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
	case *luaTable:
		if y, ok := b.(*luaTable); ok && x != y && !raw {
			if result, ok := callMetamethod(x, y, "__eq", state); ok {
				return convertToBoolean(result)
			}
		}
		return a == b
	case *userdata:
		if y, ok := b.(*userdata); ok {
			return x.data == y.data
		}
		return false
	default:
		return a == b
	}
}

func (state *luaState) lt(a, b luaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x < y
		case int64:
			return x < float64(y)
		}
	}

	if result, ok := callMetamethod(a, b, "__lt", state); ok {
		return convertToBoolean(result)
	}
	typeName1 := state.TypeName(typeOf(a))
	typeName2 := state.TypeName(typeOf(b))
	panic("attempt to compare " + typeName1 + " with " + typeName2)
}

func (state *luaState) le(a, b luaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x <= y
		case int64:
			return x <= float64(y)
		}
	}

	if result, ok := callMetamethod(a, b, "__le", state); ok {
		return convertToBoolean(result)
	}
	if result, ok := callMetamethod(b, a, "__lt", state); ok {
		return !convertToBoolean(result)
	}
	typeName1 := state.TypeName(typeOf(a))
	typeName2 := state.TypeName(typeOf(b))
	panic("attempt to compare " + typeName1 + " with " + typeName2)
}
