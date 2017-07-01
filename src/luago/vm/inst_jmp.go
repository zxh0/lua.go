package vm

import . "luago/api"

// pc+=sBx; if (A) close all upvalues >= R(A - 1)
func jmp(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()

	vm.AddPC(sBx)
	if a != 0 {
		panic("todo: jmp!")
	}
}
