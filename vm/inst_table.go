package vm

import . "github.com/zxh0/lua.go/api"

/* number of list items to accumulate before a SETLIST instruction */
const LFIELDS_PER_FLUSH = 50

// R(A) := {} (size = B,C)
func newTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	//vm.CheckStack(1)
	vm.CreateTable(Fb2int(b), Fb2int(c)) // ~/{}
	vm.Replace(a)                        // ~

	// TODO
	vm.AddPC(1)
}

// R[A] := R[B][R[C]]
func getTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	//vm.CheckStack(1)
	vm.GetReg(c)   // ~/r[c]
	vm.GetTable(b) // ~/r[b][r[c]]
	vm.Replace(a)  // ~
}

// R[A][R[B]] := RK(C)
func setTable(i Instruction, vm LuaVM) {
	a, b, c, k := i.ABCk()
	a += 1

	//vm.CheckStack(2)
	vm.GetReg(b)    // ~/r[b]
	vm.GetRK2(c, k) // ~/r[b]/rk[c]
	vm.SetTable(a)  // ~
}

// R[A] := R[B][K[C]:string]
func getField(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	//vm.CheckStack(1)
	f := vm.GetConstStr(c)
	vm.GetField(b, f) // ~/r[b][k[c]]
	vm.Replace(a)     // ~
}

// R[A][K[B]:string] := RK(C)
func setField(i Instruction, vm LuaVM) {
	a, b, c, k := i.ABCk()
	a += 1

	//vm.CheckStack(2)
	f := vm.GetConstStr(b)
	vm.GetRK2(c, k)   // ~rk[c]
	vm.SetField(a, f) // ~
}

// R[A] := R[B][C]
func getI(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	//vm.CheckStack(1)
	vm.GetI(b, int64(c)) // ~/r[b][c]
	vm.Replace(a)        // ~
}

// R[A][B] := RK(C)
func setI(i Instruction, vm LuaVM) {
	a, b, c, k := i.ABCk()
	a += 1

	//vm.CheckStack(2)
	vm.GetRK2(c, k)      // ~rk[c]
	vm.SetI(a, int64(b)) // ~
}

// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
func setList(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	if c > 0 {
		c = c - 1
	} else {
		c = Instruction(vm.Fetch()).Ax()
	}

	bIsZero := b == 0
	if bIsZero {
		b = int(vm.ToInteger(-1)) - a - 1
		vm.Pop(1)
	}

	vm.CheckStack(1)
	idx := int64(c * LFIELDS_PER_FLUSH)
	for j := 1; j <= b; j++ {
		idx++
		vm.PushValue(a + j) // ~/r[a+j]
		vm.SetI(a, idx)     // ~
	}

	if bIsZero {
		for j := vm.RegisterCount() + 1; j <= vm.GetTop(); j++ {
			idx++
			vm.PushValue(j)
			vm.SetI(a, idx)
		}

		// clear stack
		vm.SetTop(vm.RegisterCount())
	}
}
