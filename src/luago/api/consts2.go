package api

const (
	LUA_VERSION_MAJOR   = "5"
	LUA_VERSION_MINOR   = "3"
	LUA_VERSION_NUM     = 503
	LUA_VERSION_RELEASE = "3"

	LUA_VERSION = "Lua " + LUA_VERSION_MAJOR + "." + LUA_VERSION_MINOR
	LUA_RELEASE = LUA_VERSION + "." + LUA_VERSION_RELEASE
	//LUA_COPYRIGHT	LUA_RELEASE "  Copyright (C) 1994-2016 Lua.org, PUC-Rio"
	//LUA_AUTHORS	"R. Ierusalimschy, L. H. de Figueiredo, W. Celes"
)

/*
@@ LUAI_MAXSTACK limits the size of the Lua stack.
** CHANGE it if you need a different limit. This limit is arbitrary;
** its only purpose is to stop Lua from consuming unlimited stack
** space (and to reserve some numbers for pseudo-indices).
*/
const LUAI_MAXSTACK = 1000000

const LUA_PATH_DEFAULT = "?.lua" // todo

/* option for multiple returns in 'lua_pcall' and 'lua_call' */
const LUA_MULTRET = -1

/*
** Pseudo-indices
** (-LUAI_MAXSTACK is the minimum valid index; we keep some free empty
** space after that to help overflow detection)
 */
const LUA_REGISTRYINDEX = -LUAI_MAXSTACK - 1000

func LuaUpvalueIndex(i int) int {
	return LUA_REGISTRYINDEX - i
}

/* minimum Lua stack available to a C function */
const LUA_MINSTACK = 20

/* predefined values in the registry */
const LUA_RIDX_MAINTHREAD int64 = 1
const LUA_RIDX_GLOBALS int64 = 2
const LUA_RIDX_LAST = LUA_RIDX_GLOBALS

/*
** maximum number of upvalues in a closure (both C and Lua). (Value
** must fit in a VM register.)
 */
const MAXUPVAL = 255
