package api

type LuaType int
type ArithOp int
type CompareOp int
type ThreadStatus int

// type LuaBoolean bool
// type LuaInteger int64
// type LuaNumber float64
type GoFunction func(LuaState) int
type UserData interface{}

type LuaState interface {
	BasicAPI
	DebugAPI
	AuxLib
	String() string // debug
}

type BasicAPI interface {
	/* state manipulation */
	Close()                               //
	AtPanic(panicf GoFunction) GoFunction // set panic function
	Version() float64                     // get version number
	/* basic stack manipulation */
	AbsIndex(idx int) int     // abs(idx)
	GetTop() int              // stack.top
	SetTop(idx int)           // stack.top = idx
	Pop(n int)                // pop(n)
	PushValue(idx int)        // push(r[idx])
	Rotate(idx, n int)        // r[idx, -1] >> n
	Insert(idx int)           // r[idx, -1] >> 1
	Remove(idx int)           // remove(r[idx])
	Replace(idx int)          // r[idx] = pop()
	Copy(fromIdx, toIdx int)  // r[toIdx] = r[fromidx]
	CheckStack(n int) bool    //
	XMove(to LuaState, n int) // to.push(pop(n))
	/* access functions (stack -> C) */
	Type(idx int) LuaType              // r[idx].type
	TypeName(tp LuaType) string        // r[idx].type.name
	IsNumber(idx int) bool             // r[idx] ~= LuaNumber
	IsString(idx int) bool             // r[idx] ~= LuaString
	IsGoFunction(idx int) bool         // r[idx].type == LUA_TLCL || LUA_TGCL
	IsInteger(idx int) bool            // r[idx].type == LUA_TINTEGER
	IsUserData(idx int) bool           // r[idx].type == LUA_TUSERDATA
	IsNone(idx int) bool               // r[idx].type == LUA_TNONE
	IsNil(idx int) bool                // r[idx].type == LUA_TNIL
	IsNoneOrNil(idx int) bool          // r[idx].type == LUA_TNONE || LUA_TNIL
	IsBoolean(idx int) bool            // r[idx].type == LUA_TBOOLEAN
	IsTable(idx int) bool              // r[idx].type == LUA_TTABLE
	IsFunction(idx int) bool           // r[idx].type == LUA_TFUNCTION
	IsThread(idx int) bool             // r[idx].type == LUA_TTHREAD
	ToNumberX(idx int) (float64, bool) // r[idx] as LuaNumber
	ToIntegerX(idx int) (int64, bool)  // r[idx] as LuaInteger
	ToBoolean(idx int) bool            // r[idx] as bool
	ToString(idx int) (string, bool)   // r[idx] as string
	ToGoFunction(idx int) GoFunction   // r[idx] as GoFunction
	ToUserData(idx int) UserData       // r[idx] as UserData
	ToThread(idx int) LuaState         // r[idx] as LuaThread
	ToPointer(idx int) interface{}     // r[idx] as interface{}
	ToInteger(idx int) int64           // r[idx] as LuaInteger
	ToNumber(idx int) float64          // r[idx] as LuaNumber
	RawLen(idx int) uint               // len(r[idx])
	/* Comparison and arithmetic functions */
	Arith(op ArithOp)                          // b=pop(); a=pop(); push(a op b)
	Compare(idx1, idx2 int, op CompareOp) bool // r[idx1] op r[idx2]
	RawEqual(idx1, idx2 int) bool              // r[idx1] == r[idx2]
	/* push functions (C -> stack) */
	PushBoolean(b bool)                 // push(b)
	PushGoClosure(fn GoFunction, n int) // push(f)
	PushGoFunction(f GoFunction)        // push(f)
	PushFString(fmt string)             // todo
	PushGlobalTable()                   // push(global)
	PushInteger(n int64)                // push(n)
	PushUserData(d UserData)            // push(d)
	PushNil()                           // push(nil)
	PushNumber(n float64)               // push(n)
	PushString(s string)                // push(s)
	PushThread(ls LuaState) bool        // push(ls)
	PushVFString()                      // todo
	/* get functions (Lua -> stack) */
	GetGlobal(name string) LuaType       // push(global[name])
	GetTable(idx int) LuaType            // push(r[idx][pop()])
	GetField(idx int, k string) LuaType  // push(r[idx][k])
	GetI(idx int, i int64) LuaType       // push(r[idx][i])
	RawGet(idx int) LuaType              // push(r[idx][pop()])
	RawGetI(idx int, n int64) LuaType    // push(r[idx][i])
	RawGetP(idx int, p UserData) LuaType // push(r[idx][p])
	CreateTable(nArr, nRec int)          // push({})
	GetMetatable(idx int) bool           // push(r[idx].metatable)?
	GetUserValue(idx int) LuaType        // push(r[idx].userValue)
	NewTable()                           // push({})
	/* set functions (stack -> Lua) */
	Register(name string, f GoFunction) // global[name] = f
	SetGlobal(name string)              // global[name] = pop()
	SetTable(idx int)                   // v=pop(); k=pop(); r[idx][k] = v
	SetField(idx int, k string)         // r[idx][k] = pop()
	SetI(idx int, n int64)              // r[idx][n] = pop()
	SetMetatable(idx int)               // r[idx].metatable = pop()
	SetUserValue(idx int)               // r[idx].userValue = pop()
	RawSet(idx int)                     // v=pop(); k=pop(); r[idx][k] = v
	RawSetI(idx int, i int64)           // r[idx][n] = pop()
	RawSetP(idx int, p UserData)        // r[idx][p] = pop()
	/* 'load' and 'call' functions (load and run Lua code) */
	CallK()                                                 //
	Call(nArgs, nResults int)                               // args=pop(nArgs); f=pop(); f(args)
	PCall(nArgs, nResults, msgh int) ThreadStatus           // call(nArgs, nResults) || push(err)
	PCallK()                                                //
	Load(chunk []byte, chunkName, mode string) ThreadStatus // push(compile(chunk))
	Dump()                                                  // todo
	/* miscellaneous functions */
	Error() int                   //
	Next(idx int) bool            // key=pop(); k,v=next(r[idx]); push(k,v);
	Concat(n int)                 // push(concat(pop(n)))
	Len(idx int)                  // push(len(r[idx]))
	StringToNumber(s string) bool // push(number(s))
	/* coroutine functions */
	NewThread() LuaState                          // todo
	Yield(nResults int) int                       // todo
	YieldK()                                      // todo
	Status() ThreadStatus                         // todo
	Resume(from LuaState, nArgs int) ThreadStatus // todo
	IsYieldable() bool                            // todo
	/* garbage-collection function and options */
	GC(what, data int) int //
}

// type LuaLightUserData UserData
// type LuaKContext int  // todo
// type LuaKFunction int // todo
// type LuaReader int    // todo
// type LuaWriter int    // todo

//GetAllocf()
//GetExtraSpace()
//IsLightUserData(idx int)
//NewUserData(size uint)
//PushLightUserData()
//PushLiteral
//PushLString
//SetAllocf()
//ToLString(idx int) (string, uint)
