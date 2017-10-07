package vm

import . "luago/api"

// R(A) := R(B)
func move(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(b, a)
}

// pc+=sBx; if (A) close all upvalues >= R(A - 1)
func jmp(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()

	vm.AddPC(sBx)
	if a != 0 {
		panic("todo: jmp!")
	}
}

// R(A+1) := R(B); R(A) := R(B)[RK(C)]
func _self(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.Copy(b, a+1)

	//vm.CheckStack(1)
	vm.GetRK(c)    // ~/rk[c]
	vm.GetTable(b) // ~/r[b][rk[c]]
	vm.Replace(a)  // ~
}
