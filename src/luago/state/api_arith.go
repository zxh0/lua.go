package state

import "math"
import . "luago/api"
import "luago/number"

// [-(2|1), +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_arith
func (self *luaState) Arith(op ArithOp) {
	var a, b luaValue
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		b = self.stack.pop()
	}
	a = self.stack.pop()

	switch op {
	case LUA_OPADD:
		self.add(a, b)
	case LUA_OPSUB:
		self.sub(a, b)
	case LUA_OPMUL:
		self.mul(a, b)
	case LUA_OPMOD:
		self.mod(a, b)
	case LUA_OPPOW:
		self.pow(a, b)
	case LUA_OPDIV:
		self.div(a, b)
	case LUA_OPIDIV:
		self.idiv(a, b)
	case LUA_OPBAND:
		self.band(a, b)
	case LUA_OPBOR:
		self.bor(a, b)
	case LUA_OPBXOR:
		self.bxor(a, b)
	case LUA_OPSHL:
		self.shl(a, b)
	case LUA_OPSHR:
		self.shr(a, b)
	case LUA_OPUNM:
		self.unm(a)
	case LUA_OPBNOT:
		self.bnot(a)
	default:
		panic("invalid arith op!")
	}
}

/* integer or float */

func (self *luaState) add(a, b luaValue) {
	if x, y, ok := _castToInteger(a, b); ok {
		self.stack.push(x + y)
	} else if x, y, ok := _convertToFloats(a, b); ok {
		self.stack.push(x + y)
	} else if result, ok := callMetamethod(a, b, "__add", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __add")
	}
}

func (self *luaState) sub(a, b luaValue) {
	if x, y, ok := _castToInteger(a, b); ok {
		self.stack.push(x - y)
	} else if x, y, ok := _convertToFloats(a, b); ok {
		self.stack.push(x - y)
	} else if result, ok := callMetamethod(a, b, "__sub", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __sub")
	}
}

func (self *luaState) mul(a, b luaValue) {
	if x, y, ok := _castToInteger(a, b); ok {
		self.stack.push(x * y)
	} else if x, y, ok := _convertToFloats(a, b); ok {
		self.stack.push(x * y)
	} else if result, ok := callMetamethod(a, b, "__mul", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __mul")
	}
}

func (self *luaState) idiv(a, b luaValue) {
	if x, y, ok := _castToInteger(a, b); ok {
		self.stack.push(number.IFloorDiv(x, y))
	} else if x, y, ok := _convertToFloats(a, b); ok {
		self.stack.push(number.FFloorDiv(x, y))
	} else if result, ok := callMetamethod(a, b, "__idiv", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __idiv")
	}
}

func (self *luaState) mod(a, b luaValue) {
	if x, y, ok := _castToInteger(a, b); ok {
		self.stack.push(number.IMod(x, y))
	} else if x, y, ok := _convertToFloats(a, b); ok {
		self.stack.push(number.FMod(x, y))
	} else if result, ok := callMetamethod(a, b, "__mod", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __mod")
	}
}

func (self *luaState) unm(a luaValue) {
	if x, ok := a.(int64); ok {
		self.stack.push(-x)
	} else if x, ok := convertToFloat(a); ok {
		self.stack.push(-x)
	} else if result, ok := callMetamethod(a, a, "__unm", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __unm!")
	}
}

/* float */

func (self *luaState) div(a, b luaValue) {
	if x, y, ok := _convertToFloats(a, b); ok {
		self.stack.push(x / y)
	} else if result, ok := callMetamethod(a, b, "__div", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __div")
	}
}

func (self *luaState) pow(a, b luaValue) {
	if x, y, ok := _convertToFloats(a, b); ok {
		self.stack.push(math.Pow(x, y))
	} else if result, ok := callMetamethod(a, b, "__pow", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __pow")
	}
}

/* bitwise */

func (self *luaState) band(a, b luaValue) {
	if x, y, ok := _convertToIntegers(a, b); ok {
		self.stack.push(x & y)
	} else if result, ok := callMetamethod(a, b, "__band", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __band")
	}
}

func (self *luaState) bor(a, b luaValue) {
	if x, y, ok := _convertToIntegers(a, b); ok {
		self.stack.push(x | y)
	} else if result, ok := callMetamethod(a, b, "__bor", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __bor")
	}
}

func (self *luaState) bxor(a, b luaValue) {
	if x, y, ok := _convertToIntegers(a, b); ok {
		self.stack.push(x ^ y)
	} else if result, ok := callMetamethod(a, b, "__bxor", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __bxor")
	}
}

func (self *luaState) shl(a, b luaValue) {
	if x, y, ok := _convertToIntegers(a, b); ok {
		self.stack.push(number.ShiftLeft(x, y))
	} else if result, ok := callMetamethod(a, b, "__shl", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __shl")
	}
}

func (self *luaState) shr(a, b luaValue) {
	if x, y, ok := _convertToIntegers(a, b); ok {
		self.stack.push(number.ShiftRight(x, y))
	} else if result, ok := callMetamethod(a, b, "__shr", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __shr")
	}
}

func (self *luaState) bnot(a luaValue) {
	if x, ok := convertToInteger(a); ok {
		self.stack.push(^x)
	} else if result, ok := callMetamethod(a, a, "__bnot", self); ok {
		self.stack.push(result)
	} else {
		panic("todo: __bnot!")
	}
}

/* helper */

func _castToInteger(a, b luaValue) (int64, int64, bool) {
	if x, ok := a.(int64); ok {
		if y, ok := b.(int64); ok {
			return x, y, true
		}
	}
	return 0, 0, false
}

func _convertToIntegers(a, b luaValue) (int64, int64, bool) {
	if x, ok := convertToInteger(a); ok {
		if y, ok := convertToInteger(b); ok {
			return x, y, true
		}
	}
	return 0, 0, false
}

func _convertToFloats(a, b luaValue) (float64, float64, bool) {
	if x, ok := convertToFloat(a); ok {
		if y, ok := convertToFloat(b); ok {
			return x, y, true
		}
	}
	return 0, 0, false
}
