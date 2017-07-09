package state

import "fmt"
import . "luago/api"

type luaStack struct {
	/* linked list */
	prev *luaStack
	/* call info */
	state *luaState
	luaCl *luaClosure
	goCl  *goClosure
	xArgs []luaValue // extraArgs
	pc    int
	/* virtual stack */
	slots []luaValue
	top   int
}

func newLuaStack(size int, state *luaState) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		state: state,
	}
}

func (self *luaStack) check(n int) bool {
	free := len(self.slots) - self.top
	if free >= n {
		return true
	}
	// grow
	slots := make([]luaValue, len(self.slots)+n+4)
	copy(slots, self.slots)
	self.slots = slots
	// never fails
	return true
}

func (self *luaStack) push(val luaValue) {
	if self.top == len(self.slots) {
		panic("stack overflow!")
	}
	self.slots[self.top] = val
	self.top++
}

func (self *luaStack) pop() luaValue {
	if self.top < 1 {
		panic("stack underflow!")
	}
	self.top--
	val := self.slots[self.top]
	self.slots[self.top] = nil
	return val
}

func (self *luaStack) pushN(vals []luaValue) {
	for _, val := range vals {
		self.push(val)
	}
}

func (self *luaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = self.pop()
	}
	return vals
}

func (self *luaStack) reverse(from, to int) {
	slots := self.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}

func (self *luaStack) absIndex(idx int) int {
	// zero or positive or pseudo
	if idx >= 0 || idx <= LUA_REGISTRYINDEX {
		return idx
	}
	// negative
	return idx + self.top + 1
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
func (self *luaStack) get(index int) luaValue {
	if index < LUA_REGISTRYINDEX {
		uvIdx := LUA_REGISTRYINDEX - index
		return self.goCl.upvals[uvIdx-1]
	}
	if index == LUA_REGISTRYINDEX {
		return self.state.registry
	}
	if absIdx := self.absIndex(index); absIdx > 0 {
		return self.slots[absIdx-1]
	}
	panic(fmt.Sprintf("bad index: %d", index))
}

func (self *luaStack) set(index int, val luaValue) {
	// todo: LUA_REGISTRYINDEX?
	if absIdx := self.absIndex(index); absIdx > 0 {
		self.slots[absIdx-1] = val
	} else {
		panic(fmt.Sprintf("bad index: %d", index))
	}
}
