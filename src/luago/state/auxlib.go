package state

import "fmt"
import "io/ioutil"
import "strings"
import . "luago/api"
import "luago/stdlib"

/* key, in the registry, for table of loaded modules */
const LUA_LOADED_TABLE = "_LOADED"

const LEVELS1 = 10 /* size of the first part of the stack */
const LEVELS2 = 11 /* size of the second part of the stack */

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_error
func (self *luaState) Error2(fmt string, a ...interface{}) int {
	self.PushFString(fmt, a...) // todo
	return self.Error()
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_argerror
func (self *luaState) ArgError(arg int, extraMsg string) int {
	// bad argument #arg to 'funcname' (extramsg)
	return self.Error2("bad argument #%d (%s)", arg, extraMsg) // todo
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_where
// lua-5.3.4/src/lauxlib.c#luaL_where()
func (self *luaState) Where(level int) {
	// chunkname:currentline:
	ar := LuaDebug{}
	if self.GetStack(level, &ar) { /* check function at level */
		self.GetInfo("Sl", &ar) /* get info about it */
		if ar.CurrentLine > 0 { /* is there info? */
			self.PushFString("%s:%d: ", ar.ShortSrc, ar.CurrentLine)
			return
		}
	}
	self.PushFString("") /* else, no information available... */
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkstack
// lua-5.3.4/src/lauxlib.c#luaL_checkstack()
func (self *luaState) CheckStack2(sz int, msg string) {
	if !self.CheckStack(sz) {
		if msg != "" {
			self.Error2("stack overflow (%s)", msg)
		} else {
			self.Error2("stack overflow")
		}
	}
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_argcheck
// lua-5.3.4/src/lauxlib.c#luaL_argcheck()
func (self *luaState) ArgCheck(cond bool, arg int, extraMsg string) {
	if !cond {
		self.ArgError(arg, extraMsg)
	}
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkany
// lua-5.3.4/src/lauxlib.c#luaL_checkany()
func (self *luaState) CheckAny(arg int) {
	if self.Type(arg) == LUA_TNONE {
		self.ArgError(arg, "value expected")
	}
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checktype
// lua-5.3.4/src/lauxlib.c#luaL_checktype()
func (self *luaState) CheckType(arg int, t LuaType) {
	if self.Type(arg) != t {
		self.tagError(arg, t)
	}
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkinteger
// lua-5.3.4/src/lauxlib.c#luaL_checkinteger()
func (self *luaState) CheckInteger(arg int) int64 {
	i, ok := self.ToIntegerX(arg)
	if !ok {
		self.intError(arg)
	}
	return i
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checknumber
// lua-5.3.4/src/lauxlib.c#luaL_checknumber()
func (self *luaState) CheckNumber(arg int) float64 {
	f, ok := self.ToNumberX(arg)
	if !ok {
		self.tagError(arg, LUA_TNUMBER)
	}
	return f
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkstring
// http://www.lua.org/manual/5.3/manual.html#luaL_checklstring
// lua-5.3.4/src/lauxlib.c#luaL_checklstring()
func (self *luaState) CheckString(arg int) string {
	s, ok := self.ToStringX(arg)
	if !ok {
		self.tagError(arg, LUA_TSTRING)
	}
	return s
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_optinteger
// lua-5.3.4/src/lauxlib.c#luaL_optinteger()
func (self *luaState) OptInteger(arg int, def int64) int64 {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckInteger(arg)
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_optnumber
// lua-5.3.4/src/lauxlib.c#luaL_optnumber()
func (self *luaState) OptNumber(arg int, def float64) float64 {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckNumber(arg)
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_optstring
// lua-5.3.4/src/lauxlib.c#luaL_optlstring()
func (self *luaState) OptString(arg int, def string) string {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckString(arg)
}

// [-0, +?, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_dofile
// lua-5.3.4/src/lauxlib.h#luaL_dofile()
func (self *luaState) DoFile(filename string) bool {
	return self.LoadFile(filename) == LUA_OK &&
		self.PCall(0, LUA_MULTRET, 0) == LUA_OK
}

// [-0, +?, –]
// http://www.lua.org/manual/5.3/manual.html#luaL_dostring
// lua-5.3.4/src/lauxlib.h#luaL_dostring()
func (self *luaState) DoString(str string) bool {
	return self.LoadString(str) == LUA_OK &&
		self.PCall(0, LUA_MULTRET, 0) == LUA_OK
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_loadfile
// lua-5.3.4/src/lauxlib.h#luaL_loadfile()
func (self *luaState) LoadFile(filename string) ThreadStatus {
	return self.LoadFileX(filename, "bt")
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_loadfilex
func (self *luaState) LoadFileX(filename, mode string) ThreadStatus {
	if data, err := ioutil.ReadFile(filename); err == nil {
		return self.Load(data, "@"+filename, mode)
	}
	return LUA_ERRFILE // todo
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#luaL_loadstring
func (self *luaState) LoadString(s string) ThreadStatus {
	return self.Load([]byte(s), s, "bt")
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkversion
func (self *luaState) CheckVersion() {
	//panic("todo: CheckVersion!")
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#luaL_typename
// lua-5.3.4/src/lauxlib.h#luaL_typename()
func (self *luaState) TypeName2(idx int) string {
	return self.TypeName(self.Type(idx))
}

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_len
// lua-5.3.4/src/lauxlib.c#luaL_len()
func (self *luaState) Len2(idx int) int64 {
	self.Len(idx)
	i, isNum := self.ToIntegerX(-1)
	if !isNum {
		self.Error2("object length is not an integer")
	}
	self.Pop(1)
	return i
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_tolstring
// lua-5.3.4/src/lauxlib.c#luaL_tolstring()
func (self *luaState) ToString2(idx int) string {
	if self.CallMeta(idx, "__tostring") { /* metafield? */
		if !self.IsString(-1) {
			self.Error2("'__tostring' must return a string")
		}
	} else {
		switch self.Type(idx) {
		case LUA_TNUMBER:
			if self.IsInteger(idx) {
				self.PushString(fmt.Sprintf("%d", self.ToInteger(idx))) // todo
			} else {
				self.PushString(fmt.Sprintf("%g", self.ToNumber(idx))) // todo
			}
		case LUA_TSTRING:
			self.PushValue(idx)
		case LUA_TBOOLEAN:
			if self.ToBoolean(idx) {
				self.PushString("true")
			} else {
				self.PushString("false")
			}
		case LUA_TNIL:
			self.PushString("nil")
		default:
			tt := self.GetMetafield(idx, "__name") /* try name */
			var kind string
			if tt == LUA_TSTRING {
				kind = self.CheckString(-1)
			} else {
				kind = self.TypeName2(idx)
			}

			self.PushString(fmt.Sprintf("%s: %p", kind, self.ToPointer(idx)))
			if tt != LUA_TNIL {
				self.Remove(-2) /* remove '__name' */
			}
		}
	}
	return self.CheckString(-1)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_getsubtable
// lua-5.3.4/src/lauxlib.c#luaL_getsubtable()
func (self *luaState) GetSubTable(idx int, fname string) bool {
	if self.GetField(idx, fname) == LUA_TTABLE {
		return true /* table already there */
	}
	self.Pop(1) /* remove previous result */
	idx = self.stack.absIndex(idx)
	self.NewTable()
	self.PushValue(-1)        /* copy to be left at top */
	self.SetField(idx, fname) /* assign new table to field */
	return false              /* false, because did not find table there */
}

// [-0, +(0|1), m]
// http://www.lua.org/manual/5.3/manual.html#luaL_getmetafield
// lua-5.3.4/src/lauxlib.c#luaL_getmetafield()
func (self *luaState) GetMetafield(obj int, event string) LuaType {
	if !self.GetMetatable(obj) { /* no metatable? */
		return LUA_TNIL
	}

	self.PushString(event)
	tt := self.RawGet(-2)
	if tt == LUA_TNIL { /* is metafield nil? */
		self.Pop(2) /* remove metatable and metafield */
	} else {
		self.Remove(-2) /* remove only metatable */
	}
	return tt /* return metafield type */
}

// [-0, +(0|1), e]
// http://www.lua.org/manual/5.3/manual.html#luaL_callmeta
// lua-5.3.4/src/lauxlib.c#luaL_callmeta()
func (self *luaState) CallMeta(obj int, event string) bool {
	obj = self.AbsIndex(obj)
	if self.GetMetafield(obj, event) == LUA_TNIL { /* no metafield? */
		return false
	}

	self.PushValue(obj)
	self.Call(1, 1)
	return true
}

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_openlibs
// lua-5.3.4/src/linit.c#luaL_openlibs()
func (self *luaState) OpenLibs() {
	libs := map[string]GoFunction{
		"_G":        stdlib.OpenBaseLib,
		"package":   stdlib.OpenPackageLib,
		"coroutine": stdlib.OpenCoroutineLib,
		"table":     stdlib.OpenTableLib,
		"io":        stdlib.OpenIOLib,
		"os":        stdlib.OpenOSLib,
		"string":    stdlib.OpenStringLib,
		"math":      stdlib.OpenMathLib,
		"utf8":      stdlib.OpenUTF8Lib,
		"debug":     stdlib.OpenDebugLib,
	}

	for name, fun := range libs {
		self.RequireF(name, fun, true)
		self.Pop(1)
	}
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_requiref
// lua-5.3.4/src/lauxlib.c#luaL_requiref()
func (self *luaState) RequireF(modname string, openf GoFunction, glb bool) {
	self.GetSubTable(LUA_REGISTRYINDEX, "_LOADED")
	self.GetField(-1, modname)
	if !self.ToBoolean(-1) { /* package not already loaded? */
		self.Pop(1) /* remove field */
		self.PushGoFunction(openf)
		self.PushString(modname)   /* argument to open function */
		self.Call(1, 1)            /* call 'openf' to open module */
		self.PushValue(-1)         /* make copy of module (call result) */
		self.SetField(-3, modname) /* _LOADED[modname] = module */
	}
	self.Remove(-2) /* remove _LOADED table */
	if glb {
		self.PushValue(-1)      /* copy of module */
		self.SetGlobal(modname) /* _G[modname] = module */
	}
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_newlib
// lua-5.3.4/src/lauxlib.h#luaL_newlib()
func (self *luaState) NewLib(l FuncReg) {
	self.NewLibTable(l)
	self.SetFuncs(l, 0)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_newlibtable
// lua-5.3.4/src/lauxlib.h#luaL_newlibtable()
func (self *luaState) NewLibTable(l FuncReg) {
	self.CreateTable(0, len(l))
}

// [-nup, +0, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_setfuncs
// lua-5.3.4/src/lauxlib.c#luaL_setfuncs()
func (self *luaState) SetFuncs(l FuncReg, nup int) {
	self.CheckStack2(nup, "too many upvalues")
	for name, fun := range l { /* fill the table with given functions */
		for i := 0; i < nup; i++ { /* copy upvalues to the top */
			self.PushValue(-nup)
		}
		// r[-(nup+2)][name]=fun
		self.PushGoClosure(fun, nup) /* closure with those upvalues */
		self.SetField(-(nup + 2), name)
	}
	self.Pop(nup) /* remove upvalues */
}

func (self *luaState) intError(arg int) {
	if self.IsNumber(arg) {
		self.ArgError(arg, "number has no integer representation")
	} else {
		self.tagError(arg, LUA_TNUMBER)
	}
}

func (self *luaState) tagError(arg int, tag LuaType) {
	self.typeError(arg, self.TypeName(LuaType(tag)))
}

func (self *luaState) typeError(arg int, tname string) int {
	var typeArg string /* name for the type of the actual argument */
	if self.GetMetafield(arg, "__name") == LUA_TSTRING {
		typeArg = self.ToString(-1) /* use the given type name */
	} else if self.Type(arg) == LUA_TLIGHTUSERDATA {
		typeArg = "light userdata" /* special name for messages */
	} else {
		typeArg = self.TypeName2(arg) /* standard name */
	}
	msg := tname + " expected, got " + typeArg
	self.PushString(msg)
	return self.ArgError(arg, msg)
}

func (ls *luaState) Traceback(ls1 LuaState, msg string, level int) {
	ar := LuaDebug{}
	top := ls.GetTop()
	last := lastlevel(ls1)
	n1 := -1
	if last-level > LEVELS1+LEVELS2 {
		n1 = LEVELS1
	}
	if msg != "" {
		ls.PushFString("%s\n", msg)
	}
	ls.CheckStack2(10, "")
	ls.PushString("stack traceback:")
	for ls1.GetStack(level, &ar) {
		level++
		if n1--; n1 == 0 { /* too many levels? */
			ls.PushString("\n\t...")   /* add a '...' */
			level = last - LEVELS2 + 1 /* and skip to last ones */
		} else {
			ls1.GetInfo("Slnt", &ar)
			ls.PushFString("\n\t%s:", ar.ShortSrc)
			if ar.CurrentLine > 0 {
				ls.PushFString("%d:", ar.CurrentLine)
			}
			ls.PushString(" in ")
			pushfuncname(ls, &ar)
			if ar.IsTailCall {
				ls.PushString("\n\t(...tail calls...)")
			}
			ls.Concat(ls.GetTop() - top)
		}
	}
	ls.Concat(ls.GetTop() - top)
}

func lastlevel(ls LuaState) int {
	ar := LuaDebug{}
	li := 1
	le := 1
	/* find an upper bound */
	for ls.GetStack(le, &ar) {
		li = le
		le *= 2
	}
	/* do a binary search */
	for li < le {
		m := (li + le) / 2
		if ls.GetStack(m, &ar) {
			li = m + 1
		} else {
			le = m
		}
	}
	return le - 1
}

func pushfuncname(ls LuaState, ar *LuaDebug) {
	if pushglobalfuncname(ls, ar) { /* try first a global name */
		ls.PushFString("function '%s'", ls.ToString(-1))
		ls.Remove(-2) /* remove name */
	} else if ar.NameWhat != "" { /* is there a name from code? */
		ls.PushFString("%s '%s'", ar.NameWhat, ar.Name) /* use it */
	} else if ar.What == "main" { /* main? */
		ls.PushString("main chunk")
	} else if ar.What != "C" { /* for Lua functions, use <file:line> */
		ls.PushFString("function <%s:%d>", ar.ShortSrc, ar.LineDefined)
	} else { /* nothing left... */
		ls.PushString("?")
	}
}

/*
** Search for a name for a function in all loaded modules
 */
func pushglobalfuncname(ls LuaState, ar *LuaDebug) bool {
	top := ls.GetTop()
	ls.GetInfo("f", ar) /* push function */
	ls.GetField(LUA_REGISTRYINDEX, LUA_LOADED_TABLE)
	if findfield(ls, top+1, 2) {
		name := ls.ToString(-1)
		if strings.HasPrefix(name, "_G.") { /* name start with '_G.'? */
			ls.PushString(name[3:]) /* push name without prefix */
			ls.Remove(-2)           /* remove original name */
		}
		ls.Copy(-1, top+1) /* move name to proper place */
		ls.Pop(2)          /* remove pushed values */
		return true
	} else {
		ls.SetTop(top) /* remove function and global table */
		return false
	}
}

/*
** search for 'objidx' in table at index -1.
** return 1 + string at top if find a good name.
 */
func findfield(ls LuaState, objidx, level int) bool {
	if level == 0 || !ls.IsTable(-1) {
		return false /* not found */
	}
	ls.PushNil()      /* start 'next' loop */
	for ls.Next(-2) { /* for each pair in table */
		if ls.Type(-2) == LUA_TSTRING { /* ignore non-string keys */
			if ls.RawEqual(objidx, -1) { /* found object? */
				ls.Pop(1) /* remove value (but keep name) */
				return true
			} else if findfield(ls, objidx, level-1) { /* try recursively */
				ls.Remove(-2) /* remove table (but keep name) */
				ls.PushString(".")
				ls.Insert(-2) /* place '.' between the two names */
				ls.Concat(3)
				return true
			}
		}
		ls.Pop(1) /* remove value */
	}
	return false /* not found */
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_getmetatable
// lua-5.3.4/src/lauxlib.h#luaL_getmetatable()
// func (self *luaState) GetMetatable2(tname string) LuaType {
// 	return self.GetField(LUA_REGISTRYINDEX, tname)
// }
