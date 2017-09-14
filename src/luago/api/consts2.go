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

const LUA_PATH_DEFAULT = "?.lua" // todo

/*
** maximum number of upvalues in a closure (both C and Lua). (Value
** must fit in a VM register.)
 */
const MAXUPVAL = 255

func LuaUpvalueIndex(i int) int {
	return LUA_REGISTRYINDEX - i
}
