package state

import . "luago/lua"

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_absindex
func (self *luaState) AbsIndex(idx int) int {
	if idx > 0 || IsPseudo(idx) {
		return idx
	}
	return self.stack.absIndex(idx)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_checkstack
func (self *luaState) CheckStack(n int) bool {
	return self.stack.check(n)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_gettop
func (self *luaState) GetTop() int {
	return self.stack.sp
}

// [-?, +?, –]
// http://www.lua.org/manual/5.3/manual.html#lua_settop
func (self *luaState) SetTop(index int) {
	if index < 0 {
		index = self.stack.absIndex(index)
	}

	n := self.stack.sp - index
	if n > 0 {
		for i := 0; i < n; i++ {
			self.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			self.stack.push(nil)
		}
	}
}

// [-n, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pop
func (self *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		self.stack.pop()
	}
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushvalue
func (self *luaState) PushValue(index int) {
	val := self.stack.get(index)
	self.stack.push(val)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_copy
func (self *luaState) Copy(fromIdx, toIdx int) {
	val := self.stack.get(fromIdx)
	self.stack.set(toIdx, val)
}

// [-1, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_insert
func (self *luaState) Insert(index int) {
	self.Rotate(index, 1)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_remove
func (self *luaState) Remove(index int) {
	self.Rotate(index, -1)
	self.Pop(1)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_replace
func (self *luaState) Replace(index int) {
	self.Copy(-1, index)
	self.Pop(1)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rotate
func (self *luaState) Rotate(idx, n int) {
	stack := self.stack
	t := stack.sp - 1        /* end of stack segment being rotated */
	p := stack.absIndex(idx) /* start of segment */
	p -= 1
	var m int /* end of prefix */
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	stack.reverse(p, m)   /* reverse the prefix with length 'n' */
	stack.reverse(m+1, t) /* reverse the suffix */
	stack.reverse(p, t)   /* reverse the entire segment */
}

// [-?, +?, –]
// http://www.lua.org/manual/5.3/manual.html#lua_xmove
func (self *luaState) XMove(to LuaState, n int) {
	fromLs := self
	toLs := to.(*luaState)

	x := fromLs.stack.popN(n)
	toLs.stack.pushN(x)
}
