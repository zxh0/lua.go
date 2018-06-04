package stdlib

import . "luago/api"

var coFuncs = map[string]GoFunction{
	"create":      coCreate,
	"resume":      coResume,
	"yield":       coYield,
	"status":      coStatus,
	"isyieldable": coYieldable,
	"running":     coRunning,
	"wrap":        coWrap,
}

func OpenCoroutineLib(L LuaState) int {
	luaL_newlib(L, coFuncs)
	return 1
}

// coroutine.create (f)
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.create
// lua-5.3.4/src/lcorolib.c#luaB_cocreate()
func coCreate(L LuaState) int {
	luaL_checktype(L, 1, LUA_TFUNCTION)
	NL := lua_newthread(L)
	lua_pushvalue(L, 1) /* move function to top */
	lua_xmove(L, NL, 1) /* move function from L to NL */
	return 1
}

// coroutine.resume (co [, val1, ···])
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.resume
// lua-5.3.4/src/lcorolib.c#luaB_coresume()
func coResume(L LuaState) int {
	co := getco(L)
	r := auxresume(L, co, lua_gettop(L)-1)
	if r < 0 {
		lua_pushboolean(L, false)
		lua_insert(L, -2)
		return 2 /* return false + error message */
	} else {
		lua_pushboolean(L, true)
		lua_insert(L, -(r + 1))
		return r + 1 /* return true + 'resume' returns */
	}
}

func getco(L LuaState) LuaState {
	co := lua_tothread(L, 1)
	luaL_argcheck(L, co != nil, 1, "thread expected")
	return co
}

func auxresume(L, co LuaState, narg int) int {
	if !lua_checkstack(co, narg) {
		lua_pushliteral(L, "too many arguments to resume")
		return -1 /* error flag */
	}
	if lua_status(co) == LUA_OK && lua_gettop(co) == 0 {
		lua_pushliteral(L, "cannot resume dead coroutine")
		return -1 /* error flag */
	}
	lua_xmove(L, co, narg)
	status := lua_resume(co, L, narg)
	if status == LUA_OK || status == LUA_YIELD {
		nres := lua_gettop(co)
		if !lua_checkstack(L, nres+1) {
			lua_pop(co, nres) /* remove results anyway */
			lua_pushliteral(L, "too many results to resume")
			return -1 /* error flag */
		}
		lua_xmove(co, L, nres) /* move yielded values */
		return nres
	} else {
		lua_xmove(co, L, 1) /* move error message */
		return -1           /* error flag */
	}
}

// coroutine.yield (···)
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.yield
// lua-5.3.4/src/lcorolib.c#luaB_yield()
func coYield(L LuaState) int {
	return lua_yield(L, lua_gettop(L))
}

// coroutine.status (co)
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.status
// lua-5.3.4/src/lcorolib.c#luaB_costatus()
func coStatus(L LuaState) int {
	co := getco(L)
	if L == co {
		lua_pushliteral(L, "running")
	} else {
		switch lua_status(co) {
		case LUA_YIELD:
			lua_pushliteral(L, "suspended")
		case LUA_OK:
			ar := LuaDebug{}
			if lua_getstack(co, 0, &ar) { /* does it have frames? */
				lua_pushliteral(L, "normal") /* it is running */
			} else if lua_gettop(co) == 0 {
				lua_pushliteral(L, "dead")
			} else {
				lua_pushliteral(L, "suspended") /* initial state */
			}
		default: /* some error occurred */
			lua_pushliteral(L, "dead")
			break
		}
	}
	return 1
}

// coroutine.isyieldable ()
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.isyieldable
// lua-5.3.4/src/lcorolib.c#luaB_yieldable()
func coYieldable(L LuaState) int {
	lua_pushboolean(L, lua_isyieldable(L))
	return 1
}

// coroutine.running ()
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.running
func coRunning(L LuaState) int {
	ismain := lua_pushthread(L)
	lua_pushboolean(L, ismain)
	return 2
}

// coroutine.wrap (f)
// http://www.lua.org/manual/5.3/manual.html#pdf-coroutine.wrap
func coWrap(ls LuaState) int {
	panic("todo: coWrap!")
}
