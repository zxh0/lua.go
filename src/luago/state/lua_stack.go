package state

import "bytes"
import "fmt"
import . "luago/api"

/*
 sp->[ ] -.
     [ ]  |- stack
 bp->[ ] -'
     [ ] -.
     [ ]  |- registers
     [ ] -'
*/
type luaStack struct {
	prev  *luaStack
	state *luaState   // todo: remove?
	luaCl *luaClosure // todo: move to luaState?
	goCl  *goClosure  // todo: move to luaState?
	xArgs []luaValue  // extraArgs
	slots []luaValue  // registers+stack
	bp    int         // stack base pointer
	sp    int         // stack pointer
	pc    int         // todo: move to somewhere?
}

func newLuaStack(nSlots, nRegs int, state *luaState) *luaStack {
	return &luaStack{
		state: state,
		slots: make([]luaValue, nSlots),
		sp:    nRegs,
		bp:    nRegs,
	}
}

func (self *luaStack) check(n int) bool {
	free := len(self.slots) - self.sp
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

func (self *luaStack) absIndex(idx int) int {
	if idx > 0 && idx <= self.sp {
		return idx
	}
	if idx < 0 && idx >= -self.sp {
		return idx + self.sp + 1
	}
	return 0
}

/* registers */

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

func (self *luaStack) reverse(from, to int) {
	slots := self.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}

/* stack */

func (self *luaStack) push(val luaValue) {
	if self.sp == len(self.slots) {
		panic(fmt.Sprintf("stack overflow! sp=%d", self.sp))
	}
	self.slots[self.sp] = val
	self.sp++
}

func (self *luaStack) pop() luaValue {
	if self.sp-self.bp < 1 {
		panic(fmt.Sprintf("stack underflow! sp=%d", self.sp))
	}
	self.sp--
	val := self.slots[self.sp]
	self.slots[self.sp] = nil
	return val
}

func (self *luaStack) pushN(vals []luaValue) {
	for _, val := range vals {
		self.push(val)
	}
}

func (self *luaStack) popN(n int) []luaValue {
	if self.sp-self.bp < n {
		panic(fmt.Sprintf("stack underflow! n=%d", n))
	}
	self.sp -= n
	return self.slots[self.sp : self.sp+n]
}

func (self *luaStack) popAll() []luaValue {
	return self.popN(self.sp - self.bp)
}

func (self *luaStack) popLuaGoFunction() LuaGoFunction {
	val := self.pop()
	if f, ok := val.(LuaGoFunction); ok {
		return f
	}
	panic("not LuaGoFunction!")
}

/* debug */

func (self *luaStack) toString() string {
	var buf bytes.Buffer

	for i := 0; i < self.sp; i++ {
		if i == self.bp {
			buf.WriteString("~")
		}
		buf.WriteString("[")
		buf.WriteString(valToString(self.slots[i]))
		buf.WriteString("]")
	}

	return buf.String()
}
