package api

type LuaType int
type LuaArithOp int
type LuaCompareOp int
type LuaThreadStatus int

// type LuaBoolean bool
// type LuaInteger int64
// type LuaNumber float64
type LuaGoFunction func(LuaState) int
type LuaUserData interface{}
type LuaRegMap map[string]LuaGoFunction

type LuaState interface {
	BasicAPI
	AuxLib
	String() string // debug
}

type BasicAPI interface {
	AbsIndex(idx int) int                                      // abs(idx)
	Arith(op LuaArithOp)                                       // b=pop(); a=pop(); push(a op b)
	AtPanic(panicf LuaGoFunction) LuaGoFunction                //
	Call(nArgs, nResults int)                                  // args=pop(nArgs); f=pop(); f(args)
	CallK()                                                    //
	CheckStack(n int) bool                                     //
	Close()                                                    //
	Compare(index1, index2 int, op LuaCompareOp) bool          // r[index1] op r[index2]
	Concat(n int)                                              // push(concat(pop(n)))
	Copy(fromIdx, toIdx int)                                   // r[toIdx] = r[fromidx]
	CreateTable(nArr, nRec int)                                // push({})
	Dump()                                                     //
	Error() int                                                //
	GC(what, data int) int                                     //
	GetField(index int, k string) LuaType                      // push(r[index][k])
	GetGlobal(name string) LuaType                             // push(global[name])
	GetI(index int, i int64) LuaType                           // push(r[index][i])
	GetMetaTable(index int) bool                               // push(r[index].metaTable)?
	GetTable(index int) LuaType                                // push(r[index][pop()])
	GetTop() int                                               // stack.top
	GetUserValue(index int) LuaType                            // push(r[index].userValue)
	Insert(index int)                                          // r[index, -1] >> 1
	IsBoolean(index int) bool                                  // r[index].type == LUA_TBOOLEAN
	IsGoFunction(index int) bool                               // r[index].type == LUA_TLCL || LUA_TGCL
	IsFunction(index int) bool                                 // r[index].type == LUA_TFUNCTION
	IsInteger(index int) bool                                  // r[index].type == LUA_TINTEGER
	IsNil(index int) bool                                      // r[index].type == LUA_TNIL
	IsNone(index int) bool                                     // r[index].type == LUA_TNONE
	IsNoneOrNil(index int) bool                                // r[index].type == LUA_TNONE || LUA_TNIL
	IsNumber(index int) bool                                   // r[index] ~= LuaNumber
	IsString(index int) bool                                   // r[index] ~= LuaString
	IsTable(index int) bool                                    // r[index].type == LUA_TTABLE
	IsThread(index int) bool                                   // r[index].type == LUA_TTHREAD
	IsUserData(index int) bool                                 // r[index].type == LUA_TUSERDATA
	IsYieldable() bool                                         //
	Len(index int)                                             // push(len(r[index]))
	Load(chunk []byte, chunkName, mode string) LuaThreadStatus // push(compile(chunk))
	NewTable()                                                 // push({})
	NewThread() LuaState                                       // todo
	Next(index int) bool                                       // key=pop(); k,v=next(r[index]); push(k,v);
	PCall(nArgs, nResults, msgh int) LuaThreadStatus           // call(nArgs, nResults) || push(err)
	PCallK()                                                   //
	Pop(n int)                                                 // pop(n)
	PushBoolean(b bool)                                        // push(b)
	PushGoClosure(fn LuaGoFunction, n int)                     // push(f)
	PushGoFunction(f LuaGoFunction)                            // push(f)
	PushFString(fmt string)                                    // todo
	PushGlobalTable()                                          // push(global)
	PushInteger(n int64)                                       // push(n)
	PushUserData(d LuaUserData)                                // push(d)
	PushNil()                                                  // push(nil)
	PushNumber(n float64)                                      // push(n)
	PushString(s string)                                       // push(s)
	PushThread(ls LuaState) bool                               // push(ls)
	PushValue(index int)                                       // push(r[index])
	PushVFString()                                             // todo
	RawEqual(index1, index2 int) bool                          // r[index1] == r[index2]
	RawGet(index int) LuaType                                  // push(r[index][pop()])
	RawGetI(index int, n int64) LuaType                        // push(r[index][i])
	RawGetP(index int, p LuaUserData) LuaType                  // push(r[index][p])
	RawLen(index int) uint                                     // len(r[index])
	RawSet(index int)                                          // v=pop(); k=pop(); r[index][k] = v
	RawSetI(index int, i int64)                                // r[index][n] = pop()
	RawSetP(index int, p LuaUserData)                          // r[index][p] = pop()
	Register(name string, f LuaGoFunction)                     // global[name] = f
	Remove(index int)                                          // remove(r[index])
	Replace(index int)                                         // r[index] = pop()
	Resume(from LuaState, nArgs int)                           // todo
	Rotate(idx, n int)                                         // r[idx, -1] >> n
	SetField(index int, k string)                              // r[index][k] = pop()
	SetGlobal(name string)                                     // global[name] = pop()
	SetI(index int, n int64)                                   // r[index][n] = pop()
	SetMetaTable(index int)                                    // r[index].metatable = pop()
	SetTable(index int)                                        // v=pop(); k=pop(); r[index][k] = v
	SetTop(index int)                                          // stack.top = index
	SetUserValue(index int)                                    // r[index].userValue = pop()
	Status() int                                               // todo
	StringToNumber(s string) bool                              // push(number(s))
	ToBoolean(index int) bool                                  // r[index] as bool
	ToGoFunction(index int) LuaGoFunction                      // r[index] as LuaGoFunction
	ToInteger(index int) int64                                 // r[index] as LuaInteger
	ToIntegerX(index int) (int64, bool)                        // r[index] as LuaInteger
	ToNumber(index int) float64                                // r[index] as LuaNumber
	ToNumberX(index int) (float64, bool)                       // r[index] as LuaNumber
	ToPointer(index int) interface{}                           // r[index] as interface{}
	ToString(index int) (string, bool)                         // r[index] as string
	ToThread(index int) LuaState                               // r[index] as LuaThread
	ToUserData(index int) LuaUserData                          // r[index] as LuaUserData
	Type(index int) LuaType                                    // r[index].type
	TypeName(tp LuaType) string                                // r[index].type.name
	Version() float64                                          // todo
	XMove(to LuaState, n int)                                  // to.push(pop(n))
	Yield(nResults int) int                                    // todo
	YieldK()                                                   // todo
}

