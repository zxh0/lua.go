package vm

import . "github.com/zxh0/lua.go/api"

// R(A)-=R(A+2); pc+=sBx
// <check values and prepare counters>;
// if not to run then pc+=Bx+1;
func forPrep(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	//vm.CheckStack(2)
	if vm.Type(a) == LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a))
		vm.Replace(a)
	}
	if vm.Type(a+1) == LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 1))
		vm.Replace(a + 1)
	}
	if vm.Type(a+2) == LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 2))
		vm.Replace(a + 2)
	}

	vm.PushValue(a)     // ~/r[a]
	vm.PushValue(a + 2) // ~/r[a]/r[a+2]
	vm.Arith(LUA_OPSUB) // ~/r[a]-r[a+2]
	vm.Replace(a)       // ~
	vm.AddPC(bx)
}

// R(A)+=R(A+2);
// if R(A) <?= R(A+1) then {
//   pc+=sBx; R(A+3)=R(A)
// }
// update counters; if loop continues then pc-=Bx;
func forLoop(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	//vm.CheckStack(2)
	// R(A)+=R(A+2);
	vm.PushValue(a + 2) // ~/r[a+2]
	vm.PushValue(a)     // ~/r[a+2]/r[a]
	vm.Arith(LUA_OPADD) // ~/r[a]+r[a+2]
	vm.Replace(a)       // ~

	isPositiveStep := vm.ToNumber(a+2) >= 0
	if isPositiveStep && vm.Compare(a, a+1, LUA_OPLE) ||
		!isPositiveStep && vm.Compare(a+1, a, LUA_OPLE) {

		// pc+=sBx; R(A+3)=R(A)
		vm.AddPC(-bx)
		vm.Copy(a, a+3)
	}
}

// create upvalue for R[A + 3]; pc+=Bx
func tForPrep(i Instruction, vm LuaVM) {
	_, bx := i.ABx()
	// TODO
	vm.AddPC(bx)
}

// if R(A+1) ~= nil then {
//   R(A)=R(A+1); pc += sBx
// }
// if R[A+2] ~= nil then { R[A]=R[A+2]; pc -= Bx } ???
func tForLoop(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1
	if !vm.IsNil(a + 4) {
		vm.Copy(a+4, a+2)
		vm.AddPC(-bx)
	}
}
