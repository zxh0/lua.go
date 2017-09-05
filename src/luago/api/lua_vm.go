package api

type LuaVM interface {
	LuaState
	AddPC(n int)         // pc += n
	MaxStackSize() int   // proto.MaxStackSize
	GetConst(idx int)    // push(const[idx])
	GetRK(rk int)        // rk > 0xFF ? GetConst(rk & 0xFF) : PushValue(rk + 1)
	GetUpvalue2(idx int) // push(upvalue[idx])
	SetUpvalue2(idx int) // upvalue[idx] = pop()
	LoadProto(idx int)   // push(proto[idx] as LuaFunction)
	LoadVararg(n int)    // push(vararg[0], ..., vararg[n-1])
}
