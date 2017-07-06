package vm

import . "luago/api"

func add(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPADD) }  // +
func sub(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSUB) }  // -
func mul(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMUL) }  // *
func mod(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMOD) }  // %
func pow(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPPOW) }  // ^
func div(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPDIV) }  // /
func idiv(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPIDIV) } // //
func band(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBAND) } // &
func bor(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPBOR) }  // |
func bxor(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBXOR) } // ~
func shl(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHL) }  // <<
func shr(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHR) }  // >>
func unm(i Instruction, vm LuaVM)  { _unaryArith(i, vm, LUA_OPUNM) }   // -
func bnot(i Instruction, vm LuaVM) { _unaryArith(i, vm, LUA_OPBNOT) }  // ~

// R(A) := RK(B) op RK(C)
func _binaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, c := i.ABC()
	a += 1

	vm.CheckStack(2)
	vm.GetRK(b)   // ~/rk[b]
	vm.GetRK(c)   // ~/rk[b]/rk[c]
	vm.Arith(op)  // ~/result
	vm.Replace(a) // ~
}

// R(A) := op R(B)
func _unaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.CheckStack(1)
	vm.PushValue(b) // ~/r[b]
	vm.Arith(op)    // ~/result
	vm.Replace(a)   // ~
}

// R(A) := not R(B)
func not(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.CheckStack(1)
	vm.PushBoolean(!vm.ToBoolean(b)) // ~/!r[b]
	vm.Replace(a)                    // ~
}

// R(A) := length of R(B)
func _len(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.CheckStack(1)
	vm.Len(b)     // ~/#r[b]
	vm.Replace(a) // ~
}

// R(A) := R(B).. ... ..R(C)
func concat(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	c += 1

	n := c - b + 1
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		vm.PushValue(i) // ~/r[b]/.../r[c]
	}
	vm.Concat(n)  // ~/result
	vm.Replace(a) // ~
}
