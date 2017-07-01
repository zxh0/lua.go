package vm

import . "luago/api"

// R(A) := {} (size = B,C)
func newTable(i Instruction, vm VM) {
	a, b, c := i.ABC()
	a += 1

	vm.CheckStack(1)
	vm.CreateTable(FB2INT(b), FB2INT(c)) // ~/{}
	vm.Replace(a)        // ~
}

// R(A) := R(B)[RK(C)]
func getTable(i Instruction, vm VM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.CheckStack(1)
	vm.GetRK(c)    // ~/rk[c]
	vm.GetTable(b) // ~/r[b][rk[c]]
	vm.Replace(a)  // ~
}

// R(A)[RK(B)] := RK(C)
func setTable(i Instruction, vm VM) {
	a, b, c := i.ABC()
	a += 1

	vm.CheckStack(2)
	vm.GetRK(b)    // ~/rk[b]
	vm.GetRK(c)    // ~/rk[b]/rk[c]
	vm.SetTable(a) // ~
}

// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
func setList(i Instruction, vm VM) {
	a, b, c := i.ABC()
	a += 1

	if c == 0 {
		panic("todo: c == 0") // #list > 25550
	}

	bIsZero := b == 0
	if bIsZero {
		lastArgIdx := int(vm.ToInteger(-1))
		vm.Pop(1)

		b = lastArgIdx - a
	}

	vm.CheckStack(1)
	for j := 1; j <= b; j++ {
		n := (c-1)*LFIELDS_PER_FLUSH + j
		vm.PushValue(a + j)  // ~/r[a+j]
		vm.SetI(a, int64(n)) // ~
	}

	// clear stack
	if bIsZero {
		vm.Pop(vm.GetTop() - vm.GetBottom())
	}
}
