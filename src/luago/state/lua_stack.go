package state

import "fmt"

type luaStack struct {
	prev  *luaStack
	luaCl *luaClosure // todo: move to luaState?
	goCl  *goClosure  // todo: move to luaState?
	xArgs []luaValue  // extraArgs
	slots []luaValue  // registers+stack
	top   int         // stack pointer
	pc    int         // todo: move to somewhere?
}

func newLuaStack(nSlots, nRegs int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, nSlots),
		top:   nRegs,
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
		panic(fmt.Sprintf("stack overflow! top=%d", self.top))
	}
	self.slots[self.top] = val
	self.top++
}

func (self *luaStack) pop() luaValue {
	if self.top < 1 {
		panic(fmt.Sprintf("stack underflow! top=%d", self.top))
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
	if self.top < n {
		panic(fmt.Sprintf("stack underflow! n=%d", n))
	}
	self.top -= n
	return self.slots[self.top : self.top+n]
}

func (self *luaStack) reverse(from, to int) {
	slots := self.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}
