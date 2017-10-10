package state

import . "luago/api"

type luaStack struct {
	/* linked list */
	prev *luaStack
	/* call info */
	state *luaState
	luaCl *luaClosure
	goCl  *goClosure
	vargs []luaValue // varargs
	pc    int
	/* virtual stack */
	slots []luaValue
	top   int
}

func newLuaStack(size int, state *luaState) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
		state: state,
	}
}

func (self *luaStack) check(n int) {
	free := len(self.slots) - self.top
	for i := free; i < n; i++ {
		self.slots = append(self.slots, nil)
	}
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

func (self *luaStack) absIndex(idx int) int {
	// zero or positive or pseudo
	if idx >= 0 || idx <= LUA_REGISTRYINDEX {
		return idx
	}
	// negative
	return idx + self.top + 1
}

func (self *luaStack) isValid(idx int) bool {
	if idx < LUA_REGISTRYINDEX { /* upvalues */
		uvIdx := LUA_REGISTRYINDEX - idx
		return self.goCl != nil && uvIdx <= len(self.goCl.upvals)
	}
	if idx == LUA_REGISTRYINDEX {
		return true
	}
	absIdx := self.absIndex(idx)
	return absIdx > 0 || absIdx <= self.top
}

func (self *luaStack) get(idx int) luaValue {
	if idx < LUA_REGISTRYINDEX { /* upvalues */
		uvIdx := LUA_REGISTRYINDEX - idx
		//if uvIdx > MAXUPVAL + 1 {
		//	panic("upvalue index too large!")
		//}
		if self.goCl == nil || len(self.goCl.upvals) < uvIdx {
			return nil
		}
		return self.goCl.upvals[uvIdx-1]
	}

	if idx == LUA_REGISTRYINDEX {
		return self.state.registry
	}

	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		return self.slots[absIdx-1]
	}
	return nil
}

func (self *luaStack) set(idx int, val luaValue) {
	// todo: LUA_REGISTRYINDEX?
	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		self.slots[absIdx-1] = val
		return
	}
	panic("todo!")
}

func (self *luaStack) reverse(from, to int) {
	slots := self.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}
