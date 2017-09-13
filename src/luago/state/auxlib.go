package state

import "io/ioutil"
import . "luago/api"
import "luago/stdlib"

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_error
func (self *luaState) Error2(fmt string) {
	panic("todo: Error2!")
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_argerror
func (self *luaState) ArgError(arg int, extraMsg string) int {
	// bad argument #arg to 'funcname' (extramsg)
	panic("todo: ArgError!")
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkstack
// lua-5.3.4/src/lauxlib.c#luaL_checkstack()
func (self *luaState) CheckStack2(sz int, msg string) {
	if !self.CheckStack(sz) {
		if msg != "" {
			self.Error2("stack overflow (" + msg + ")")
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
// http://www.lua.org/manual/5.3/manual.html#luaL_checkinteger
// lua-5.3.4/src/lauxlib.c#luaL_checkinteger()
func (self *luaState) CheckInteger(arg int) int64 {
	if i, ok := self.ToIntegerX(arg); ok {
		return i
	} else {
		self.intError(arg)
		panic("unreachable!")
	}
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checknumber
// lua-5.3.4/src/lauxlib.c#luaL_checknumber()
func (self *luaState) CheckNumber(arg int) float64 {
	if f, ok := self.ToNumberX(arg); ok {
		return f
	} else {
		self.tagError(arg, LUA_TNUMBER)
		panic("unreachable!")
	}
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkstring
// http://www.lua.org/manual/5.3/manual.html#luaL_checklstring
// lua-5.3.4/src/lauxlib.c#luaL_checklstring()
func (self *luaState) CheckString(arg int) string {
	if s, ok := self.ToString(arg); ok {
		return s
	} else {
		self.tagError(arg, LUA_TSTRING)
		panic("unreachable!")
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
// http://www.lua.org/manual/5.3/manual.html#luaL_optinteger
// lua-5.3.4/src/lauxlib.c#luaL_optinteger()
func (self *luaState) OptInteger(arg int, def int64) int64 {
	if self.IsNoneOrNil(arg) {
		return def
	} else {
		return self.CheckInteger(arg)
	}
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_optnumber
// lua-5.3.4/src/lauxlib.c#luaL_optnumber()
func (self *luaState) OptNumber(arg int, def float64) float64 {
	if self.IsNoneOrNil(arg) {
		return def
	} else {
		return self.CheckNumber(arg)
	}
}

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_optstring
// lua-5.3.4/src/lauxlib.c#luaL_optlstring()
func (self *luaState) OptString(arg int, def string) string {
	if self.IsNoneOrNil(arg) {
		return def
	} else {
		return self.CheckString(arg)
	}
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
	return self.LoadFileX(filename, "")
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_loadfilex
func (self *luaState) LoadFileX(filename, mode string) ThreadStatus {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	self.Load(data, filename, mode)
	// panic("todo!")
	return LUA_OK
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#luaL_loadstring
func (self *luaState) LoadString(s string) ThreadStatus {
	panic("todo: LoadString!")
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_getmetatable
// lua-5.3.4/src/lauxlib.h#luaL_getmetatable()
func (self *luaState) GetMetatable2(tname string) LuaType {
	return self.GetField(LUA_REGISTRYINDEX, tname)
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
func (self *luaState) RequireF(modname string, openf GoFunction, glb bool) {
	self.GetSubTable(LUA_REGISTRYINDEX, "_LOADED") // ~/_LOADED
	self.GetField(-1, modname)                     // ~/_LOADED/_LOADED[modname]
	if !self.ToBoolean(-1) {                       /* package not already loaded? */
		self.Pop(1)                // ~/_LOADED               /* remove field */
		self.PushGoFunction(openf) // ~/_LOADED/openf
		self.PushString(modname)   // ~/_LOADED/openf/modname /* argument to open function */
		self.Call(1, 1)            // ~/_LOADED/module        /* call 'openf' to open module */
		self.PushValue(-1)         // ~/_LOADED/module/module /* make copy of module (call result) */
		self.SetField(-3, modname) // ~/_LOADED/module        /* _LOADED[modname] = module */
	}
	self.Remove(-2) // ~/module /* remove _LOADED table */
	if glb {
		self.PushValue(-1)      /* copy of module */
		self.SetGlobal(modname) /* _G[modname] = module */
	}
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_newlib
func (self *luaState) NewLib(l FuncReg) {
	self.NewLibTable(l)
	self.SetFuncs(l, 0)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_newlibtable
func (self *luaState) NewLibTable(l FuncReg) {
	self.CreateTable(0, len(l))
}

// [-nup, +0, m]
// http://www.lua.org/manual/5.3/manual.html#luaL_setfuncs
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

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_getsubtable
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

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_len
func (self *luaState) Len2(index int) int64 {
	self.Len(index)
	if i, ok := self.ToIntegerX(-1); ok {
		self.Pop(1)
		return i
	} else {
		panic("todo!")
	}
}

func (self *luaState) TypeName2(index int) string {
	panic("todo!")
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#luaL_tolstring
func (self *luaState) ToString2(idx int) string {
	panic("todo!")
}

/*
LUALIB_API const char *luaL_tolstring (lua_State *L, int idx, size_t *len) {
  if (luaL_callmeta(L, idx, "__tostring")) {  /* metafield? * /
    if (!lua_isstring(L, -1))
      luaL_error(L, "'__tostring' must return a string");
  }
  else {
    switch (lua_type(L, idx)) {
      case LUA_TNUMBER: {
        if (lua_isinteger(L, idx))
          lua_pushfstring(L, "%I", (LUAI_UACINT)lua_tointeger(L, idx));
        else
          lua_pushfstring(L, "%f", (LUAI_UACNUMBER)lua_tonumber(L, idx));
        break;
      }
      case LUA_TSTRING:
        lua_pushvalue(L, idx);
        break;
      case LUA_TBOOLEAN:
        lua_pushstring(L, (lua_toboolean(L, idx) ? "true" : "false"));
        break;
      case LUA_TNIL:
        lua_pushliteral(L, "nil");
        break;
      default: {
        int tt = luaL_getmetafield(L, idx, "__name");  /* try name * /
        const char *kind = (tt == LUA_TSTRING) ? lua_tostring(L, -1) :
                                                 luaL_typename(L, idx);
        lua_pushfstring(L, "%s: %p", kind, lua_topointer(L, idx));
        if (tt != LUA_TNIL)
          lua_remove(L, -2);  /* remove '__name' * /
        break;
      }
    }
  }
  return lua_tolstring(L, -1, len);
}
*/

// [-0, +0, v]
// http://www.lua.org/manual/5.3/manual.html#luaL_checkversion
func (self *luaState) CheckVersion() {
	//panic("todo: CheckVersion!")
}


func (self *luaState) intError(arg int) {
	if self.IsNumber(arg) {
		self.ArgError(arg, "number has no integer representation")
	} else {

		self.tagError(arg, LUA_TNUMBER)
	}
}

func (self *luaState) tagError(arg int, tag LuaType) {
	//self.typeError(arg, self.TypeName(LuaType(tag)))
	panic("todo!")
}

// func (self *luaState) typeError(arg int, tname string) int {
// 	var typearg string /* name for the type of the actual argument */
// 	if self.GetMetafield(arg, "__name") == LUA_TSTRING {
// 		typearg, _ = self.ToString(-1) /* use the given type name */
// 	//} else if self.Type(arg) == LUA_TLIGHTUSERDATA {
// 	//	typearg = "light userdata" /* special name for messages */
// 	} else {
// 		typearg = self.TypeName2(arg) /* standard name */
// 	}
// 	msg := self.PushFString("%s expected, got %s", tname, typearg)
// 	return self.ArgError(arg, msg)
// }

