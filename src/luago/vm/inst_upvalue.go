package vm

import . "luago/api"

// R(A) := UpValue[B]
func getUpval(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(LuaUpvalueIndex(b), a)
}

// UpValue[B] := R(A)
func setUpval(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(a, LuaUpvalueIndex(b))
}

// R(A) := UpValue[B][RK(C)]
func getTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	//vm.CheckStack(1)
	vm.GetRK(c)                     // ~/rk[c]
	vm.GetTable(LuaUpvalueIndex(b)) // ~/uv[b][rk[c]]
	vm.Replace(a)                   // ~/
}

// UpValue[A][RK(B)] := RK(C)
func setTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	//vm.CheckStack(2)
	vm.GetRK(b)                     // ~/rk[b]
	vm.GetRK(c)                     // ~/rk[b]/rk[c]
	vm.SetTable(LuaUpvalueIndex(a)) // ~/
}
