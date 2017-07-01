package api

/* basic types */
const (
	LUA_TNONE    = -1
	LUA_TNIL     = 0
	LUA_TBOOLEAN = 1
	// LUA_TLIGHTUSERDATA = 2
	LUA_TNUMBER   = 3
	LUA_TSTRING   = 4
	LUA_TTABLE    = 5
	LUA_TFUNCTION = 6
	LUA_TUSERDATA = 7
	LUA_TTHREAD   = 8
)

/* type variants */
const (
	/* Variant tags for numbers */
	LUA_TNUMFLT = LUA_TNUMBER | (0 << 4) /* float numbers */
	LUA_TNUMINT = LUA_TNUMBER | (1 << 4) /* integer numbers */
	/* Variant tags for strings */
	LUA_TSHRSTR = LUA_TSTRING | (0 << 4) /* short strings */
	LUA_TLNGSTR = LUA_TSTRING | (1 << 4) /* long strings */
	/* Variant tags for functions */
	LUA_TLCL = LUA_TFUNCTION | (0 << 4) /* Lua closure */
	LUA_TLGF = LUA_TFUNCTION | (1 << 4) /* light Go function */
	LUA_TGCL = LUA_TFUNCTION | (2 << 4) /* Go closure */
)

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

/* comparison functions */
const (
	LUA_OPEQ = 0 // ==
	LUA_OPLT = 1 // <
	LUA_OPLE = 2 // <=
)

/* thread status */
const (
	LUA_OK        = 0
	LUA_YIELD     = 1
	LUA_ERRRUN    = 2
	LUA_ERRSYNTAX = 3
	LUA_ERRMEM    = 4
	LUA_ERRGCMM   = 5
	LUA_ERRERR    = 6
)
