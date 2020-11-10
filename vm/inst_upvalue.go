package vm

import . "github.com/zxh0/lua.go/api"

// R[A] := UpValue[B]
func getUpval(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(vm.UpvalueIndex(b), a)
}

// UpValue[B] := R[A]
func setUpval(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(a, vm.UpvalueIndex(b))
}

// R[A] := UpValue[B][K[C]:string]
func getTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	//vm.CheckStack(1)
	vm.GetConst(c)                  // ~/rk[c]
	vm.GetTable(vm.UpvalueIndex(b)) // ~/uv[b][k[c]]
	vm.Replace(a)                   // ~/
}

// UpValue[A][RK(B)] := RK(C)
func setTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	//vm.CheckStack(2)
	vm.GetRK(b)                     // ~/rk[b]
	vm.GetRK(c)                     // ~/rk[b]/rk[c]
	vm.SetTable(vm.UpvalueIndex(a)) // ~/
}

// close all upvalues >= R[A]
func closeUV(i Instruction, vm LuaVM) {
	//panic("TODO:close")
}
