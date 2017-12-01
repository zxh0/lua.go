package api

// Event codes
const (
	LUA_HOOKCALL = iota
	LUA_HOOKRET
	LUA_HOOKLINE
	LUA_HOOKCOUNT
	LUA_HOOKTAILCALL
)

// http://www.lua.org/manual/5.3/manual.html#lua_Debug
type LuaDebug struct {
	Event           int
	Name            string /* (n) */
	NameWhat        string /* (n) */
	What            string /* (S) */
	Source          string /* (S) */
	CurrentLine     int    /* (l) */
	LineDefined     int    /* (S) */
	LastLineDefined int    /* (S) */
	NUps            int    /* (u) number of upvalues */
	NParams         int    /* (u) number of parameters */
	IsVararg        bool   /* (u) */
	IsTailCall      bool   /* (t) */
	ShortSrc        string /* (S) */
	/* private part */
	// other fields
}

type LuaHook func(ls LuaState, ar *LuaDebug)

type DebugAPI interface {
	GetHook() LuaHook
	SetHook(f LuaHook, mask, count int)
	GetHookCount() int
	GetHookMask() int
	GetStack(level int, ar *LuaDebug) bool
	GetInfo(what string, ar *LuaDebug) bool
	GetLocal(ar *LuaDebug, n int) string
	SetLocal(ar *LuaDebug, n int) string
	GetUpvalue(funcIdx, n int) string
	SetUpvalue(funcIdx, n int) string
	UpvalueId(funcIdx, n int) interface{}
	UpvalueJoin(funcIdx1, n1, funcIdx2, n2 int)
}
