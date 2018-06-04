package stdlib

import . "luago/api"

type any = interface{}

/* state manipulation */
func lua_close(ls LuaState)                            { ls.Close() }
func lua_newthread(ls LuaState) LuaState               { return ls.NewThread() }
func lua_version(ls LuaState) float64                  { return ls.Version() }
func lua_atpanic(ls LuaState, f GoFunction) GoFunction { return ls.AtPanic(f) }

/* basic stack manipulation */
func lua_absindex(ls LuaState, idx int) int       { return ls.AbsIndex(idx) }
func lua_gettop(ls LuaState) int                  { return ls.GetTop() }
func lua_settop(ls LuaState, idx int)             { ls.SetTop(idx) }
func lua_pushvalue(ls LuaState, idx int)          { ls.PushValue(idx) }
func lua_rotate(ls LuaState, idx, n int)          { ls.Rotate(idx, n) }
func lua_copy(ls LuaState, fromIdx, toIdx int)    { ls.Copy(fromIdx, toIdx) }
func lua_checkstack(ls LuaState, n int) bool      { return ls.CheckStack(n) }
func lua_xmove(from LuaState, to LuaState, n int) { from.XMove(to, n) }
func lua_pop(ls LuaState, n int)                  { ls.Pop(n) }
func lua_insert(ls LuaState, idx int)             { ls.Insert(idx) }
func lua_remove(ls LuaState, idx int)             { ls.Remove(idx) }
func lua_replace(ls LuaState, idx int)            { ls.Replace(idx) }

/* access functions (stack -> C) */
func lua_isnumber(ls LuaState, idx int) bool             { return ls.IsNumber(idx) }
func lua_isstring(ls LuaState, idx int) bool             { return ls.IsString(idx) }
func lua_isgofunction(ls LuaState, idx int) bool         { return ls.IsGoFunction(idx) }
func lua_isinteger(ls LuaState, idx int) bool            { return ls.IsInteger(idx) }
func lua_isuserdata(ls LuaState, idx int) bool           { return ls.IsUserData(idx) }
func lua_type(ls LuaState, idx int) LuaType              { return ls.Type(idx) }
func lua_typename(ls LuaState, tp LuaType) string        { return ls.TypeName(tp) }
func lua_tonumberx(ls LuaState, idx int) (float64, bool) { return ls.ToNumberX(idx) }
func lua_tointegerx(ls LuaState, idx int) (int64, bool)  { return ls.ToIntegerX(idx) }
func lua_toboolean(ls LuaState, idx int) bool            { return ls.ToBoolean(idx) }
func lua_tostring(ls LuaState, idx int) string           { return ls.ToString(idx) }
func lua_tostringx(ls LuaState, idx int) (string, bool)  { return ls.ToStringX(idx) }
func lua_rawlen(ls LuaState, idx int) uint               { return ls.RawLen(idx) }
func lua_togofunction(ls LuaState, idx int) GoFunction   { return ls.ToGoFunction(idx) }
func lua_touserdata(ls LuaState, idx int) UserData       { return ls.ToUserData(idx) }
func lua_tothread(ls LuaState, idx int) LuaState         { return ls.ToThread(idx) }
func lua_topointer(ls LuaState, idx int) any             { return ls.ToPointer(idx) }
func lua_tonumber(ls LuaState, idx int) float64          { return ls.ToNumber(idx) }
func lua_tointeger(ls LuaState, idx int) int64           { return ls.ToInteger(idx) }
func lua_isfunction(ls LuaState, idx int) bool           { return ls.IsFunction(idx) }
func lua_istable(ls LuaState, idx int) bool              { return ls.IsTable(idx) }
func lua_isnil(ls LuaState, idx int) bool                { return ls.IsNil(idx) }
func lua_isboolean(ls LuaState, idx int) bool            { return ls.IsBoolean(idx) }
func lua_isthread(ls LuaState, idx int) bool             { return ls.IsThread(idx) }
func lua_isnone(ls LuaState, idx int) bool               { return ls.IsNone(idx) }
func lua_isnoneornil(ls LuaState, idx int) bool          { return ls.IsNoneOrNil(idx) }

/* push functions (C -> stack) */
func lua_pushnil(ls LuaState)                                  { ls.PushNil() }
func lua_pushnumber(ls LuaState, n float64)                    { ls.PushNumber(n) }
func lua_pushinteger(ls LuaState, n int64)                     { ls.PushInteger(n) }
func lua_pushstring(ls LuaState, s string)                     { ls.PushString(s) }
func lua_pushfstring(ls LuaState, fmt string, a ...any) string { return ls.PushFString(fmt, a...) }
func lua_pushgoclosure(ls LuaState, f GoFunction, n int)       { ls.PushGoClosure(f, n) }
func lua_pushboolean(ls LuaState, b bool)                      { ls.PushBoolean(b) }
func lua_pushlightuserdata(ls LuaState, d UserData)            { ls.PushLightUserData(d) }
func lua_pushthread(ls LuaState) bool                          { return ls.PushThread() }
func lua_pushgofunction(ls LuaState, f GoFunction)             { ls.PushGoFunction(f) }
func lua_pushliteral(ls LuaState, s string)                    { ls.PushString(s) }
func lua_pushglobaltable(ls LuaState)                          { ls.PushGlobalTable() }

