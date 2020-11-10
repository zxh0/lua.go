package state

import (
	. "github.com/zxh0/lua.go/api"
	"github.com/zxh0/lua.go/binchunk"
	"github.com/zxh0/lua.go/compiler"
	"github.com/zxh0/lua.go/vm"
)

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_dump
func (state *luaState) Dump(strip bool) []byte {
	v := state.stack.get(-1)
	if c, ok := v.(*closure); ok {
		return binchunk.Dump(c.proto)
	}
	panic("todo!")
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_load
func (state *luaState) Load(chunk []byte, chunkName, mode string) (status ThreadStatus) {
	status = LUA_ERRSYNTAX

	// catch error
	defer func() {
		if r := recover(); r != nil {
			state.stack.push(_getErrObj(r))
		}
	}()

	var proto *binchunk.Prototype
	if binchunk.IsBinaryChunk(chunk) {
		if mode == "t" {
			panic("attempt to load a binary chunk (mode is '" + mode + "')")
		}
		proto = binchunk.Undump(chunk)
	} else {
		if mode == "b" {
			panic("attempt to load a text chunk (mode is '" + mode + "')")
		}
		proto = compiler.Compile(string(chunk), chunkName)
	}

	c := newLuaClosure(proto)
	if len(proto.Upvalues) > 0 {
		env := state.registry.get(LUA_RIDX_GLOBALS)
		c.upvals[0] = &upvalue{&env}
	}
	state.stack.push(c)
	status = LUA_OK
	return
}

// [-(nargs+1), +nresults, e]
// http://www.lua.org/manual/5.3/manual.html#lua_call
func (state *luaState) Call(nArgs, nResults int) {
	val := state.stack.get(-(nArgs + 1))

	c, ok := val.(*closure)
	if !ok {
		if mf := getMetafield(val, "__call", state); mf != nil {
			if c, ok = mf.(*closure); ok {
				state.stack.push(val)
				state.Insert(-(nArgs + 2))
				nArgs += 1
			}
		}
	}

	if ok {
		if c.proto != nil {
			state.callLuaClosure(nArgs, nResults, c)
		} else {
			state.callGoClosure(nArgs, nResults, c)
		}
	} else {
		typeName := state.TypeName(typeOf(val))
		panic("attempt to call a " + typeName + " value")
	}
}

func (state *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	// create new lua stack
	newStack := newLuaStack(nArgs+LUA_MINSTACK, state)
	newStack.closure = c

	// pass args, pop func
	if nArgs > 0 {
		args := state.stack.popN(nArgs)
		newStack.pushN(args, nArgs)
	}
	state.stack.pop()

	// run closure
	state.pushLuaStack(newStack)
	r := c.goFunc(state)
	state.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(r)
		state.stack.check(len(results))
		state.stack.pushN(results, nResults)
	}
}

func (state *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// create new lua stack
	newStack := newLuaStack(nRegs+LUA_MINSTACK, state)
	newStack.closure = c

	// pass args, pop func
	funcAndArgs := state.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	// run closure
	state.pushLuaStack(newStack)
	state.runLuaClosure()
	state.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		state.stack.check(len(results))
		state.stack.pushN(results, nResults)
	}
}

func (state *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(state.Fetch())
		inst.Execute(state)

		// indent := fmt.Sprintf("%%%ds", state.callDepth*2)
		// fmt.Printf(indent+"[%02d: %s] => %s\n",
		// 	"", pc+1, inst.OpName(), state)
		//println(inst.OpName())

		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

// Calls a function in protected mode.
// http://www.lua.org/manual/5.3/manual.html#lua_pcall
func (state *luaState) PCall(nArgs, nResults, msgh int) (status ThreadStatus) {
	status = LUA_ERRRUN
	caller := state.stack
	handler := state.stack.get(msgh)

	// catch error
	defer func() {
		if r := recover(); r != nil { // todo
			err := _getErrObj(r)
			if msgh != 0 {
				if handler == nil {
					panic(err) // todo
				}

				state.stack.push(handler)
				state.stack.push(err)
				state.PCall(1, 1, 0)
				err = state.stack.pop()
			}
			for state.stack != caller {
				state.popLuaStack()
			}
			state.stack.push(err)
		}
	}()

	state.Call(nArgs, nResults)
	status = LUA_OK
	return
}

func _getErrObj(err interface{}) luaValue {
	if t, ok := err.(*luaTable); ok {
		return t.get("_ERR")
	}

	// runtime error
	switch x := err.(type) {
	case string:
		return x
	case error:
		return x.Error()
	default:
		return "unknown error"
	}
}

// [-(nargs + 1), +nresults, e]
// http://www.lua.org/manual/5.3/manual.html#lua_callk
func (state *luaState) CallK() {
	panic("todo: CallK!")
}

// [-(nargs + 1), +(nresults|1), –]
// http://www.lua.org/manual/5.3/manual.html#lua_pcallk
func (state *luaState) PCallK() {
	panic("todo: PCallK!")
}
