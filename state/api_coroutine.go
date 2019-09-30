package state

import (
	. "github.com/zxh0/lua.go/api"
)

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_newthread
// lua-5.3.4/src/lstate.c#lua_newthread()
func (state *luaState) NewThread() LuaState {
	t := &luaState{registry: state.registry}
	t.pushLuaStack(newLuaStack(LUA_MINSTACK, t))
	state.stack.push(t)
	return t
}

// [-?, +?, –]
// http://www.lua.org/manual/5.3/manual.html#lua_resume
func (state *luaState) Resume(from LuaState, nArgs int) ThreadStatus {
	lsFrom := from.(*luaState)
	if lsFrom.coChan == nil {
		lsFrom.coChan = make(chan int)
	}

	if state.coChan == nil {
		// start coroutine
		state.coChan = make(chan int)
		state.coCaller = lsFrom
		go func() {
			state.coStatus = state.PCall(nArgs, -1, 0)
			lsFrom.coChan <- 1
		}()
	} else {
		// resume coroutine
		if state.coStatus != LUA_YIELD { // todo
			state.stack.push("cannot resume non-suspended coroutine")
			return LUA_ERRRUN
		}
		state.coStatus = LUA_OK
		state.coChan <- 1
	}

	<-lsFrom.coChan // wait coroutine to finish or yield
	return state.coStatus
}

// [-?, +?, e]
// http://www.lua.org/manual/5.3/manual.html#lua_yield
func (state *luaState) Yield(nResults int) int {
	if state.coCaller == nil { // todo
		panic("attempt to yield from outside a coroutine")
	}
	state.coStatus = LUA_YIELD
	state.coCaller.coChan <- 1
	<-state.coChan
	return state.GetTop()
}

// [-?, +?, e]
// http://www.lua.org/manual/5.3/manual.html#lua_yieldk
func (state *luaState) YieldK() {
	panic("todo!")
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isyieldable
func (state *luaState) IsYieldable() bool {
	if state.isMainThread() {
		return false
	}
	return state.coStatus != LUA_YIELD // todo
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_status
// lua-5.3.4/src/lapi.c#lua_status()
func (state *luaState) Status() ThreadStatus {
	return state.coStatus
}