/* Comparison and arithmetic functions */
func lua_arith(ls LuaState, op ArithOp)                          { ls.Arith(op) }
func lua_compare(ls LuaState, idx1, idx2 int, op CompareOp) bool { return ls.Compare(idx1, idx2, op) }
func lua_rawequal(ls LuaState, idx1, idx2 int) bool              { return ls.RawEqual(idx1, idx2) }

/* get functions (Lua -> stack) */
func lua_getglobal(ls LuaState, name string) LuaType       { return ls.GetGlobal(name) }
func lua_gettable(ls LuaState, idx int) LuaType            { return ls.GetTable(idx) }
func lua_getfield(ls LuaState, idx int, k string) LuaType  { return ls.GetField(idx, k) }
func lua_geti(ls LuaState, idx int, i int64) LuaType       { return ls.GetI(idx, i) }
func lua_rawget(ls LuaState, idx int) LuaType              { return ls.RawGet(idx) }
func lua_rawgeti(ls LuaState, idx int, i int64) LuaType    { return ls.RawGetI(idx, i) }
func lua_rawgetp(ls LuaState, idx int, p UserData) LuaType { return ls.RawGetP(idx, p) }
func lua_createtable(ls LuaState, nArr, nRec int)          { ls.CreateTable(nArr, nRec) }
func lua_getmetatable(ls LuaState, idx int) bool           { return ls.GetMetatable(idx) }
func lua_getuservalue(ls LuaState, idx int) LuaType        { return ls.GetUserValue(idx) }
func lua_newtable(ls LuaState)                             { ls.NewTable() }

/* set functions (stack -> Lua) */
func lua_setglobal(ls LuaState, name string)              { ls.SetGlobal(name) }
func lua_settable(ls LuaState, idx int)                   { ls.SetTable(idx) }
func lua_setfield(ls LuaState, idx int, k string)         { ls.SetField(idx, k) }
func lua_seti(ls LuaState, idx int, i int64)              { ls.SetI(idx, i) }
func lua_rawset(ls LuaState, idx int)                     { ls.RawSet(idx) }
func lua_rawseti(ls LuaState, idx int, i int64)           { ls.RawSetI(idx, i) }
func lua_rawsetp(ls LuaState, idx int, p UserData)        { ls.RawSetP(idx, p) }
func lua_setmetatable(ls LuaState, idx int)               { ls.SetMetatable(idx) }
func lua_setuservalue(ls LuaState, idx int)               { ls.SetUserValue(idx) }
func lua_register(ls LuaState, name string, f GoFunction) { ls.Register(name, f) }

/* 'load' and 'call' functions (load and run Lua code) */
func lua_callk(ls LuaState)                                      { ls.CallK() }
func lua_call(ls LuaState, nArgs, nResults int)                  { ls.Call(nArgs, nResults) }
func lua_pcallk(ls LuaState)                                     { ls.PCallK() }
func lua_pcall(ls LuaState, nArgs, nRets, msgh int) ThreadStatus { return ls.PCall(nArgs, nRets, msgh) }
func lua_load(ls LuaState, c []byte, cn, m string) ThreadStatus  { return ls.Load(c, cn, m) }
func lua_dump(ls LuaState, strip bool) []byte                    { return ls.Dump(strip) }

/* coroutine functions */
func lua_yieldk(ls LuaState)                                        { ls.YieldK() }
func lua_resume(ls LuaState, from LuaState, nArgs int) ThreadStatus { return ls.Resume(from, nArgs) }
func lua_status(ls LuaState) ThreadStatus                           { return ls.Status() }
func lua_isyieldable(ls LuaState) bool                              { return ls.IsYieldable() }
func lua_yield(ls LuaState, nResults int) int                       { return ls.Yield(nResults) }

/* garbage-collection function and options */
func lua_gc(ls LuaState, what, data int) int { return ls.GC(what, data) }

/* miscellaneous functions */
func lua_concat(ls LuaState, n int)                 { ls.Concat(n) }
func lua_len(ls LuaState, idx int)                  { ls.Len(idx) }
func lua_next(ls LuaState, idx int) bool            { return ls.Next(idx) }
func lua_stringtonumber(ls LuaState, s string) bool { return ls.StringToNumber(s) }
func lua_error(ls LuaState) int                     { return ls.Error() }

