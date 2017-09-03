package stdlib

import "fmt"
import . "luago/api"
import "luago/luanum"

var baseFuncs = map[string]GoFunction{
	"print":          basePrint,
	"assert":         baseAssert,
	"error":          baseError,
	"select":         baseSelect,
	"ipairs":         baseIPairs,
	"pairs":          basePairs,
	"next":           baseNext,
	"load":           baseLoad,
	"loadfile":       baseLoadFile,
	"dofile":         baseDoFile,
	"pcall":          basePCall,
	"xpcall":         baseXpcall,
	"getmetatable":   baseGetMetatable,
	"setmetatable":   baseSetMetatable,
	"rawequal":       baseRawEqual,
	"rawlen":         baseRawLen,
	"rawget":         baseRawGet,
	"rawset":         baseRawSet,
	"tonumber":       baseToNumber,
	"tostring":       baseToString,
	"type":           baseType,
	"collectgarbage": baseCollectGarbage,
	/* placeholders */
	"_G":       nil,
	"_VERSION": nil,
}

func OpenBaseLib(ls LuaState) int {
	/* open lib into global table */
	ls.PushGlobalTable()
	ls.SetFuncs(baseFuncs, 0)
	/* set global _G */
	ls.PushValue(-1)
	ls.SetField(-2, "_G")
	/* set global _VERSION */
	ls.PushString(LUA_VERSION)
	ls.SetField(-2, "_VERSION")
	return 1
}

