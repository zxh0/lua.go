package state

import . "luago/api"

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_gethook
func (self *luaState) GetHook() LuaHook {
	return self.hook
}

func (self *luaState) SetHook(f LuaHook, mask, count int) {
	panic("todo: SetHook!")
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_gethookcount
func (self *luaState) GetHookCount() int {
	return 0 // todo
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_gethookmask
func (self *luaState) GetHookMask() int {
	return self.hookMask
}

func (self *luaState) GetStack(level int, ar *LuaDebug) bool {
	if level < 0 || level > self.callDepth {
		return false
	}
	if self.callDepth > 1 {
		return true
	}
	return false
	// todo
}

func (self *luaState) GetInfo(what string, ar *LuaDebug) bool {
	if len(what) > 0 && what[0] == '>' {
		what = what[1:]
		val := self.stack.pop()
		c, ok := val.(*closure)
		if !ok {
			panic("function expected")
		}

		for len(what) > 0 {
			switch what[0] {
			case 'n': // fills in the field name and namewhat;
				// ar.Name = proto.Source
				ar.NameWhat = "" // todo
			case 'S': // fills in the fields source, short_src, linedefined, lastlinedefined, and what;
				_loadFuncInfoS(ar, c)
			case 'l': // fills in the field currentline;
				ar.CurrentLine = 0 // todo
			case 't': // fills in the field istailcall;
				ar.IsTailCall = false // todo
			case 'u': // fills in the fields nups, nparams, and isvararg;
				ar.NUps = 0         // todo
				ar.NParams = 0      // todo
				ar.IsVararg = false // todo
			case 'f': // pushes onto the stack the function that is running at the given level;
				self.stack.push(val)
			case 'L': // pushes onto the stack a table whose indices are the numbers of the lines that are valid on the function.
				//panic("todo: what->L")
				self.PushNil()
			default:
				return false
			}
			what = what[1:]
		}

		return true
	}
	panic("todo: GetInfo! what=" + what)
}

// the string "Lua" if the function is a Lua function,
// "C" if it is a C function, "main" if it is the main part of a chunk.
func _loadFuncInfoS(ar *LuaDebug, c *closure) {
	if c.proto == nil {
		ar.Source = "=[C]"
		ar.LineDefined = -1
		ar.LastLineDefined = -1
		ar.What = "C"
	} else {
		p := c.proto
		if p.Source == "" {
			ar.Source = "=?"
		} else {
			ar.Source = p.Source
		}
		ar.LineDefined = int(p.LineDefined)
		ar.LastLineDefined = int(p.LastLineDefined)
		if ar.LineDefined == 0 {
			ar.What = "main"
		} else {
			ar.What = "Lua"
		}
	}
	ar.ShortSrc = _getShortSrc(ar.Source)
}

func _getShortSrc(src string) string {
	if len(src) > 0 {	
		if src[0] == '=' {
			return src[1:]
		}
	}
	return src
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
			self.stack.push(c.getUpvalue(n - 1))
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
			c.setUpvalue(n-1, self.stack.pop())
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
			return c.upvals[n-1]
		}
	}
	return nil // todo
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_upvaluejoin
func (self *luaState) UpvalueJoin(funcIdx1, n1, funcIdx2, n2 int) {
	v1 := self.stack.get(funcIdx1)
	v2 := self.stack.get(funcIdx2)
	if c1, ok := v1.(*closure); ok && len(c1.upvals) >= n1 {
		if c2, ok := v2.(*closure); ok && len(c2.upvals) >= n2 {
			c1.upvals[n1-1] = c2.upvals[n2-1]
		}
	}
}