/* Debug API */
func lua_getstack(ls LuaState, level int, ar *LuaDebug) bool  { return ls.GetStack(level, ar) }
func lua_getinfo(ls LuaState, what string, ar *LuaDebug) bool { return ls.GetInfo(what, ar) }
func lua_getlocal(ls LuaState, ar *LuaDebug, n int) string    { return ls.GetLocal(ar, n) }
func lua_setlocal(ls LuaState, ar *LuaDebug, n int) string    { return ls.SetLocal(ar, n) }
func lua_getupvalue(ls LuaState, funcIdx, n int) string       { return ls.GetUpvalue(funcIdx, n) }
func lua_setupvalue(ls LuaState, funcIdx, n int) string       { return ls.SetUpvalue(funcIdx, n) }
func lua_upvalueid(ls LuaState, funcIdx, n int) any           { return ls.UpvalueId(funcIdx, n) }
func lua_upvaluejoin(ls LuaState, f1, n1, f2, n2 int)         { ls.UpvalueJoin(f1, n1, f2, n2) }
func lua_gethook(ls LuaState) LuaHook                         { return ls.GetHook() }
func lua_sethook(ls LuaState, f LuaHook, mask, count int)     { ls.SetHook(f, mask, count) }
func lua_gethookmask(ls LuaState) int                         { return ls.GetHookMask() }
func lua_gethookcount(ls LuaState) int                        { return ls.GetHookCount() }

/* auxiliary library */
func luaL_checkversion(ls LuaState)                                  { ls.CheckVersion() }
func luaL_getmetafield(ls LuaState, obj int, e string) LuaType       { return ls.GetMetafield(obj, e) }
func luaL_callmeta(ls LuaState, obj int, e string) bool              { return ls.CallMeta(obj, e) }
func luaL_tostring(ls LuaState, idx int) string                      { return ls.ToString2(idx) }
func luaL_argerror(ls LuaState, arg int, extraMsg string) int        { return ls.ArgError(arg, extraMsg) }
func luaL_checkstring(ls LuaState, arg int) string                   { return ls.CheckString(arg) }
func luaL_optstring(ls LuaState, arg int, d string) string           { return ls.OptString(arg, d) }
func luaL_checknumber(ls LuaState, arg int) float64                  { return ls.CheckNumber(arg) }
func luaL_optnumber(ls LuaState, arg int, d float64) float64         { return ls.OptNumber(arg, d) }
func luaL_checkinteger(ls LuaState, arg int) int64                   { return ls.CheckInteger(arg) }
func luaL_optinteger(ls LuaState, arg int, d int64) int64            { return ls.OptInteger(arg, d) }
func luaL_checkstack(ls LuaState, sz int, msg string)                { ls.CheckStack2(sz, msg) }
func luaL_checktype(ls LuaState, arg int, t LuaType)                 { ls.CheckType(arg, t) }
func luaL_checkany(ls LuaState, arg int)                             { ls.CheckAny(arg) }
func luaL_where(ls LuaState, lvl int)                                { ls.Where(lvl) }
func luaL_error(ls LuaState, fmt string, a ...any) int               { return ls.Error2(fmt, a...) }
func luaL_loadfilex(ls LuaState, fn, mode string) ThreadStatus       { return ls.LoadFileX(fn, mode) }
func luaL_loadfile(ls LuaState, filename string) ThreadStatus        { return ls.LoadFile(filename) }
func luaL_loadstring(ls LuaState, s string) ThreadStatus             { return ls.LoadString(s) }
func luaL_len(ls LuaState, idx int) int64                            { return ls.Len2(idx) }
func luaL_setfuncs(ls LuaState, l FuncReg, nup int)                  { ls.SetFuncs(l, nup) }
func luaL_getsubtable(ls LuaState, idx int, fname string) bool       { return ls.GetSubTable(idx, fname) }
func luaL_requiref(ls LuaState, n string, f GoFunction, glb bool)    { ls.RequireF(n, f, glb) }
func luaL_newlibtable(ls LuaState, l FuncReg)                        { ls.NewLibTable(l) }
func luaL_newlib(ls LuaState, l FuncReg)                             { ls.NewLib(l) }
func luaL_argcheck(ls LuaState, cond bool, arg int, extraMsg string) { ls.ArgCheck(cond, arg, extraMsg) }
func luaL_typename(ls LuaState, idx int) string                      { return ls.TypeName2(idx) }
func luaL_dofile(ls LuaState, filename string) bool                  { return ls.DoFile(filename) }
func luaL_dostring(ls LuaState, str string) bool                     { return ls.DoString(str) }
func luaL_openlibs(ls LuaState)                                      { ls.OpenLibs() }
