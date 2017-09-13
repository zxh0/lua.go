package api

type FuncReg map[string]GoFunction

// auxiliary library
type AuxLib interface {
	/* Error-report functions */
	Error2(fmt string, a ...interface{}) int // todo
	ArgError(arg int, extraMsg string) int   // todo
	Where(lvl int)                           //
	/* Argument check functions */
	CheckStack2(sz int, msg string)               //
	ArgCheck(cond bool, arg int, extraMsg string) //
	CheckAny(arg int)                             // r[arg] is None ?
	CheckInteger(arg int) int64                   // r[arg] is LuaInteger ?
	CheckNumber(arg int) float64                  // r[arg] is LuaNumber ?
	CheckString(arg int) string                   // r[arg] is string ?
	CheckType(arg int, t LuaType)                 // r[arg] is LuaType ?
	OptInteger(arg int, d int64) int64            // r[arg] or d
	OptNumber(arg int, d float64) float64         // r[arg] or d
	OptString(arg int, d string) string           // r[arg] or d
	/* Load functions */
	DoFile(filename string) bool                  //
	DoString(str string) bool                     //
	LoadFile(filename string) ThreadStatus        //
	LoadFileX(filename, mode string) ThreadStatus //
	LoadString(s string) ThreadStatus             //
	/* Other functions */
	GetMetatable2(tname string) LuaType                  // v=registry[tname]; push(v.mt)
	GetMetafield(obj int, e string) LuaType              // v=r[obj]; mt=v.mt; f=mt[e]; push(f)
	CallMeta(obj int, e string) bool                     // v=r[obj]; mt=v.mt; f=mt[e]; f(v)
	OpenLibs()                                           //
	RequireF(modname string, openf GoFunction, glb bool) //
	NewLib(l FuncReg)                                    //
	NewLibTable(l FuncReg)                               //
	SetFuncs(l FuncReg, nup int)                         // l.each{name,func => r[-1][name]=func}
	GetSubTable(idx int, fname string) bool              // push(r[idx][fname] || {})
	Len2(idx int) int64                                  // #(r[idx])
	TypeName2(idx int) string                            // typename(type(idx))
	ToString2(idx int) string                            //
	CheckVersion()                                       //
}

// luaL_fileresult
// luaL_execresult
// luaL_checkoption
// luaL_newmetatable
// luaL_setmetatable
// luaL_checkudata
// luaL_testudata
// luaL_traceback
// luaL_gsub
// luaL_newstate
// luaL_ref
// luaL_unref
