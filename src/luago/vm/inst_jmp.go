package vm

import . "luago/lua"

// pc+=sBx; if (A) close all upvalues >= R(A - 1)
func jmp(i Instruction, vm VM) {
	a, sBx := i.AsBx()

	vm.AddPC(sBx)
	if a != 0 {
		panic("todo: jmp!")
	}
}
