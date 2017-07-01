package vm

import . "luago/api"

func eq(i Instruction, vm VM) { _compare(i, vm, LUA_OPEQ) } // ==
func lt(i Instruction, vm VM) { _compare(i, vm, LUA_OPLT) } // <
func le(i Instruction, vm VM) { _compare(i, vm, LUA_OPLE) } // <=

// if ((RK(B) op RK(C)) ~= A) then pc++
func _compare(i Instruction, vm VM, op LuaCompareOp) {
	a, b, c := i.ABC()

	vm.CheckStack(2)
	vm.GetRK(b) // ~/rk[b]
	vm.GetRK(c) // ~/rk[b]/rk[c]
	if vm.Compare(-2, -1, op) != (a != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2) // ~
}

// if not (R(A) <=> C) then pc++
func test(i Instruction, vm VM) {
	a, _, c := i.ABC()
	a += 1

	if vm.ToBoolean(a) != (c != 0) {
		vm.AddPC(1)
	}
}

// if (R(B) <=> C) then R(A) := R(B) else pc++
func testSet(i Instruction, vm VM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	if vm.ToBoolean(b) == (c != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPC(1)
	}
}
