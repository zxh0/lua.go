package vm

import . "luago/lua"

// R(A), R(A+1), ..., R(A+B-2) = vararg
func vararg(i Instruction, vm VM) {
	a, b, _ := i.ABC()
	a += 1

	if b < 0 {
		panic("b < 0!")
	} else if b != 1 { // b==0 or b>1
		vm.LoadVararg(b - 1)
		_moveResults(a, b, vm)
	}
}

// R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));
func tForCall(i Instruction, vm VM) {
	a, _, c := i.ABC()
	a += 1

	_pushFuncAndArgs(a, 3, vm)
	vm.Call(2, c)
	_moveResults(a+3, c+1, vm)
}

// R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
func call(i Instruction, vm VM) {
	a, b, c := i.ABC()
	a += 1

	// println(":::"+ vm.StackToString())
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_moveResults(a, c, vm)
}

// return R(A)(R(A+1), ... ,R(A+B-1))
func tailCall(i Instruction, vm VM) {
	a, b, _ := i.ABC()
	a += 1

	// todo: optimize tail call!
	c := 0
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_moveResults(a, c, vm)
}

func _pushFuncAndArgs(a, b int, vm VM) (nArgs int) {
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
		top := vm.GetTop()
		btm := vm.GetBottom()

		if lastArgIdx <= btm {
			vm.CheckStack(lastArgIdx - a + 1)
			for i := a; i <= lastArgIdx; i++ {
				vm.PushValue(i)
			}
		} else {
			vm.CheckStack(btm - a + 1)
			for i := a; i <= btm; i++ {
				vm.PushValue(i)
			}
			if top > btm {
				vm.Rotate(btm+1, btm-top)
			}
		}
	}
	return
}

func _moveResults(a, c int, vm VM) {
	if c == 1 {
		// no results
	} else if c > 1 {
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else {
		top := vm.GetTop()
		btm := vm.GetBottom()
		nRets := top - btm
		if nRets > 0 {
			vm.Rotate(a, a-btm-1)
			if btm+1-a >= nRets {
				vm.Pop(nRets)
			} else {
				vm.Pop(btm + 1 - a)
			}
		}
		vm.PushInteger(int64(a + nRets - 1))
	}
}

// return R(A), ... ,R(A+B-2)
func _return(i Instruction, vm VM) {
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
		// todo: panic("todo: return & b == 0!")
		lastRetIdx := int(vm.ToInteger(-1))
		vm.Pop(1)

		top := vm.GetTop()
		btm := vm.GetBottom()

		if lastRetIdx <= btm {
			vm.CheckStack(btm - lastRetIdx + 1)
			for i := a; i <= lastRetIdx; i++ {
				vm.PushValue(i)
			}
		} else {
			vm.CheckStack(btm - a + 1)
			for i := a; i <= btm; i++ {
				vm.PushValue(i)
			}
			if top > btm {
				vm.Rotate(btm+1, btm-top)
			}
		}
	}
}
