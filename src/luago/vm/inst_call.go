package vm

import . "luago/api"

// R(A) := closure(KPROTO[Bx])
func closure(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	//vm.CheckStack(1)
	vm.LoadProto(bx) // ~/closure
	vm.Replace(a)    // ~
}

// R(A), R(A+1), ..., R(A+B-2) = vararg
func vararg(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b < 0 {
		panic("b < 0!")
	} else if b != 1 { // b==0 or b>1
		vm.LoadVararg(b - 1)
		_popResults(a, b, vm)
	}
}

// R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));
func tForCall(i Instruction, vm LuaVM) {
	a, _, c := i.ABC()
	a += 1

	_pushFuncAndArgs(a, 3, vm)
	vm.Call(2, c)
	_popResults(a+3, c+1, vm)
}

// return R(A)(R(A+1), ... ,R(A+B-1))
func tailCall(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	// todo: optimize tail call!
	c := 0
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

// R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
func call(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	// println(":::"+ vm.StackToString())
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

func _pushFuncAndArgs(a, b int, vm LuaVM) (nArgs int) {
	if b == 1 {
		nArgs = 0
		vm.CheckStack(1)
		vm.PushValue(a)
	} else if b > 1 {
		nArgs = b - 1
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
	} else {
		lastArgIdx := int(vm.ToInteger(-1))
		vm.Pop(1)

		nArgs = lastArgIdx - a
		nRegs := vm.MaxStackSize()

		if lastArgIdx <= nRegs {
			vm.CheckStack(nArgs + 1)
			for i := a; i <= lastArgIdx; i++ {
				vm.PushValue(i)
			}
		} else {
			vm.CheckStack(nRegs - a + 1)
			vm.SetTop(nRegs + nArgs + 1)
			for i := lastArgIdx; i >= a; i-- {
				vm.Copy(i, nRegs-a+1+i)
			}
		}
	}
	return
}

func _popResults(a, c int, vm LuaVM) {
	if c == 1 {
		// no results
	} else if c > 1 {
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else {
		nRegs := vm.MaxStackSize()
		nRets := vm.GetTop() - nRegs
		if nRets > 0 {
			//vm.Rotate(a, a-nRegs-1)
			for i := 0; i < nRets; i++ {
				vm.Copy(nRegs+1+i, a+i)
			}
			if nRegs-a+1 >= nRets {
				vm.Pop(nRets)
			} else {
				vm.Pop(nRegs - a + 1)
			}
		}
		vm.PushInteger(int64(a + nRets - 1))
	}
}

// return R(A), ... ,R(A+B-2)
func _return(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b == 1 {
		// no return values
	} else if b > 1 {
		// b-1 return values
		vm.CheckStack(b - 1)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
	} else {
		lastRetIdx := int(vm.ToInteger(-1))
		vm.Pop(1)

		nRets := lastRetIdx - a + 1
		nRegs := vm.MaxStackSize()

		if lastRetIdx <= nRegs {
			vm.CheckStack(nRets)
			for i := a; i <= lastRetIdx; i++ {
				vm.PushValue(i)
			}
		} else {
			vm.CheckStack(nRegs - a + 1)
			vm.SetTop(nRegs + nRets)
			for i := lastRetIdx; i >= a; i-- {
				vm.Copy(i, nRegs-a+1+i)
			}
		}
	}
}
