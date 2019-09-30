package state

import (
	. "github.com/zxh0/lua.go/api"
)

type luaStack struct {
	/* virtual stack */
	slots []luaValue
	top   int
	/* call info */
	state   *luaState
	closure *closure
	varargs []luaValue
	openuvs map[int]*upvalue
	pc      int
	/* linked list */
	prev *luaStack
}

func newLuaStack(size int, state *luaState) *luaStack {
	return &luaStack{
		state: state,
		slots: make([]luaValue, size),
	}
}

func (stack *luaStack) check(n int) {
	free := len(stack.slots) - stack.top
	for i := free; i < n; i++ {
		stack.slots = append(stack.slots, nil)
	}
}

func (stack *luaStack) push(val luaValue) {
	if stack.top == len(stack.slots) {
		panic("stack overflow!")
	}
	stack.slots[stack.top] = val
	stack.top++
}

func (stack *luaStack) pop() luaValue {
	if stack.top < 1 {
		panic("stack underflow!")
	}
	stack.top--
	val := stack.slots[stack.top]
	stack.slots[stack.top] = nil
	return val
}

func (stack *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}

	for i := 0; i < n; i++ {
		if i < nVals {
			stack.push(vals[i])
		} else {
			stack.push(nil)
		}
	}
}

func (stack *luaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = stack.pop()
	}
	return vals
}

func (stack *luaStack) absIndex(idx int) int {
	// zero or positive or pseudo
	if idx >= 0 || idx <= LUA_REGISTRYINDEX {
		return idx
	}
	// negative
	return idx + stack.top + 1
}

func (stack *luaStack) isValid(idx int) bool {
	if idx < LUA_REGISTRYINDEX { /* upvalues */
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := stack.closure
		return c != nil && uvIdx < len(c.upvals)
	}
	if idx == LUA_REGISTRYINDEX {
		return true
	}
	absIdx := stack.absIndex(idx)
	return absIdx > 0 && absIdx <= stack.top
}

func (stack *luaStack) get(idx int) luaValue {
	if idx < LUA_REGISTRYINDEX { /* upvalues */
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := stack.closure
		if c == nil || uvIdx >= len(c.upvals) {
			return nil
		}
		return *(c.upvals[uvIdx].val)
	}

	if idx == LUA_REGISTRYINDEX {
		return stack.state.registry
	}

	absIdx := stack.absIndex(idx)
	if absIdx > 0 && absIdx <= stack.top {
		return stack.slots[absIdx-1]
	}
	return nil
}

func (stack *luaStack) set(idx int, val luaValue) {
	if idx < LUA_REGISTRYINDEX { /* upvalues */
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := stack.closure
		if c != nil && uvIdx < len(c.upvals) {
			*(c.upvals[uvIdx].val) = val
		}
		return
	}

	if idx == LUA_REGISTRYINDEX {
		stack.state.registry = val.(*luaTable)
		return
	}

	absIdx := stack.absIndex(idx)
	if absIdx > 0 && absIdx <= stack.top {
		stack.slots[absIdx-1] = val
		return
	}
	panic("invalid index!")
}

func (stack *luaStack) reverse(from, to int) {
	slots := stack.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}
