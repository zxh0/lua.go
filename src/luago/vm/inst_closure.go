package vm

import . "luago/api"

// R(A) := closure(KPROTO[Bx])
func closure(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	vm.CheckStack(1)
	vm.LoadProto(bx) // ~/closure
	vm.Replace(a)    // ~
}
