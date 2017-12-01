package state

import . "luago/api"

func (self *luaState) GetHook() LuaHook {
	panic("todo: GetHook!")
}

func (self *luaState) SetHook(f LuaHook, mask, count int) {
	panic("todo: SetHook!")
}

func (self *luaState) GetHookCount() int {
	panic("todo: GetHookCount!")
}

func (self *luaState) GetHookMask() int {
	panic("todo: GetHookMask!")
}

func (self *luaState) GetStack(level int, ar *LuaDebug) bool {
	//panic("todo: GetStack!")
	// todo
	if self.callDepth > 1 {
		return true
	}
	return false
}

func (self *luaState) GetInfo(what string, ar *LuaDebug) bool {
	panic("todo: GetInfo!")
}

func (self *luaState) GetLocal(ar *LuaDebug, n int) string {
	panic("todo: GetLocal!")
}

func (self *luaState) SetLocal(ar *LuaDebug, n int) string {
	panic("todo: SetLocal!")
}

// [-0, +(0|1), –]
// http://www.lua.org/manual/5.3/manual.html#lua_getupvalue
func (self *luaState) GetUpvalue(funcIdx, n int) string {
	val := self.stack.get(funcIdx)
	if c, ok := val.(*closure); ok {
		if len(c.upvals) >= n {
			uv := *(c.upvals[n-1])
			self.stack.push(uv)
			return c.getUpvalueName(n - 1)
		}
	}
	return ""
}

// [-(0|1), +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_setupvalue
func (self *luaState) SetUpvalue(funcIdx, n int) string {
	val := self.stack.get(funcIdx)
	if c, ok := val.(*closure); ok {
		if len(c.upvals) >= n {
			*(c.upvals[n-1]) = self.stack.pop()
			return c.getUpvalueName(n - 1)
		}
	}
	return ""
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_upvalueid
func (self *luaState) UpvalueId(funcIdx, n int) interface{} {
	val := self.stack.get(funcIdx)
	if c, ok := val.(*closure); ok {
		if len(c.upvals) >= n {
			return *(c.upvals[n-1])
		}
	}
	return nil // todo
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_upvaluejoin
func (self *luaState) UpvalueJoin(funcIdx1, n1, funcIdx2, n2 int) {
	panic("todo: UpvalueJoin!")
}