// print (···)
// http://www.lua.org/manual/5.3/manual.html#pdf-print
func basePrint(ls LuaState) int {
	nArgs := ls.GetTop()
	ls.CheckStack(2)
	for i := 1; i <= nArgs; i++ {
		ls.PushGoFunction(baseToString)
		ls.PushValue(i)
		ls.Call(1, 1)
		s, _ := ls.ToString(-1)
		ls.Pop(1)

		fmt.Print(s)
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}

// assert (v [, message])
// http://www.lua.org/manual/5.3/manual.html#pdf-assert
func baseAssert(ls LuaState) int {
	if ls.ToBoolean(1) { /* condition is true? */
		return ls.GetTop() /* return all arguments */
	} else { /* error */
		ls.CheckAny(1)                     /* there must be a condition */
		ls.Remove(1)                       /* remove it */
		ls.PushString("assertion failed!") /* default message */
		ls.SetTop(1)                       /* leave only message (default if no other one) */
		return baseError(ls)               /* call 'error' */
	}
}

// error (message [, level])
// http://www.lua.org/manual/5.3/manual.html#pdf-error
func baseError(ls LuaState) int {
	panic("todo! baseError")
}

// select (index, ···)
// http://www.lua.org/manual/5.3/manual.html#pdf-select
func baseSelect(ls LuaState) int {
	n := int64(ls.GetTop())
	if ls.Type(1) == LUA_TSTRING && ls.CheckString(1) == "#" {
		ls.PushInteger(n - 1)
		return 1
	} else {
		i := ls.CheckInteger(1)
		if i < 0 {
			i = n + i
		} else if i > n {
			i = n
		}
		ls.ArgCheck(1 <= i, 1, "index out of range")
		return int(n - i)
	}
}

// ipairs (t)
// http://www.lua.org/manual/5.3/manual.html#pdf-ipairs
func baseIPairs(ls LuaState) int {
	ls.CheckAny(1)
	ls.PushGoFunction(iPairsAux) /* iteration function */
	ls.PushValue(1)              /* state */
	ls.PushInteger(0)            /* initial value */
	return 3
}

func iPairsAux(ls LuaState) int {
	i := ls.CheckInteger(2) + 1
	ls.PushInteger(i)
	if ls.GetI(1, i) == LUA_TNIL {
		return 1
	} else {
		return 2
	}
}

// pairs (t)
// http://www.lua.org/manual/5.3/manual.html#pdf-pairs
func basePairs(ls LuaState) int {
	ls.PushGoFunction(baseNext) /* will return generator, */
	ls.PushValue(1)             /* state, */
	ls.PushNil()
	return 3
}

// next (table [, index])
// http://www.lua.org/manual/5.3/manual.html#pdf-next
func baseNext(ls LuaState) int {
	if ls.GetTop() < 2 {
		ls.PushNil()
	}
	if ls.Next(1) {
		return 2
	} else {
		ls.PushNil()
		return 1
	}
}

// load (chunk [, chunkname [, mode [, env]]])
// http://www.lua.org/manual/5.3/manual.html#pdf-load
func baseLoad(ls LuaState) int {
	panic("todo! baseLoad")
}

// loadfile ([filename [, mode [, env]]])
// http://www.lua.org/manual/5.3/manual.html#pdf-loadfile
func baseLoadFile(ls LuaState) int {
	panic("todo! baseLoadFile")
}

// dofile ([filename])
// http://www.lua.org/manual/5.3/manual.html#pdf-dofile
func baseDoFile(ls LuaState) int {
	fname := ls.OptString(1, "")
	ls.SetTop(1)
	if ls.LoadFile(fname) != LUA_OK {
		//return lua_error(L);
		panic("todo!")
	}
	//ls.CallK(0, LUA_MULTRET, 0, dofilecont);
	//return dofilecont(L, 0, 0);
	panic("todo!")
}

// pcall (f [, arg1, ···])
// http://www.lua.org/manual/5.3/manual.html#pdf-pcall
func basePCall(ls LuaState) int {
	nArgs := ls.GetTop() - 1
	status := ls.PCall(nArgs, -1, 0)
	if status == LUA_OK {
		ls.PushBoolean(true)
	} else {
		ls.PushBoolean(false)
	}
	ls.Rotate(1, 1)
	return ls.GetTop()
}

// xpcall (f, msgh [, arg1, ···])
// http://www.lua.org/manual/5.3/manual.html#pdf-xpcall
func baseXpcall(ls LuaState) int {
	panic("todo! baseXpcall")
}

// getmetatable (object)
// http://www.lua.org/manual/5.3/manual.html#pdf-getmetatable
func baseGetMetatable(ls LuaState) int {
	if !ls.GetMetatable(1) {
		ls.PushNil()
		return 1
	}

	ls.PushString("__metatable")
	ls.GetTable(-2)
	if !ls.IsNil(-1) {
		return 1
	} else {
		ls.Pop(1)
		return 1
	}
}

// setmetatable (table, metatable)
// http://www.lua.org/manual/5.3/manual.html#pdf-setmetatable
func baseSetMetatable(ls LuaState) int {
	if ls.GetMetatable(1) {
		ls.PushString("__metatable")
		ls.GetTable(-2)
		if !ls.IsNil(-1) {
			panic("cannot change a protected metatable") // todo
		} else {
			ls.Pop(2)
		}
	}

	ls.SetMetatable(1)
	return 1
}

// rawequal (v1, v2)
// http://www.lua.org/manual/5.3/manual.html#pdf-rawequal
func baseRawEqual(ls LuaState) int {
	ls.PushBoolean(ls.RawEqual(1, 2))
	return 1
}

// rawlen (v)
// http://www.lua.org/manual/5.3/manual.html#pdf-rawlen
func baseRawLen(ls LuaState) int {
	rl := int64(ls.RawLen(1))
	ls.PushInteger(rl)
	return 1
}

// rawget (table, index)
// http://www.lua.org/manual/5.3/manual.html#pdf-rawget
func baseRawGet(ls LuaState) int {
	// todo
	ls.RawGet(1)
	return 1
}

// rawset (table, index, value)
// http://www.lua.org/manual/5.3/manual.html#pdf-rawset
func baseRawSet(ls LuaState) int {
	// todo
	ls.RawSet(1)
	return 1
}

// tonumber (e [, base])
// http://www.lua.org/manual/5.3/manual.html#pdf-tonumber
// lua-5.3.4/src/lbaselib.c#luaB_tonumber()
func baseToNumber(ls LuaState) int {
	if ls.IsNoneOrNil(2) { /* standard conversion? */
		ls.CheckAny(1)
		if ls.Type(1) == LUA_TNUMBER { /* already a number? */
			ls.SetTop(1) /* yes; return it */
			return 1
		} else {
			if s, ok := ls.ToString(1); ok {
				if ok && ls.StringToNumber(s) {
					return 1 /* successful conversion to number */
				} /* else not a number */
			}
		}
	} else {
		ls.CheckType(1, LUA_TSTRING) /* no numbers as strings */
		s, _ := ls.ToString(1)
		base := int(ls.CheckInteger(2))
		ls.ArgCheck(2 <= base && base <= 36, 2, "base out of range")
		if n, ok := luanum.ParseInteger(s, base); ok {
			ls.PushInteger(n)
			return 1
		} /* else not a number */
	} /* else not a number */
	ls.PushNil() /* not a number */
	return 1
}

// tostring (v)
// http://www.lua.org/manual/5.3/manual.html#pdf-tostring
func baseToString(ls LuaState) int {
	ls.CheckStack(4)
	ls.GetMetatable(1)  // v/mt
	if ls.IsTable(-1) { //
		ls.PushString("__tostring") // v/mt/"__tostring"
		ls.GetTable(-2)             // v/mt/__tostring
		if ls.IsFunction(-1) {      //
			ls.PushValue(1) // v/mt/__tostring/v
			ls.Call(1, 1)   // v/mt/result
			s := castToString(ls, -1)
			ls.PushString(s)
			return 1
		}
	}

	s := castToString(ls, 1)
	ls.PushString(s)
	return 1
}

// type (v)
// http://www.lua.org/manual/5.3/manual.html#pdf-type
func baseType(ls LuaState) int {
	luaType := ls.Type(1)
	typeName := ls.TypeName(luaType)
	ls.PushString(typeName)
	return 1
}

// collectgarbage ([opt [, arg]])
// http://www.lua.org/manual/5.3/manual.html#pdf-collectgarbage
func baseCollectGarbage(ls LuaState) int {
	panic("todo! baseCollectGarbage")
}
