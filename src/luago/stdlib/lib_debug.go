package stdlib

import "strings"
import . "luago/api"

var dbLib = map[string]GoFunction{
	"debug":        dbDebug,
	"getinfo":      dbGetInfo,
	"getregistry":  dbGetRegistry,
	"traceback":    dbTraceback,
	"gethook":      dbGetHook,
	"sethook":      dbSetHook,
	"getlocal":     dbGetLocal,
	"setlocal":     dbSetLocal,
	"getmetatable": dbGetMetatable,
	"setmetatable": dbSetMetatable,
	"getupvalue":   dbGetUpvalue,
	"setupvalue":   dbSetUpvalue,
	"upvalueid":    dbUpvalueId,
	"upvaluejoin":  dbUpvalueJoin,
	"getuservalue": dbGetUserValue,
	"setuservalue": dbSetUserValue,
}

func OpenDebugLib(ls LuaState) int {
	ls.NewLib(dbLib)
	return 1
}

func dbDebug(ls LuaState) int {
	panic("todo: dbDebug!")
}

// debug.getinfo ([thread,] f [, what])
// http://www.lua.org/manual/5.3/manual.html#pdf-debug.getinfo
// lua-5.3.4/src/ldblib.c#db_getinfo()
func dbGetInfo(L LuaState) int {
	ar := &LuaDebug{}
	var arg int
	L1 := getthread(L, &arg)
	options := luaL_optstring(L, arg+2, "flnStu")
	checkstack(L, L1, 3)
	if lua_isfunction(L, arg+1) { /* info about a function? */
		options = lua_pushfstring(L, ">%s", options) /* add '>' to 'options' */
		lua_pushvalue(L, arg+1)                      /* move function to 'L1' stack */
		lua_xmove(L, L1, 1)
	} else { /* stack level */
		if !lua_getstack(L1, int(luaL_checkinteger(L, arg+1)), ar) {
			lua_pushnil(L) /* level out of range */
			return 1
		}
	}
	if !lua_getinfo(L1, options, ar) {
		return luaL_argerror(L, arg+2, "invalid option")
	}
	lua_newtable(L) /* table to collect results */
	if strchr(options, 'S') {
		settabss(L, "source", ar.Source)
		settabss(L, "short_src", ar.ShortSrc)
		settabsi(L, "linedefined", ar.LineDefined)
		settabsi(L, "lastlinedefined", ar.LastLineDefined)
		settabss(L, "what", ar.What)
	}
	if strchr(options, 'l') {
		settabsi(L, "currentline", ar.CurrentLine)
	}
	if strchr(options, 'u') {
		settabsi(L, "nups", ar.NUps)
		settabsi(L, "nparams", ar.NParams)
		settabsb(L, "isvararg", ar.IsVararg)
	}
	if strchr(options, 'n') {
		settabss(L, "name", ar.Name)
		settabss(L, "namewhat", ar.NameWhat)
	}
	if strchr(options, 't') {
		settabsb(L, "istailcall", ar.IsTailCall)
	}
	if strchr(options, 'L') {
		treatstackoption(L, L1, "activelines")
	}
	if strchr(options, 'f') {
		treatstackoption(L, L1, "func")
	}
	return 1 /* return table */
}

// debug.getregistry ()
// http://www.lua.org/manual/5.3/manual.html#pdf-debug.getregistry
// lua-5.3.4/src/ldblib.c#db_getregistry()
func dbGetRegistry(ls LuaState) int {
	ls.PushValue(LUA_REGISTRYINDEX)
	return 1
}

func dbTraceback(ls LuaState) int {
	panic("todo: dbTraceback!")
}

func dbGetHook(L LuaState) int {
	var arg int
	L1 := getthread(L, &arg)
	mask := lua_gethookmask(L1)
	hook := lua_gethook(L1)
	if hook == nil { /* no hook? */
		lua_pushnil(L)
		// } else if hook != hookf { /* external hook? */
		// 	lua_pushliteral(L, "external hook")
		// } else { /* hook table must exist */
		// 	lua_rawgetp(L, LUA_REGISTRYINDEX, &HOOKKEY)
		// 	checkstack(L, L1, 1)
		// 	lua_pushthread(L1)
		// 	lua_xmove(L1, L, 1)
		// 	lua_rawget(L, -2) /* 1st result = hooktable[L1] */
		// 	lua_remove(L, -2) /* remove hook table */
	}
	lua_pushstring(L, unmakemask(mask))             /* 2nd result = mask */
	lua_pushinteger(L, int64(lua_gethookcount(L1))) /* 3rd result = count */
	return 3
}

func dbSetHook(ls LuaState) int {
	panic("todo: dbSetHook!")
}

func dbGetLocal(ls LuaState) int {
	panic("todo: dbGetLocal!")
}

func dbSetLocal(ls LuaState) int {
	panic("todo: dbSetLocal!")
}

// debug.getmetatable (value)
// http://www.lua.org/manual/5.3/manual.html#pdf-debug.getmetatable
// lua-5.3.4/src/ldblib.c#db_getmetatable()
func dbGetMetatable(ls LuaState) int {
	ls.CheckAny(1)
	if !ls.GetMetatable(1) {
		ls.PushNil() /* no metatable */
	}
	return 1
}

