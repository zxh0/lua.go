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
func dbGetInfo(ls LuaState) int {
	ar := &LuaDebug{}
	arg, ls1 := _getThread(ls)
	options := ls.OptString(arg+2, "flnStu")
	_checkStack(ls, ls1, 3)
	if ls.IsFunction(arg + 1) { /* info about a function? */
		options = ">" + options /* add '>' to 'options' */
		ls.PushString(options)
		ls.PushValue(arg + 1) /* move function to 'L1' stack */
		ls.XMove(ls1, 1)
	} else { /* stack level */
		if !ls1.GetStack(int(ls.CheckInteger(arg+1)), ar) {
			ls.PushNil() /* level out of range */
			return 1
		}
	}
	if !ls1.GetInfo(options, ar) {
		return ls.ArgError(arg+2, "invalid option")
	}
	ls.NewTable() /* table to collect results */
	if strings.IndexByte(options, 'S') >= 0 {
		_setTabSS(ls, "source", ar.Source)
		_setTabSS(ls, "short_src", ar.ShortSrc)
		_setTabSI(ls, "linedefined", ar.LineDefined)
		_setTabSI(ls, "lastlinedefined", ar.LastLineDefined)
		_setTabSS(ls, "what", ar.What)
	}
	if strings.IndexByte(options, 'l') >= 0 {
		_setTabSI(ls, "currentline", ar.CurrentLine)
	}
	if strings.IndexByte(options, 'u') >= 0 {
		_setTabSI(ls, "nups", ar.NUps)
		_setTabSI(ls, "nparams", ar.NParams)
		_setTabSB(ls, "isvararg", ar.IsVararg)
	}
	if strings.IndexByte(options, 'n') >= 0 {
		_setTabSS(ls, "name", ar.Name)
		_setTabSS(ls, "namewhat", ar.NameWhat)
	}
	if strings.IndexByte(options, 't') >= 0 {
		_setTabSB(ls, "istailcall", ar.IsTailCall)
	}
	if strings.IndexByte(options, 'L') >= 0 {
		_treatStackOption(ls, ls1, "activelines")
	}
	if strings.IndexByte(options, 'f') >= 0 {
		_treatStackOption(ls, ls1, "func")
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

func dbGetHook(ls LuaState) int {
	panic("todo: dbGetHook!")
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
func _checkStack(ls, ls1 LuaState, n int) {
	if ls != ls1 && !ls1.CheckStack(n) {
		ls.Error2("stack overflow")
	}
}

/*
** Auxiliary function used by several library functions: check for
** an optional thread as function's first argument and set 'arg' with
** 1 if this argument is present (so that functions can skip it to
** access their other arguments)
 */
func _getThread(ls LuaState) (arg int, ls1 LuaState) {
	if ls.IsThread(1) {
		return 1, ls.ToThread(1)
	}
	return 0, ls /* function will operate over current thread */
}

/*
** Variations of 'lua_settable', used by 'db_getinfo' to put results
** from 'lua_getinfo' into result table. Key is always a string;
** value can be a string, an int, or a boolean.
 */
func _setTabSS(ls LuaState, k string, v string) {
	ls.PushString(v)
	ls.SetField(-2, k)
}

func _setTabSI(ls LuaState, k string, v int) {
	ls.PushInteger(int64(v))
	ls.SetField(-2, k)
}

func _setTabSB(ls LuaState, k string, v bool) {
	ls.PushBoolean(v)
	ls.SetField(-2, k)
}

/*
** In function 'db_getinfo', the call to 'lua_getinfo' may push
** results on the stack; later it creates the result table to put
** these objects. Function 'treatstackoption' puts the result from
** 'lua_getinfo' on top of the result table so that it can call
** 'lua_setfield'.
 */
func _treatStackOption(ls, ls1 LuaState, fname string) {
	if ls == ls1 {
		ls.Rotate(-2, 1) /* exchange object and table */
	} else {
		ls1.XMove(ls, 1) /* move object to the "main" stack */
	}
	ls.SetField(-2, fname) /* put object into table */
}
