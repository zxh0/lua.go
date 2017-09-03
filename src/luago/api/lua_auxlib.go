package api

type FuncReg map[string]GoFunction

// auxiliary library
type AuxLib interface {
	ArgCheck(cond bool, arg int, extraMsg string)        //
	ArgError(arg int, extraMsg string) int               //
	CheckAny(arg int)                                    // r[arg] is None ?
	CheckInteger(arg int) int64                          // r[arg] is LuaInteger ?
	CheckNumber(arg int) float64                         // r[arg] is LuaNumber ?
	CheckStack2(sz int, msg string)                      //
	CheckString(arg int) string                          // r[arg] is string ?
	CheckType(arg int, t LuaType)                        // r[arg] is LuaType ?
	CheckVersion()                                       //
	DoFile(filename string) bool                         //
	DoString(str string) bool                            //
	Error2(fmt string)                                   // todo
	GetMetaField(obj int, e string) LuaType              //
	GetMetatableL(tname string) LuaType                  //
	GetSubTable(idx int, fname string) bool              // push(r[idx][fname] || {})
	Len2(idx int) int64                                  // #(r[idx])
	LoadFile(filename string) ThreadStatus               //
	LoadFileX(filename, mode string) ThreadStatus        //
	LoadString(s string) ThreadStatus                    //
	NewLib(l FuncReg)                                    //
	NewLibTable(l FuncReg)                               //
	OpenLibs()                                           //
	OptInteger(arg int, d int64) int64                   // r[arg] or d
	OptNumber(arg int, d float64) float64                // r[arg] or d
	OptString(arg int, d string) string                  // r[arg] or d
	RequireF(modname string, openf GoFunction, glb bool) //
	SetFuncs(l FuncReg, nup int)                         // l.each{name,func => r[-1][name]=func}
	TypeName2(idx int) string                            //
}
