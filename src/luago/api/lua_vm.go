package api

type LuaVM interface {
	LuaState
	AddPC(n int)          // pc += n
	GetBottom() int       // stack.bottom
	LoadProto(index int)  // push(proto[index] as LuaFunction)
	LoadVararg(n int)     // push(vararg[0], ..., vararg[n-1])
	GetRK(rk int)         // rk > 0xFF ? GetConst(rk & 0xFF) : vm.PushValue(rk + 1)
	GetConst(index int)   // push(const[index])
	GetUpvalue(index int) // push(upvalue[index])
	SetUpvalue(index int) // upvalue[index] = pop()
}
