package vm

import . "github.com/zxh0/lua.go/api"

/* opK */

func addK(i Instruction, vm LuaVM)  { _binaryArithK(i, vm, LUA_OPADD) }  // +
func subK(i Instruction, vm LuaVM)  { _binaryArithK(i, vm, LUA_OPSUB) }  // -
func mulK(i Instruction, vm LuaVM)  { _binaryArithK(i, vm, LUA_OPMUL) }  // *
func modK(i Instruction, vm LuaVM)  { _binaryArithK(i, vm, LUA_OPMOD) }  // %
func powK(i Instruction, vm LuaVM)  { _binaryArithK(i, vm, LUA_OPPOW) }  // ^
func divK(i Instruction, vm LuaVM)  { _binaryArithK(i, vm, LUA_OPDIV) }  // /
func idivK(i Instruction, vm LuaVM) { _binaryArithK(i, vm, LUA_OPIDIV) } // //
func bandK(i Instruction, vm LuaVM) { _binaryArithK(i, vm, LUA_OPBAND) } // &
func borK(i Instruction, vm LuaVM)  { _binaryArithK(i, vm, LUA_OPBOR) }  // |
func bxorK(i Instruction, vm LuaVM) { _binaryArithK(i, vm, LUA_OPBXOR) } // ~

// OP_ADDK, /*	A B C	R[A] := R[B] + K[C]				*/
// OP_SUBK, /*	A B C	R[A] := R[B] - K[C]				*/
// OP_MULK, /*	A B C	R[A] := R[B] * K[C]				*/
// OP_MODK, /*	A B C	R[A] := R[B] % K[C]				*/
// OP_POWK, /*	A B C	R[A] := R[B] ^ K[C]				*/
// OP_DIVK, /*	A B C	R[A] := R[B] / K[C]				*/
// OP_IDIVK,/*	A B C	R[A] := R[B] // K[C]			*/
// OP_BANDK,/*	A B C	R[A] := R[B] & K[C]:integer		*/
// OP_BORK, /*	A B C	R[A] := R[B] | K[C]:integer		*/
// OP_BXORK,/*	A B C	R[A] := R[B] ~ K[C]:integer		*/

func _binaryArithK(i Instruction, vm LuaVM, op ArithOp) {
	a, b, c := i.ABC()
	a += 1

	//vm.CheckStack(2)
	vm.GetReg(b)   // ~/r[b]
	vm.GetConst(c) // ~/r[b]/k[c]
	vm.Arith(op)   // ~/result
	vm.Replace(a)  // ~

	// TODO
	vm.AddPC(1)
}

/* opI */

// R[A] := R[B] + sC
func addI(i Instruction, vm LuaVM) {
	a, b, sC := i.ABsC()
	a += 1

	//vm.CheckStack(2)
	vm.GetReg(b)              // ~/r[b]
	vm.PushInteger(int64(sC)) // ~/r[b]/sC
	vm.Arith(LUA_OPADD)       // ~/result
	vm.Replace(a)             // ~

	// TODO
	vm.AddPC(1)
}

// R[A] := R[B] >> sC
func shrI(i Instruction, vm LuaVM) {
	a, b, sC := i.ABsC()
	a += 1

	//vm.CheckStack(2)
	vm.GetReg(b)              // ~/r[b]
	vm.PushInteger(int64(sC)) // ~/r[b]/sC
	vm.Arith(LUA_OPSHR)       // ~/result
	vm.Replace(a)             // ~

	// TODO
	vm.AddPC(1)
}

// R[A] := sC << R[B]
func shlI(i Instruction, vm LuaVM) {
	a, b, sC := i.ABsC()
	a += 1

	//vm.CheckStack(2)
	vm.PushInteger(int64(sC)) // ~/sC
	vm.GetReg(b)              // ~/sC/r[b]
	vm.Arith(LUA_OPSHL)       // ~/result
	vm.Replace(a)             // ~

	// TODO
	vm.AddPC(1)
}

/* compare */

// if ((R[A] == K[B]) ~= k) then pc++
func eqK(i Instruction, vm LuaVM) {
	a, b, _, k := i.ABCk()

	//vm.CheckStack(2)
	vm.GetReg(a)   // ~/r[a]
	vm.GetConst(b) // ~/r[a]/k[b]
	if vm.Compare(-2, -1, LUA_OPEQ) != (k != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2) // ~
}

// if ((R[A] == sB) ~= k) then pc++
func eqI(i Instruction, vm LuaVM) {
	a, sB, k := i.AsBk()
	a += 1

	if vm.ToInteger(a) == int64(sB) != (k != 0) {
		vm.AddPC(1)
	}
}

// if ((R[A] < sB) ~= k) then pc++
func ltI(i Instruction, vm LuaVM) {
	a, sB, k := i.AsBk()
	a += 1

	if vm.ToInteger(a) < int64(sB) != (k != 0) {
		vm.AddPC(1)
	}
}

// if ((R[A] <= sB) ~= k) then pc++
func leI(i Instruction, vm LuaVM) {
	a, sB, k := i.AsBk()
	a += 1

	if vm.ToInteger(a) <= int64(sB) != (k != 0) {
		vm.AddPC(1)
	}
}

// if ((R[A] > sB) ~= k) then pc++
func gtI(i Instruction, vm LuaVM) {
	a, sB, k := i.AsBk()
	a += 1

	if vm.ToInteger(a) > int64(sB) != (k != 0) {
		vm.AddPC(1)
	}
}

// if ((R[A] >= sB) ~= k) then pc++
func geI(i Instruction, vm LuaVM) {
	a, sB, k := i.AsBk()
	a += 1

	if vm.ToInteger(a) >= int64(sB) != (k != 0) {
		vm.AddPC(1)
	}
}
