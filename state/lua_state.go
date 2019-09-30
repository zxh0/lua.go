package state

import (
	. "github.com/zxh0/lua.go/api"
)

type luaState struct {
	/* global state */
	hook     LuaHook
	hookMask int
	panicf   GoFunction
	registry *luaTable
	/* stack */
	stack     *luaStack
	callDepth int
	/* coroutine */
	coStatus ThreadStatus
	coCaller *luaState
	coChan   chan int
}

func New() LuaState {
	ls := &luaState{}

	registry := newLuaTable(8, 0)
	registry.put(LUA_RIDX_MAINTHREAD, ls)
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 20))

	ls.registry = registry
	ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))
	return ls
}

func (state *luaState) isMainThread() bool {
	return state.registry.get(LUA_RIDX_MAINTHREAD) == state
}

func (state *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = state.stack
	state.stack = stack
	state.callDepth++
}

func (state *luaState) popLuaStack() {
	stack := state.stack
	state.stack = stack.prev
	stack.prev = nil
	state.callDepth--
}

func (state *luaState) getLuaStack(level int) *luaStack {
	stack := state.stack
	for i := 0; i < level && stack != nil; i++ {
		stack = stack.prev
	}
	return stack
}

// debug
func (state *luaState) String() string {
	return stackToString(state.stack)
}
