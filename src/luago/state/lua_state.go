package state

import . "luago/api"

/*
 * luaState
 *   panicf
 *   global
 *   registry
 *   luaStack <-.
 *     prev ----'
 *     slots
 *   callDepth
 */
type luaState struct {
	panicf    GoFunction
	registry  *luaTable
	stack     *luaStack // virtual stack
	callDepth int       // todo: rename
}

// todo: rename to New()?
func NewLuaState() LuaState {
	registry := newLuaTable(8, 0)
	registry.put(LUA_RIDX_MAINTHREAD, nil) // todo
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(8, 0))

	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(16, 0, ls))
	return ls
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

// todo: remove this method
func (self *luaState) callMetaOp1(val luaValue, mmName string) (luaValue, bool) {
	if mm := getMetaField(val, mmName); mm != nil {
		self.stack.check(4)
		self.stack.push(mm)
		self.stack.push(val)
		self.Call(1, 1)
		return self.stack.pop(), true
	} else {
		return nil, false
	}
}

func (self *luaState) callMetaOp2(a, b luaValue, mmName string) (luaValue, bool) {
	mm := getMetaField(a, mmName)
	if mm == nil {
		mm = getMetaField(b, mmName)
	}

	if mm != nil {
		self.stack.check(4)
		self.stack.push(mm)
		self.stack.push(a)
		self.stack.push(b)
		self.Call(2, 1)
		return self.stack.pop(), true
	} else {
		return nil, false
	}
}

// debug
func (self *luaState) String() string {
	return self.stack.toString()
}
