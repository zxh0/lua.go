package stdlib

import . "luago/api"

var dbLib = map[string]GoFunction{
	"debug":        dbDebug,
	"getuservalue": dbGetUserValue,
	"gethook":      dbGetHook,
	"getinfo":      dbGetInfo,
	"getlocal":     dbGetLocal,
	"getregistry":  dbGetRegistry,
	"getmetatable": dbGetMetatable,
	"getupvalue":   dbGetUpvalue,
	"upvaluejoin":  dbUpvalueJoin,
	"upvalueid":    dbUpvalueId,
	"setuservalue": dbSetUserValue,
	"sethook":      dbSetHook,
	"setlocal":     dbSetLocal,
	"setmetatable": dbSetMetatable,
	"setupvalue":   dbSetUpvalue,
	"traceback":    dbTraceback,
}

func OpenDebugLib(ls LuaState) int {
	ls.NewLib(dbLib)
	return 1
}

func dbDebug(ls LuaState) int {
	panic("todo: dbDebug!")
}

func dbGetUserValue(ls LuaState) int {
	panic("todo: dbGetUserValue!")
}

func dbGetHook(ls LuaState) int {
	panic("todo: dbGetHook!")
}

func dbGetInfo(ls LuaState) int {
	panic("todo: dbGetInfo!")
}

func dbGetLocal(ls LuaState) int {
	panic("todo: dbGetLocal!")
}

func dbGetRegistry(ls LuaState) int {
	panic("todo: dbGetRegistry!")
}

func dbGetMetatable(ls LuaState) int {
	panic("todo: dbGetMetatable!")
}

func dbGetUpvalue(ls LuaState) int {
	panic("todo: dbGetUpvalue!")
}

func dbUpvalueJoin(ls LuaState) int {
	panic("todo: dbUpvalueJoin!")
}

func dbUpvalueId(ls LuaState) int {
	panic("todo: dbUpvalueId!")
}

func dbSetUserValue(ls LuaState) int {
	panic("todo: dbSetUserValue!")
}

func dbSetHook(ls LuaState) int {
	panic("todo: dbSetHook!")
}

func dbSetLocal(ls LuaState) int {
	panic("todo: dbSetLocal!")
}

func dbSetMetatable(ls LuaState) int {
	panic("todo: dbSetMetatable!")
}

func dbSetUpvalue(ls LuaState) int {
	panic("todo: dbSetUpvalue!")
}

func dbTraceback(ls LuaState) int {
	panic("todo: dbTraceback!")
}
