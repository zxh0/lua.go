package vm

import . "github.com/zxh0/lua.go/api"

/* arith */

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

// R[A] := R[B] op R[C]
func _binaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, c := i.ABC()
	a += 1

	//vm.CheckStack(2)
	vm.GetReg(b)  // ~/r[b]
	vm.GetReg(c)  // ~/r[b]/r[c]
	vm.Arith(op)  // ~/result
	vm.Replace(a) // ~
}

// R[A] := op R[B]
func _unaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	//vm.CheckStack(1)
	vm.PushValue(b) // ~/r[b]
	vm.Arith(op)    // ~/result
	vm.Replace(a)   // ~
}

/* compare */

func eq(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPEQ) } // ==
func lt(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLT) } // <
func le(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLE) } // <=

// if ((R[A] op R[B]) ~= k) then pc++
func _compare(i Instruction, vm LuaVM, op CompareOp) {
	a, b, _, k := i.ABCk()

	//vm.CheckStack(2)
	vm.GetReg(a) // ~/r[a]
	vm.GetReg(b) // ~/r[a]/r[b]
	if vm.Compare(-2, -1, op) != (k != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2) // ~
}

/* logical */

// R[A] := not R[B]
func not(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	//vm.CheckStack(1)
	vm.PushBoolean(!vm.ToBoolean(b)) // ~/!r[b]
	vm.Replace(a)                    // ~
}

// if (not R[A] == k) then pc++
func test(i Instruction, vm LuaVM) {
	a, _, _, k := i.ABCk()
	a += 1

	if vm.ToBoolean(a) != (k != 0) {
		vm.AddPC(1)
	}
}

// if (not R[B] == k) then pc++ else R[A] := R[B]
func testSet(i Instruction, vm LuaVM) {
	a, b, _, k := i.ABCk()
	a += 1
	b += 1

	if vm.ToBoolean(b) == (k != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPC(1)
	}
}

/* len & concat */

// R[A] := length of R[B]
func length(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	//vm.CheckStack(1)
	vm.Len(b)     // ~/#r[b]
	vm.Replace(a) // ~
}

// R(A) := R(B).. ... ..R(C)
// R[A] := R[A].. ... ..R[A + B - 1]
func concat(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	vm.CheckStack(b)
	for i := 0; i < b; i++ {
		vm.PushValue(a + i) // ~/r[a]/.../r[a + b - 1]
	}
	vm.Concat(b)  // ~/result
	vm.Replace(a) // ~
}
