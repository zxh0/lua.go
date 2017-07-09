package state

import "fmt"
import . "luago/api"

type callFrame struct {
	// todo
}

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
	/* virtual stack */
	stack     *luaStack // callStack
	callDepth int // todo: rename?
	/* coroutine */
	status ThreadStatus
}

// todo: rename to New()?
func NewLuaState() LuaState {
	registry := newLuaTable(8, 0)
	registry.put(LUA_RIDX_MAINTHREAD, nil) // todo
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(8, 0))

	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(16, 0))
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

func (self *luaState) absIndex(idx int) int {
	if idx > 0 || idx <= LUA_REGISTRYINDEX {
		return idx
	}
	if idx < 0 && idx >= -self.stack.top {
		return idx + self.stack.top + 1
	}
	return 0 // todo
}

// func (self *luaStack) _get(index int) (luaValue, bool) {
// 	if index < LUA_REGISTRYINDEX { /* upvalues */
// 		uvIdx := LUA_REGISTRYINDEX - index
// 		if uvIdx > MAXUPVAL + 1 {
// 			panic("upvalue index too large!")
// 		} else if self.goCl == nil || len(self.goCl.upvals) < uvIdx {
// 			return nil, false
// 		} else {
// 			return self.goCl.upvals[uvIdx-1], true
// 		}
// 	} else if index == LUA_REGISTRYINDEX {
// 		return self.state.registry, true
// 	} else {
// 		absIdx := self.absIndex(index)
// 		if absIdx <= 0 || absIdx > len(self.slots) {
// 			return nil, false
// 		} else {
// 			return self.slots[absIdx-1], true
// 		}
// 	}
// }

// func (self *luaStack) getOrNil(index int) luaValue {
// 	if val, ok := self._get(index); ok {
// 		return val
// 	} else {
// 		return nil
// 	}
// }

// todo: move to luaState
func (self *luaState) get(index int) luaValue {
	if index < LUA_REGISTRYINDEX {
		uvIdx := LUA_REGISTRYINDEX - index
		return self.stack.goCl.upvals[uvIdx-1]
	}
	if index == LUA_REGISTRYINDEX {
		return self.registry
	}
	if absIdx := self.absIndex(index); absIdx > 0 {
		return self.stack.slots[absIdx-1]
	}
	panic(fmt.Sprintf("bad index: %d", index))
}

func (self *luaState) set(index int, val luaValue) {
	// todo: LUA_REGISTRYINDEX?
	if absIdx := self.absIndex(index); absIdx > 0 {
		self.stack.slots[absIdx-1] = val
	} else {
		panic(fmt.Sprintf("bad index: %d", index))
	}
}

func (self *luaState) getMetaTable(val luaValue) *luaTable {
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
		return x.metaTable
	case *userData:
		return x.metaTable
	default: // todo
		return nil
	}
}

func (self *luaState) setMetaTable(val luaValue, mt *luaTable) {
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
		x.metaTable = mt
	case *userData:
		x.metaTable = mt
	default:
		// todo
	}
}

func (self *luaState) getMetaField(val luaValue, fieldName string) luaValue {
	if mt := self.getMetaTable(val); mt != nil {
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
