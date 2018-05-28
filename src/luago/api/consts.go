package api

// LUA_HOOKCALL
// LUA_HOOKCOUNT
// LUA_HOOKLINE
// LUA_HOOKRET
// LUA_HOOKTAILCALL
// LUA_MASKCALL
// LUA_MASKCOUNT
// LUA_MASKLINE
// LUA_MASKRET
// LUA_NOREF
// LUA_REFNIL
// LUA_USE_APICHECK
// LUAL_BUFFERSIZE

const (
	LUA_VERSION_MAJOR   = "5"
	LUA_VERSION_MINOR   = "3"
	LUA_VERSION_NUM     = 503
	LUA_VERSION_RELEASE = "4"

	LUA_VERSION = "Lua " + LUA_VERSION_MAJOR + "." + LUA_VERSION_MINOR
	LUA_RELEASE = LUA_VERSION + "." + LUA_VERSION_RELEASE
	//LUA_COPYRIGHT	LUA_RELEASE "  Copyright (C) 1994-2016 Lua.org, PUC-Rio"
	//LUA_AUTHORS	"R. Ierusalimschy, L. H. de Figueiredo, W. Celes"
)

/* option for multiple returns in 'lua_pcall' and 'lua_call' */
const LUA_MULTRET = -1

/* minimum Lua stack available to a C function */
const LUA_MINSTACK = 20

/*
** Pseudo-indices
** (-LUAI_MAXSTACK is the minimum valid index; we keep some free empty
** space after that to help overflow detection)
 */
const LUA_REGISTRYINDEX = -LUAI_MAXSTACK - 1000

/* predefined values in the registry */
const LUA_RIDX_MAINTHREAD int64 = 1
const LUA_RIDX_GLOBALS int64 = 2
const LUA_RIDX_LAST = LUA_RIDX_GLOBALS

// lua-5.3.4/src/lvm.c
/* limit for table tag-method chains (to avoid loops) */
const MAXTAGLOOP = 2000

/* basic types */
const (
	LUA_TNONE = iota - 1 // -1
	LUA_TNIL
	LUA_TBOOLEAN
	LUA_TLIGHTUSERDATA
	LUA_TNUMBER
	LUA_TSTRING
	LUA_TTABLE
	LUA_TFUNCTION
	LUA_TUSERDATA
	LUA_TTHREAD
)

// lua-5.3.4/src/lobject.h
/* type variants */
const (
	LUA_TNUMFLT = LUA_TNUMBER | (0 << 4)   // float numbers
	LUA_TNUMINT = LUA_TNUMBER | (1 << 4)   // integer numbers
	LUA_TSHRSTR = LUA_TSTRING | (0 << 4)   // short strings
	LUA_TLNGSTR = LUA_TSTRING | (1 << 4)   // long strings
	LUA_TLCL    = LUA_TFUNCTION | (0 << 4) // Lua closure
	LUA_TLGF    = LUA_TFUNCTION | (1 << 4) // light Go function
	LUA_TGCL    = LUA_TFUNCTION | (2 << 4) // Go closure
)

/* arithmetic functions */
const (
	LUA_OPADD  = iota // +
	LUA_OPSUB         // -
	LUA_OPMUL         // *
	LUA_OPMOD         // %
	LUA_OPPOW         // ^
	LUA_OPDIV         // /
	LUA_OPIDIV        // //
	LUA_OPBAND        // &
	LUA_OPBOR         // |
	LUA_OPBXOR        // ~
	LUA_OPSHL         // <<
	LUA_OPSHR         // >>
	LUA_OPUNM         // -
	LUA_OPBNOT        // ~
)

/* comparison functions */
const (
	LUA_OPEQ = iota // ==
	LUA_OPLT        // <
	LUA_OPLE        // <=
)

/* thread status */
const (
	LUA_OK = iota
	LUA_YIELD
	LUA_ERRRUN
	LUA_ERRSYNTAX
	LUA_ERRMEM
	LUA_ERRGCMM
	LUA_ERRERR
	LUA_ERRFILE
)

// todo
func LuaUpvalueIndex(i int) int {
	return LUA_REGISTRYINDEX - i
}
