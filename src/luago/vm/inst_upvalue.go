package vm

import . "luago/api"

// R(A) := UpValue[B]
func getUpval(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	vm.CheckStack(1)
	vm.GetUpvalue2(b) // ~/uv[b]
	vm.Replace(a)     // ~
}

// UpValue[B] := R(A)
func setUpval(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	vm.CheckStack(1)
	vm.PushValue(a)   // ~/r[a]
	vm.SetUpvalue2(b) // ~
}

// R(A) := UpValue[B][RK(C)]
func getTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.CheckStack(2)
	vm.GetUpvalue2(b) // ~/uv[b]
	vm.GetRK(c)       // ~/uv[b]/rk[c]
	vm.GetTable(-2)   // ~/uv[b]/uv[b][rk[c]]
	vm.Replace(a)     // ~/uv[b]
	vm.Pop(1)         // ~
}

// UpValue[A][RK(B)] := RK(C)
func setTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()

	vm.CheckStack(3)
	vm.GetUpvalue2(a) // ~/uv[a]
	vm.GetRK(b)       // ~/uv[a]/rk[b]
	vm.GetRK(c)       // ~/uv[a]/rk[b]/rk[c]
	vm.SetTable(-3)   // ~/uv[a]
	vm.Pop(1)         // ~
}
