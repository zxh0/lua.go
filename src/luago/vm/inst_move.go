package vm

import . "luago/lua"

// R(A) := R(B)
func move(i Instruction, vm VM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(b, a)
}
