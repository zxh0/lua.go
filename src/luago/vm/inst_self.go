package vm

import . "luago/api"

// R(A+1) := R(B); R(A) := R(B)[RK(C)]
func _self(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.Copy(b, a+1)

	vm.CheckStack(1)
	vm.GetRK(c)    // ~/rk[c]
	vm.GetTable(b) // ~/r[b][rk[c]]
	vm.Replace(a)  // ~
}
