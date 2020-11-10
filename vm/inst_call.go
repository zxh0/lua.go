package vm

import . "github.com/zxh0/lua.go/api"

// R(A) := closure(KPROTO[Bx])
func closure(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	//vm.CheckStack(1)
	vm.LoadProto(bx) // ~/closure
	vm.Replace(a)    // ~
}

func varargPrep(i Instruction, vm LuaVM) {
	// TODO
}

// R(A), R(A+1), ..., R(A+B-2) = vararg
func vararg(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b != 1 { // b==0 or b>1
		vm.LoadVararg(b - 1)
		_popResults(a, b, vm)
	}
}

// R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));
// R[A+4], ... ,R[A+3+C] := R[A](R[A+1], R[A+2]);
/* 'ra' has the iterator function, 'ra + 1' has the state,
   'ra + 2' has the control variable, and 'ra + 3' has the
   to-be-closed variable. The call will use the stack after
   these values (starting at 'ra + 4')
*/
func tForCall(i Instruction, vm LuaVM) {
	a, _, c := i.ABC()
	a += 1

	_pushFuncAndArgs(a, 3, vm)
	vm.Call(2, c)
	_popResults(a+4, c+1, vm)
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

// R[A], ... ,R[A+C-2] := R[A](R[A+1], ... ,R[A+B-1])
func call(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	// println(":::"+ vm.StackToString())
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

func _pushFuncAndArgs(a, b int, vm LuaVM) (nArgs int) {
	if b >= 1 {
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else {
		_fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

func _fixStack(a int, vm LuaVM) {
	x := int(vm.ToInteger(-1))
	vm.Pop(1)

	vm.CheckStack(x - a)
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	vm.Rotate(vm.RegisterCount()+1, x-a)
}

func _popResults(a, c int, vm LuaVM) {
	if c == 1 {
		// no results
	} else if c > 1 {
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else {
		// leave results on stack
		vm.CheckStack(1)
		vm.PushInteger(int64(a))
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
		_fixStack(a, vm)
	}
}
