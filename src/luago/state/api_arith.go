package state

import "math"
import . "luago/api"
import "luago/luanum"

// [-(2|1), +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_arith
func (self *luaState) Arith(op ArithOp) {
	var operand1, operand2 luaValue
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		operand2 = self.stack.pop()
	}
	operand1 = self.stack.pop()

	switch op {
	case LUA_OPADD:
		self.add(operand1, operand2)
	case LUA_OPSUB:
		self.sub(operand1, operand2)
	case LUA_OPMUL:
		self.mul(operand1, operand2)
	case LUA_OPMOD:
		self.mod(operand1, operand2)
	case LUA_OPPOW:
		self.pow(operand1, operand2)
	case LUA_OPDIV:
		self.div(operand1, operand2)
	case LUA_OPIDIV:
		self.idiv(operand1, operand2)
	case LUA_OPBAND:
		self.band(operand1, operand2)
	case LUA_OPBOR:
		self.bor(operand1, operand2)
	case LUA_OPBXOR:
		self.bxor(operand1, operand2)
	case LUA_OPSHL:
		self.shl(operand1, operand2)
	case LUA_OPSHR:
		self.shr(operand1, operand2)
	case LUA_OPUNM:
		self.unm(operand1)
	case LUA_OPBNOT:
		self.bnot(operand1)
	default:
		panic("invalid arith op!")
	}
}

/* integer or float */

func (self *luaState) add(a, b luaValue) {
	if x, y, ok := _castToInt64s(a, b); ok {
		self.stack.push(x + y)
	} else if x, y, ok := _convertToFloat64s(a, b); ok {
		self.stack.push(x + y)
	} else if result, ok := self.callMetaOp2(a, b, "__add"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __add")
	}
}

func (self *luaState) sub(a, b luaValue) {
	if x, y, ok := _castToInt64s(a, b); ok {
		self.stack.push(x - y)
	} else if x, y, ok := _convertToFloat64s(a, b); ok {
		self.stack.push(x - y)
	} else if result, ok := self.callMetaOp2(a, b, "__sub"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __sub")
	}
}

func (self *luaState) mul(a, b luaValue) {
	if x, y, ok := _castToInt64s(a, b); ok {
		self.stack.push(x * y)
	} else if x, y, ok := _convertToFloat64s(a, b); ok {
		self.stack.push(x * y)
	} else if result, ok := self.callMetaOp2(a, b, "__mul"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __mul")
	}
}

func (self *luaState) idiv(a, b luaValue) {
	if x, y, ok := _castToInt64s(a, b); ok {
		self.stack.push(luanum.IFloorDiv(x, y))
	} else if x, y, ok := _convertToFloat64s(a, b); ok {
		self.stack.push(luanum.FFloorDiv(x, y))
	} else if result, ok := self.callMetaOp2(a, b, "__idiv"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __idiv")
	}
}

func (self *luaState) mod(a, b luaValue) {
	if x, y, ok := _castToInt64s(a, b); ok {
		self.stack.push(luanum.IMod(x, y))
	} else if x, y, ok := _convertToFloat64s(a, b); ok {
		self.stack.push(luanum.FMod(x, y))
	} else if result, ok := self.callMetaOp2(a, b, "__mod"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __mod")
	}
}

func (self *luaState) unm(a luaValue) {
	if x, ok := a.(int64); ok {
		self.stack.push(-x)
	} else if x, ok := convertToNumber(a); ok {
		self.stack.push(-x)
	} else if result, ok := self.callMetaOp1(a, "__unm"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __unm!")
	}
}

/* float */

func (self *luaState) div(a, b luaValue) {
	if x, y, ok := _convertToFloat64s(a, b); ok {
		self.stack.push(x / y)
	} else if result, ok := self.callMetaOp2(a, b, "__div"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __div")
	}
}

func (self *luaState) pow(a, b luaValue) {
	if x, y, ok := _convertToFloat64s(a, b); ok {
		self.stack.push(math.Pow(x, y))
	} else if result, ok := self.callMetaOp2(a, b, "__pow"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __pow")
	}
}

/* bitwise */

func (self *luaState) band(a, b luaValue) {
	if x, y, ok := _convertToInt64s(a, b); ok {
		self.stack.push(x & y)
	} else if result, ok := self.callMetaOp2(a, b, "__band"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __band")
	}
}

func (self *luaState) bor(a, b luaValue) {
	if x, y, ok := _convertToInt64s(a, b); ok {
		self.stack.push(x | y)
	} else if result, ok := self.callMetaOp2(a, b, "__bor"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __bor")
	}
}

func (self *luaState) bxor(a, b luaValue) {
	if x, y, ok := _convertToInt64s(a, b); ok {
		self.stack.push(x ^ y)
	} else if result, ok := self.callMetaOp2(a, b, "__bxor"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __bxor")
	}
}

func (self *luaState) shl(a, b luaValue) {
	if x, y, ok := _convertToInt64s(a, b); ok {
		self.stack.push(luanum.ShiftLeft(x, y))
	} else if result, ok := self.callMetaOp2(a, b, "__shl"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __shl")
	}
}

func (self *luaState) shr(a, b luaValue) {
	if x, y, ok := _convertToInt64s(a, b); ok {
		self.stack.push(luanum.ShiftRight(x, y))
	} else if result, ok := self.callMetaOp2(a, b, "__shr"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __shr")
	}
}

func (self *luaState) bnot(a luaValue) {
	if x, ok := convertToInteger(a); ok {
		self.stack.push(^x)
	} else if result, ok := self.callMetaOp1(a, "__bnot"); ok {
		self.stack.push(result)
	} else {
		panic("todo: __bnot!")
	}
}

/* helper */

func _castToInt64s(a, b luaValue) (int64, int64, bool) {
	if x, ok := a.(int64); ok {
		if y, ok := b.(int64); ok {
			return x, y, true
		}
	}
	return 0, 0, false
}

func _convertToInt64s(a, b luaValue) (int64, int64, bool) {
	if x, ok := convertToInteger(a); ok {
		if y, ok := convertToInteger(b); ok {
			return x, y, true
		}
	}
	return 0, 0, false
}

func _convertToFloat64s(a, b luaValue) (float64, float64, bool) {
	if x, ok := convertToNumber(a); ok {
		if y, ok := convertToNumber(b); ok {
			return x, y, true
		}
	}
	return 0, 0, false
}
