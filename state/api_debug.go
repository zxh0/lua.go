package state

import "fmt"
import "reflect"
import "runtime"
import "strings"
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

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_getstack
func (self *luaState) GetStack(level int, ar *LuaDebug) bool {
	if level < 0 || level >= self.callDepth-1 {
		return false
	}
	if self.callDepth > 1 {
		if s := self.getLuaStack(level); s != nil {
			ar.CallInfo = s
			return true
		}
	}
	return false
}

// [-(0|1), +(0|1|2), e]
// http://www.lua.org/manual/5.3/manual.html#lua_getinfo
func (self *luaState) GetInfo(what string, ar *LuaDebug) bool {
	if len(what) > 0 && what[0] == '>' {
		what = what[1:]
		val := self.stack.pop()
		if c, ok := val.(*closure); ok {
			return self.loadInfo(ar, c, what)
		}
		panic("function expected")
	}

	if ci := ar.CallInfo; ci != nil {
		if c := ci.(*luaStack).closure; c != nil {
			return self.loadInfo(ar, c, what)
		}
	}

	return false
}

func (self *luaState) loadInfo(ar *LuaDebug, c *closure, what string) bool {
	for len(what) > 0 {
		switch what[0] {
		case 'n': // fills in the field name and namewhat;
			ar.Name = _getFuncName(c)
			ar.NameWhat = "" // todo
		case 'S': // fills in the fields source, short_src, linedefined, lastlinedefined, and what;
			_setFuncInfoS(ar, c)
		case 'l': // fills in the field currentline;
			_setCurrentLine(ar, c)
		case 't': // fills in the field istailcall;
			ar.IsTailCall = false // todo
		case 'u': // fills in the fields nups, nparams, and isvararg;
			ar.NUps = 0         // todo
			ar.NParams = 0      // todo
			ar.IsVararg = false // todo
		case 'f': // pushes onto the stack the function that is running at the given level;
			self.stack.push(c)
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

func _getFuncName(c *closure) string {
	if gof := c.goFunc; gof != nil {
		pc := reflect.ValueOf(gof).Pointer()
		if f := runtime.FuncForPC(pc); f != nil {
			name := f.Name()
			if strings.HasPrefix(name, "luago/stdlib.") {
				name = name[13:]                                        // remove "luago/stdlib."
				for len(name) > 0 && name[0] >= 'a' && name[0] <= 'z' { // remove prefix
					name = name[1:]
				}
				return strings.ToLower(name)
			}
		}
	}
	return "?"
}

// the string "Lua" if the function is a Lua function,
// "C" if it is a C function, "main" if it is the main part of a chunk.
func _setFuncInfoS(ar *LuaDebug, c *closure) {
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
	if len(src) > 0 { /* 'literal' source */
		if src[0] == '=' {
			src = src[1:]
			if strLen := len(src); strLen > LUA_IDSIZE {
				src = src[:LUA_IDSIZE-1]
			}
		} else if src[0] == '@' { /* file name */
			src = src[1:]
			if strLen := len(src); strLen > LUA_IDSIZE {
				src = "..." + src[strLen-LUA_IDSIZE+4:]
			}
		} else { /* string; format as [string "source"] */
			if i := strings.IndexByte(src, '\n'); i >= 0 {
				src = src[0:i] + "..."
			}
			maxSrcLen := LUA_IDSIZE - len(`[string " "]`)
			if len(src) > maxSrcLen {
				src = src[0:maxSrcLen-3] + "..."
			}
			src = fmt.Sprintf(`[string "%s"]`, src)
		}
	}
	return src
}

func _setCurrentLine(ar *LuaDebug, c *closure) {
	if ci := ar.CallInfo; ci != nil {
		c := ci.(*luaStack).closure
		pc := ci.(*luaStack).pc
		if c.proto == nil || pc < 1 || pc > len(c.proto.LineInfo) {
			ar.CurrentLine = -1
		} else {
			ar.CurrentLine = int(c.proto.LineInfo[pc-1])
		}
	}
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