// debug.setmetatable (value, table)
// http://www.lua.org/manual/5.3/manual.html#pdf-debug.setmetatable
// lua-5.3.4/src/ldblib.c#db_setmetatable()
func dbSetMetatable(ls LuaState) int {
	t := ls.Type(2)
	ls.ArgCheck(t == LUA_TNIL || t == LUA_TTABLE, 2,
		"nil or table expected")
	ls.SetTop(2)
	ls.SetMetatable(1)
	return 1 /* return 1st argument */
}

// debug.getupvalue (f, up)
// http://www.lua.org/manual/5.3/manual.html#pdf-debug.getupvalue
// lua-5.3.4/src/ldblib.c#db_getupvalue()
func dbGetUpvalue(ls LuaState) int {
	return _auxUpvalue(ls, 1)
}

// debug.setupvalue (f, up, value)
// http://www.lua.org/manual/5.3/manual.html#pdf-debug.setupvalue
// lua-5.3.4/src/ldblib.c#db_setupvalue()
func dbSetUpvalue(ls LuaState) int {
	ls.CheckAny(3)
	return _auxUpvalue(ls, 0)
}

func _auxUpvalue(ls LuaState, get int) int {
	n := int(ls.CheckInteger(2))   /* upvalue index */
	ls.CheckType(1, LUA_TFUNCTION) /* closure */
	var name string
	if get > 0 {
		name = ls.GetUpvalue(1, n)
	} else {
		name = ls.SetUpvalue(1, n)
	}
	if name == "" {
		return 0
	}
	ls.PushString(name)
	ls.Insert(-(get + 1)) /* no-op if get is false */
	return get + 1
}

// debug.upvaluejoin (f1, n1, f2, n2)
// http://www.lua.org/manual/5.3/manual.html#pdf-debug.upvaluejoin
// lua-5.3.4/src/ldblib.c#db_upvaluejoin()
func dbUpvalueJoin(ls LuaState) int {
	n1 := _checkUpval(ls, 1, 2)
	n2 := _checkUpval(ls, 3, 4)
	ls.ArgCheck(!ls.IsGoFunction(1), 1, "Lua function expected")
	ls.ArgCheck(!ls.IsGoFunction(3), 3, "Lua function expected")
	ls.UpvalueJoin(1, n1, 3, n2)
	return 0
}

func _checkUpval(ls LuaState, argf, argnup int) int {
	nup := int(ls.CheckInteger(argnup)) /* upvalue index */
	ls.CheckType(argf, LUA_TFUNCTION)   /* closure */
	ls.ArgCheck((ls.GetUpvalue(argf, nup) != ""), argnup,
		"invalid upvalue index")
	return nup
}

// debug.upvalueid (f, n)
// http://www.lua.org/manual/5.3/manual.html#pdf-debug.upvalueid
// lua-5.3.4/src/ldblib.c#db_upvalueid()
func dbUpvalueId(ls LuaState) int {
	n := _checkUpval(ls, 1, 2)
	ls.PushLightUserData(ls.UpvalueId(1, n))
	return 1
}

func dbSetUserValue(ls LuaState) int {
	panic("todo: dbSetUserValue!")
}

func dbGetUserValue(ls LuaState) int {
	panic("todo: dbGetUserValue!")
}

/*
** If L1 != L, L1 can be in any state, and therefore there are no
** guarantees about its stack space; any push in L1 must be
** checked.
 */
func checkstack(L, L1 LuaState, n int) {
	if L != L1 && !lua_checkstack(L1, n) {
		luaL_error(L, "stack overflow")
	}
}

/*
** Auxiliary function used by several library functions: check for
** an optional thread as function's first argument and set 'arg' with
** 1 if this argument is present (so that functions can skip it to
** access their other arguments)
 */
func getthread(L LuaState, arg *int) LuaState {
	if lua_isthread(L, 1) {
		*arg = 1
		return lua_tothread(L, 1)
	} else {
		*arg = 0
		return L /* function will operate over current thread */
	}
}

/*
** Variations of 'lua_settable', used by 'db_getinfo' to put results
** from 'lua_getinfo' into result table. Key is always a string;
** value can be a string, an int, or a boolean.
 */
func settabss(L LuaState, k string, v string) {
	lua_pushstring(L, v)
	lua_setfield(L, -2, k)
}

func settabsi(L LuaState, k string, v int) {
	lua_pushinteger(L, int64(v))
	lua_setfield(L, -2, k)
}

func settabsb(L LuaState, k string, v bool) {
	lua_pushboolean(L, v)
	lua_setfield(L, -2, k)
}

/*
** In function 'db_getinfo', the call to 'lua_getinfo' may push
** results on the stack; later it creates the result table to put
** these objects. Function 'treatstackoption' puts the result from
** 'lua_getinfo' on top of the result table so that it can call
** 'lua_setfield'.
 */
func treatstackoption(L, L1 LuaState, fname string) {
	if L == L1 {
		lua_rotate(L, -2, 1) /* exchange object and table */
	} else {
		lua_xmove(L1, L, 1) /* move object to the "main" stack */
	}
	lua_setfield(L, -2, fname) /* put object into table */
}

/*
** Convert a bit mask (for 'gethook') into a string mask
 */
func unmakemask(mask int) string {
	smask := ""
	if mask&LUA_MASKCALL != 0 {
		smask += "c"
	}
	if mask&LUA_MASKRET != 0 {
		smask += "r"
	}
	if mask&LUA_MASKLINE != 0 {
		smask += "l"
	}
	return smask
}

func strchr(s string, b byte) bool {
	return strings.IndexByte(s, b) >= 0
}

func hookf(ls LuaState, ar *LuaDebug) {
	panic("todo!")
}