// auxiliary library
type AuxLib interface {
	ArgCheck(cond bool, arg int, extraMsg string)           //
	ArgError(arg int, extraMsg string) int                  //
	CheckAny(arg int)                                       // r[arg] is None ?
	CheckInteger(arg int) int64                             // r[arg] is LuaInteger ?
	CheckNumber(arg int) float64                            // r[arg] is LuaNumber ?
	CheckStackL(sz int, msg string)                         //
	CheckString(arg int) string                             // r[arg] is string ?
	CheckType(arg int, t LuaType)                           // r[arg] is LuaType ?
	CheckVersion()                                          //
	DoFile(filename string) bool                            //
	DoString(str string) bool                               //
	ErrorL(fmt string)                                      // todo
	GetMetaField(obj int, e string) LuaType                 //
	GetMetaTableL(tname string) LuaType                     //
	GetSubTable(idx int, fname string) bool                 // push(r[idx][fname] || {})
	LenL(index int) int64                                   // #(r[index])
	LoadFile(filename string) LuaThreadStatus               //
	LoadFileX(filename, mode string) LuaThreadStatus        //
	LoadString(s string) LuaThreadStatus                    //
	NewLib(l LuaRegMap)                                     //
	NewLibTable(l LuaRegMap)                                //
	OpenLibs()                                              //
	OptInteger(arg int, d int64) int64                      // r[arg] or d
	OptNumber(arg int, d float64) float64                   // r[arg] or d
	OptString(arg int, d string) string                     // r[arg] or d
	RequireF(modname string, openf LuaGoFunction, glb bool) //
	SetFuncs(l LuaRegMap, nup int)                          // l.each{name,func => r[-1][name]=func}
	TypeNameL(index int) string                             //
}

// type LuaLightUserData LuaUserData
// type LuaKContext int  // todo
// type LuaKFunction int // todo
// type LuaReader int    // todo
// type LuaWriter int    // todo

//GetAllocf()
//GetExtraSpace()
//IsLightUserData(index int)
//NewUserData(size uint)
//PushLightUserData()
//PushLiteral
//PushLString
//SetAllocf()
//ToLString(index int) (string, uint)
