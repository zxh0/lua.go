package vm

import . "luago/api"
import "luago/number"

// R(A) := {} (size = B,C)
func newTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	//vm.CheckStack(1)
	vm.CreateTable(number.Fb2int(b), number.Fb2int(c)) // ~/{}
	vm.Replace(a)                                      // ~
}

// R(A) := R(B)[RK(C)]
func getTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	//vm.CheckStack(1)
	vm.GetRK(c)    // ~/rk[c]
	vm.GetTable(b) // ~/r[b][rk[c]]
	vm.Replace(a)  // ~
}

// R(A)[RK(B)] := RK(C)
func setTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	//vm.CheckStack(2)
	vm.GetRK(b)    // ~/rk[b]
	vm.GetRK(c)    // ~/rk[b]/rk[c]
	vm.SetTable(a) // ~
}

// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
func setList(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	if c == 0 {
		vm.AddPC(1)
		c = Instruction(vm.Instruction()).Ax()
	}

	bIsZero := b == 0
	if bIsZero {
		b = int(vm.ToInteger(-1)) - a - 1
		vm.Pop(1)
	}

	vm.CheckStack(1)
	tableIdx := int64((c - 1) * LFIELDS_PER_FLUSH)
	for j := 1; j <= b; j++ {
		tableIdx++
		vm.PushValue(a + j)  // ~/r[a+j]
		vm.SetI(a, tableIdx) // ~
	}

	if bIsZero {
		n := vm.GetTop() - vm.MaxStackSize()
		for j := 1; j <= n; j++ {
			tableIdx++
			vm.PushValue(vm.MaxStackSize() + j)
			vm.SetI(a, tableIdx)
		}

		// clear stack
		vm.SetTop(vm.MaxStackSize())
	}
}
