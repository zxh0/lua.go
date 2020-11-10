package api

type LuaVM interface {
	LuaState
	AddPC(n int)         // pc += n
	Fetch() uint32       // code[pc++]
	RegisterCount() int  // proto.MaxStackSize
	GetConst(idx int)    // push(const[idx])
	GetRK(rk int)        // rk > 0xFF ? GetConst(rk & 0xFF) : PushValue(rk + 1)
	GetRK2(x, k int)     // k ? GetConst(x) : PushValue(x + 1)
	GetReg(r int)        // PushValue(r + 1)
	LoadProto(idx int)   // push(proto[idx] as LuaFunction)
	LoadVararg(n int)    // push(vararg[0], ..., vararg[n-1])
	CloseUpvalues(a int) // close all upvalues >= R(A - 1)

	GetConstStr(idx int) string
	GetConstInt(idx int) int64
}
