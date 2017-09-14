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

/* option for multiple returns in 'lua_pcall' and 'lua_call' */
const LUA_MULTRET = -1

/* minimum Lua stack available to a C function */
const LUA_MINSTACK = 20

/*
@@ LUAI_MAXSTACK limits the size of the Lua stack.
** CHANGE it if you need a different limit. This limit is arbitrary;
** its only purpose is to stop Lua from consuming unlimited stack
** space (and to reserve some numbers for pseudo-indices).
*/
const LUAI_MAXSTACK = 1000000

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

const (
	LUA_MAXINTEGER = 1<<63 - 1
	LUA_MININTEGER = -1 << 63
)

// lua-5.3.4/src/lua.h
/* basic types */
const (
	LUA_TNONE          = -1
	LUA_TNIL           = 0
	LUA_TBOOLEAN       = 1
	LUA_TLIGHTUSERDATA = 2
	LUA_TNUMBER        = 3
	LUA_TSTRING        = 4
	LUA_TTABLE         = 5
	LUA_TFUNCTION      = 6
	LUA_TUSERDATA      = 7
	LUA_TTHREAD        = 8
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

// lua-5.3.4/src/lua.h
/* arithmetic functions */
const (
	LUA_OPADD  = 0  // +
	LUA_OPSUB  = 1  // -
	LUA_OPMUL  = 2  // *
	LUA_OPMOD  = 3  // %
	LUA_OPPOW  = 4  // ^
	LUA_OPDIV  = 5  // /
	LUA_OPIDIV = 6  // //
	LUA_OPBAND = 7  // &
	LUA_OPBOR  = 8  // |
	LUA_OPBXOR = 9  // ~
	LUA_OPSHL  = 10 // <<
	LUA_OPSHR  = 11 // >>
	LUA_OPUNM  = 12 // -
	LUA_OPBNOT = 13 // ~
)

// lua-5.3.4/src/lua.h
/* comparison functions */
const (
	LUA_OPEQ = 0 // ==
	LUA_OPLT = 1 // <
	LUA_OPLE = 2 // <=
)

// lua-5.3.4/src/lua.h
/* thread status */
const (
	LUA_OK        = 0
	LUA_YIELD     = 1
	LUA_ERRRUN    = 2
	LUA_ERRSYNTAX = 3
	LUA_ERRMEM    = 4
	LUA_ERRGCMM   = 5
	LUA_ERRERR    = 6
	LUA_ERRFILE   = 7
)
