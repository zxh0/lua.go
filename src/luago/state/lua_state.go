package state

import . "luago/api"

type luaState struct {
	/* global state */
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

func (self *luaState) isMainThread() bool {
	return self.registry.get(LUA_RIDX_MAINTHREAD) == self
}

func (self *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = self.stack
	self.stack = stack
	self.callDepth++
}

func (self *luaState) popLuaStack() {
	stack := self.stack
	self.stack = stack.prev
	stack.prev = nil
	self.callDepth--
}

// debug
func (self *luaState) String() string {
	return stackToString(self.stack)
}
