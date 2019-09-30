package state

import (
	"fmt"

	. "github.com/zxh0/lua.go/api"
	"github.com/zxh0/lua.go/stdlib"
	"io/ioutil"
	"strings"
)

/* key, in the registry, for table of loaded modules */
const LUA_LOADED_TABLE = "_LOADED"

const LEVELS1 = 10 /* size of the first part of the stack */
const LEVELS2 = 11 /* size of the second part of the stack */

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_error
func (state *luaState) Error2(fmt string, a ...interface{}) int {
	state.PushFString(fmt, a...) // todo
	return state.Error()
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_argerror
func (state *luaState) ArgError(arg int, extraMsg string) int {
	// bad argument #arg to 'funcname' (extramsg)
	return state.Error2("bad argument #%d (%s)", arg, extraMsg) // todo
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_where
// lua-5.3.4/src/lauxlib.c#luaL_where()
func (state *luaState) Where(level int) {
	// chunkname:currentline:
	ar := LuaDebug{}
	if state.GetStack(level, &ar) { /* check function at level */
		state.GetInfo("Sl", &ar) /* get info about it */
		if ar.CurrentLine > 0 {  /* is there info? */
			state.PushFString("%s:%d: ", ar.ShortSrc, ar.CurrentLine)
			return
		}
	}
	state.PushFString("") /* else, no information available... */
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkstack
// lua-5.3.4/src/lauxlib.c#luaL_checkstack()
func (state *luaState) CheckStack2(sz int, msg string) {
	if !state.CheckStack(sz) {
		if msg != "" {
			state.Error2("stack overflow (%s)", msg)
		} else {
			state.Error2("stack overflow")
		}
	}
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_argcheck
// lua-5.3.4/src/lauxlib.c#luaL_argcheck()
func (state *luaState) ArgCheck(cond bool, arg int, extraMsg string) {
	if !cond {
		state.ArgError(arg, extraMsg)
	}
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkany
// lua-5.3.4/src/lauxlib.c#luaL_checkany()
func (state *luaState) CheckAny(arg int) {
	if state.Type(arg) == LUA_TNONE {
		state.ArgError(arg, "value expected")
	}
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checktype
// lua-5.3.4/src/lauxlib.c#luaL_checktype()
func (state *luaState) CheckType(arg int, t LuaType) {
	if state.Type(arg) != t {
		state.tagError(arg, t)
	}
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkinteger
// lua-5.3.4/src/lauxlib.c#luaL_checkinteger()
func (state *luaState) CheckInteger(arg int) int64 {
	i, ok := state.ToIntegerX(arg)
	if !ok {
		state.intError(arg)
	}
	return i
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checknumber
// lua-5.3.4/src/lauxlib.c#luaL_checknumber()
func (state *luaState) CheckNumber(arg int) float64 {
	f, ok := state.ToNumberX(arg)
	if !ok {
		state.tagError(arg, LUA_TNUMBER)
	}
	return f
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkstring
// http://www.lua.org/manual/5.3/manual.html#luaL_checklstring
// lua-5.3.4/src/lauxlib.c#luaL_checklstring()
func (state *luaState) CheckString(arg int) string {
	s, ok := state.ToStringX(arg)
	if !ok {
		state.tagError(arg, LUA_TSTRING)
	}
	return s
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_optinteger
// lua-5.3.4/src/lauxlib.c#luaL_optinteger()
func (state *luaState) OptInteger(arg int, def int64) int64 {
	if state.IsNoneOrNil(arg) {
		return def
	}
	return state.CheckInteger(arg)
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_optnumber
// lua-5.3.4/src/lauxlib.c#luaL_optnumber()
func (state *luaState) OptNumber(arg int, def float64) float64 {
	if state.IsNoneOrNil(arg) {
		return def
	}
	return state.CheckNumber(arg)
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_optstring
// lua-5.3.4/src/lauxlib.c#luaL_optlstring()
func (state *luaState) OptString(arg int, def string) string {
	if state.IsNoneOrNil(arg) {
		return def
	}
	return state.CheckString(arg)
}

// [-0, +?, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_dofile
// lua-5.3.4/src/lauxlib.h#luaL_dofile()
func (state *luaState) DoFile(filename string) bool {
	return state.LoadFile(filename) == LUA_OK &&
		state.PCall(0, LUA_MULTRET, 0) == LUA_OK
}

// [-0, +?, –]
// http://www.lua.org/manual/5.3/manual.html#luaL_dostring
// lua-5.3.4/src/lauxlib.h#luaL_dostring()
func (state *luaState) DoString(str string) bool {
	return state.LoadString(str) == LUA_OK &&
		state.PCall(0, LUA_MULTRET, 0) == LUA_OK
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_loadfile
// lua-5.3.4/src/lauxlib.h#luaL_loadfile()
func (state *luaState) LoadFile(filename string) ThreadStatus {
	return state.LoadFileX(filename, "bt")
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_loadfilex
func (state *luaState) LoadFileX(filename, mode string) ThreadStatus {
	if data, err := ioutil.ReadFile(filename); err == nil {
		return state.Load(data, "@"+filename, mode)
	}
	return LUA_ERRFILE // todo
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#luaL_loadstring
func (state *luaState) LoadString(s string) ThreadStatus {
	return state.Load([]byte(s), s, "bt")
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkversion
func (state *luaState) CheckVersion() {
	//panic("todo: CheckVersion!")
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#luaL_typename
// lua-5.3.4/src/lauxlib.h#luaL_typename()
func (state *luaState) TypeName2(idx int) string {
	return state.TypeName(state.Type(idx))
}

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_len
// lua-5.3.4/src/lauxlib.c#luaL_len()
func (state *luaState) Len2(idx int) int64 {
	state.Len(idx)
	i, isNum := state.ToIntegerX(-1)
	if !isNum {
		state.Error2("object length is not an integer")
	}
	state.Pop(1)
	return i
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_tolstring
// lua-5.3.4/src/lauxlib.c#luaL_tolstring()
func (state *luaState) ToString2(idx int) string {
	if state.CallMeta(idx, "__tostring") { /* metafield? */
		if !state.IsString(-1) {
			state.Error2("'__tostring' must return a string")
		}
	} else {
		switch state.Type(idx) {
		case LUA_TNUMBER:
			if state.IsInteger(idx) {
				state.PushString(fmt.Sprintf("%d", state.ToInteger(idx))) // todo
			} else {
				state.PushString(fmt.Sprintf("%g", state.ToNumber(idx))) // todo
			}
		case LUA_TSTRING:
			state.PushValue(idx)
		case LUA_TBOOLEAN:
			if state.ToBoolean(idx) {
				state.PushString("true")
			} else {
				state.PushString("false")
			}
		case LUA_TNIL:
			state.PushString("nil")
		default:
			tt := state.GetMetafield(idx, "__name") /* try name */
			var kind string
			if tt == LUA_TSTRING {
				kind = state.CheckString(-1)
			} else {
				kind = state.TypeName2(idx)
			}

			state.PushString(fmt.Sprintf("%s: %p", kind, state.ToPointer(idx)))
			if tt != LUA_TNIL {
				state.Remove(-2) /* remove '__name' */
			}
		}
	}
	return state.CheckString(-1)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_getsubtable
// lua-5.3.4/src/lauxlib.c#luaL_getsubtable()
func (state *luaState) GetSubTable(idx int, fname string) bool {
	if state.GetField(idx, fname) == LUA_TTABLE {
		return true /* table already there */
	}
	state.Pop(1) /* remove previous result */
	idx = state.stack.absIndex(idx)
	state.NewTable()
	state.PushValue(-1)        /* copy to be left at top */
	state.SetField(idx, fname) /* assign new table to field */
	return false               /* false, because did not find table there */
}

// [-0, +(0|1), m]
// http://www.lua.org/manual/5.3/manual.html#luaL_getmetafield
// lua-5.3.4/src/lauxlib.c#luaL_getmetafield()
func (state *luaState) GetMetafield(obj int, event string) LuaType {
	if !state.GetMetatable(obj) { /* no metatable? */
		return LUA_TNIL
	}

	state.PushString(event)
	tt := state.RawGet(-2)
	if tt == LUA_TNIL { /* is metafield nil? */
		state.Pop(2) /* remove metatable and metafield */
	} else {
		state.Remove(-2) /* remove only metatable */
	}
	return tt /* return metafield type */
}

// [-0, +(0|1), e]
// http://www.lua.org/manual/5.3/manual.html#luaL_callmeta
// lua-5.3.4/src/lauxlib.c#luaL_callmeta()
func (state *luaState) CallMeta(obj int, event string) bool {
	obj = state.AbsIndex(obj)
	if state.GetMetafield(obj, event) == LUA_TNIL { /* no metafield? */
		return false
	}

	state.PushValue(obj)
	state.Call(1, 1)
	return true
}

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_openlibs
// lua-5.3.4/src/linit.c#luaL_openlibs()
func (state *luaState) OpenLibs() {
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
		state.RequireF(name, fun, true)
		state.Pop(1)
	}
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_requiref
// lua-5.3.4/src/lauxlib.c#luaL_requiref()
func (state *luaState) RequireF(modname string, openf GoFunction, glb bool) {
	state.GetSubTable(LUA_REGISTRYINDEX, "_LOADED")
	state.GetField(-1, modname)
	if !state.ToBoolean(-1) { /* package not already loaded? */
		state.Pop(1) /* remove field */
		state.PushGoFunction(openf)
		state.PushString(modname)   /* argument to open function */
		state.Call(1, 1)            /* call 'openf' to open module */
		state.PushValue(-1)         /* make copy of module (call result) */
		state.SetField(-3, modname) /* _LOADED[modname] = module */
	}
	state.Remove(-2) /* remove _LOADED table */
	if glb {
		state.PushValue(-1)      /* copy of module */
		state.SetGlobal(modname) /* _G[modname] = module */
	}
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_newlib
// lua-5.3.4/src/lauxlib.h#luaL_newlib()
func (state *luaState) NewLib(l FuncReg) {
	state.NewLibTable(l)
	state.SetFuncs(l, 0)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_newlibtable
// lua-5.3.4/src/lauxlib.h#luaL_newlibtable()
func (state *luaState) NewLibTable(l FuncReg) {
	state.CreateTable(0, len(l))
}

// [-nup, +0, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_setfuncs
// lua-5.3.4/src/lauxlib.c#luaL_setfuncs()
func (state *luaState) SetFuncs(l FuncReg, nup int) {
	state.CheckStack2(nup, "too many upvalues")
	for name, fun := range l { /* fill the table with given functions */
		for i := 0; i < nup; i++ { /* copy upvalues to the top */
			state.PushValue(-nup)
		}
		// r[-(nup+2)][name]=fun
		state.PushGoClosure(fun, nup) /* closure with those upvalues */
		state.SetField(-(nup + 2), name)
	}
	state.Pop(nup) /* remove upvalues */
}

func (state *luaState) intError(arg int) {
	if state.IsNumber(arg) {
		state.ArgError(arg, "number has no integer representation")
	} else {
		state.tagError(arg, LUA_TNUMBER)
	}
}

func (state *luaState) tagError(arg int, tag LuaType) {
	state.typeError(arg, state.TypeName(LuaType(tag)))
}

func (state *luaState) typeError(arg int, tname string) int {
	var typeArg string /* name for the type of the actual argument */
	if state.GetMetafield(arg, "__name") == LUA_TSTRING {
		typeArg = state.ToString(-1) /* use the given type name */
	} else if state.Type(arg) == LUA_TLIGHTUSERDATA {
		typeArg = "light userdata" /* special name for messages */
	} else {
		typeArg = state.TypeName2(arg) /* standard name */
	}
	msg := tname + " expected, got " + typeArg
	state.PushString(msg)
	return state.ArgError(arg, msg)
}

func (state *luaState) Traceback(ls1 LuaState, msg string, level int) {
	ar := LuaDebug{}
	top := state.GetTop()
	last := lastlevel(ls1)
	n1 := -1
	if last-level > LEVELS1+LEVELS2 {
		n1 = LEVELS1
	}
	if msg != "" {
		state.PushFString("%s\n", msg)
	}
	state.CheckStack2(10, "")
	state.PushString("stack traceback:")
	for ls1.GetStack(level, &ar) {
		level++
		if n1--; n1 == 0 { /* too many levels? */
			state.PushString("\n\t...") /* add a '...' */
			level = last - LEVELS2 + 1  /* and skip to last ones */
		} else {
			ls1.GetInfo("Slnt", &ar)
			state.PushFString("\n\t%s:", ar.ShortSrc)
			if ar.CurrentLine > 0 {
				state.PushFString("%d:", ar.CurrentLine)
			}
			state.PushString(" in ")
			pushfuncname(state, &ar)
			if ar.IsTailCall {
				state.PushString("\n\t(...tail calls...)")
			}
			state.Concat(state.GetTop() - top)
		}
	}
	state.Concat(state.GetTop() - top)
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
// func (state *luaState) GetMetatable2(tname string) LuaType {
// 	return state.GetField(LUA_REGISTRYINDEX, tname)
// }
