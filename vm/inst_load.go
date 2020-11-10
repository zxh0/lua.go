package vm

import . "github.com/zxh0/lua.go/api"

// R[A], R[A+1], ..., R[A+B] := nil
func loadNil(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	//vm.CheckStack(1)
	vm.PushNil() // ~/nil
	for i := a; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1) // ~
}

// R[A] := false; pc++
func lFalseSkip(i Instruction, vm LuaVM) {
	a, _, _ := i.ABC()
	a += 1

	vm.PushBoolean(false)
	vm.Replace(a)
	vm.AddPC(1)
}

// R[A] := false
func loadFalse(i Instruction, vm LuaVM) {
	a, _, _ := i.ABC()
	a += 1

	vm.PushBoolean(false)
	vm.Replace(a)
}

// R[A] := true
func loadTrue(i Instruction, vm LuaVM) {
	a, _, _ := i.ABC()
	a += 1

	vm.PushBoolean(true)
	vm.Replace(a)
}

// R[A] := sBx
func loadI(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	vm.PushInteger(int64(sBx))
	vm.Replace(a)
}

// R[A] := (lua_Number)sBx
func loadF(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	vm.PushNumber(float64(sBx))
	vm.Replace(a)
}

// R[A] := K[Bx]
func loadK(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	//vm.CheckStack(1)
	vm.GetConst(bx) // ~/k[bx]
	vm.Replace(a)   // ~
}

// R[A] := K[extra arg]
func loadKx(i Instruction, vm LuaVM) {
	a, _ := i.ABx()
	a += 1
	ax := Instruction(vm.Fetch()).Ax()

	//vm.CheckStack(1)
	vm.GetConst(ax) // ~/k[ax]
	vm.Replace(a)   // ~
}
