package state

import . "luago/api"

type luaState struct {
	/* global state */
	panicf     GoFunction
	registry   *luaTable
	mtOfNil    *luaTable // ?
	mtOfBool   *luaTable
	mtOfNumber *luaTable
	mtOfString *luaTable
	mtOfFunc   *luaTable
	mtOfThread *luaTable
	/* stack */
	stack     *luaStack
	callDepth int
	/* coroutine */
	coStatus ThreadStatus
	coCaller *luaState
	coChan   chan int
}

// todo: rename to New()?
func New() LuaState {
	registry := newLuaTable(8, 0)
	registry.put(LUA_RIDX_MAINTHREAD, nil) // todo
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(8, 0))

	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))
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

func (self *luaState) getMetatable(val luaValue) *luaTable {
	switch x := val.(type) {
	case nil:
		return self.mtOfNil
	case bool:
		return self.mtOfBool
	case int64, float64:
		return self.mtOfNumber
	case string:
		return self.mtOfString
	case *luaClosure, *goClosure, GoFunction:
		return self.mtOfFunc
	case *luaTable:
		return x.metatable
	case *userData:
		return x.metatable
	default: // todo
		return nil
	}
}

func (self *luaState) setMetatable(val luaValue, mt *luaTable) {
	switch x := val.(type) {
	case nil:
		self.mtOfNil = mt
	case bool:
		self.mtOfBool = mt
	case int64, float64:
		self.mtOfNumber = mt
	case string:
		self.mtOfString = mt
	case *luaClosure, *goClosure, GoFunction:
		self.mtOfFunc = mt
	case *luaTable:
		x.metatable = mt
	case *userData:
		x.metatable = mt
	default:
		// todo
	}
}

func (self *luaState) getMetaField(val luaValue, fieldName string) luaValue {
	if mt := self.getMetatable(val); mt != nil {
		return mt.get(fieldName)
	} else {
		return nil
	}
}

// todo: remove this method
func (self *luaState) callMetaOp1(val luaValue, mmName string) (luaValue, bool) {
	if mm := self.getMetaField(val, mmName); mm != nil {
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
	mm := self.getMetaField(a, mmName)
	if mm == nil {
		mm = self.getMetaField(b, mmName)
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
	return stackToString(self.stack)
}
