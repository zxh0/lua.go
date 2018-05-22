package state

// import "fmt"
import . "luago/api"
import "luago/binchunk"
import "luago/compiler"
import "luago/vm"

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_dump
func (self *luaState) Dump(strip bool) []byte {
	v := self.stack.get(-1)
	if c, ok := v.(*closure); ok {
		return binchunk.Dump(c.proto)
	}
	panic("todo!")
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_load
func (self *luaState) Load(chunk []byte, chunkName, mode string) (status ThreadStatus) {
	status = LUA_ERRSYNTAX

	// catch error
	defer func() {
		if r := recover(); r != nil {
			self.stack.push(_getErrObj(r))
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
		env := self.registry.get(LUA_RIDX_GLOBALS)
		c.upvals[0] = &upvalue{&env}
	}
	self.stack.push(c)
	status = LUA_OK
	return
}

// [-(nargs+1), +nresults, e]
// http://www.lua.org/manual/5.3/manual.html#lua_call
func (self *luaState) Call(nArgs, nResults int) {
	val := self.stack.get(-(nArgs + 1))

	c, ok := val.(*closure)
	if !ok {
		if mf := getMetafield(val, "__call", self); mf != nil {
			if c, ok = mf.(*closure); ok {
				self.stack.push(val)
				self.Insert(-(nArgs + 2))
				nArgs += 1
			}
		}
	}

	if ok {
		if c.proto != nil {
			self.callLuaClosure(nArgs, nResults, c)
		} else {
			self.callGoClosure(nArgs, nResults, c)
		}
	} else {
		typeName := self.TypeName(typeOf(val))
		panic("attempt to call a " + typeName + " value")
	}
}

func (self *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	// create new lua stack
	newStack := newLuaStack(nArgs+LUA_MINSTACK, self)
	newStack.closure = c

	// pass args, pop func
	if nArgs > 0 {
		args := self.stack.popN(nArgs)
		newStack.pushN(args, nArgs)
	}
	self.stack.pop()

	// run closure
	self.pushLuaStack(newStack)
	r := c.goFunc(self)
	self.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(r)
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}
}

func (self *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// create new lua stack
	newStack := newLuaStack(nRegs+LUA_MINSTACK, self)
	newStack.closure = c

	// pass args, pop func
	funcAndArgs := self.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	// run closure
	self.pushLuaStack(newStack)
	self.runLuaClosure()
	self.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}
}

func (self *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(self.Fetch())
		inst.Execute(self)

		// indent := fmt.Sprintf("%%%ds", self.callDepth*2)
		// fmt.Printf(indent+"[%02d: %s] => %s\n",
		// 	"", pc+1, inst.OpName(), self)

		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

// Calls a function in protected mode.
// http://www.lua.org/manual/5.3/manual.html#lua_pcall
func (self *luaState) PCall(nArgs, nResults, msgh int) (status ThreadStatus) {
	status = LUA_ERRRUN
	caller := self.stack

	// catch error
	defer func() {
		if r := recover(); r != nil { // todo
			if msgh < 0 {
				panic(_getErrObj(r))
			} else if msgh > 0 {
				panic("todo: msgh > 0")
			} else {
				for self.stack != caller {
					self.popLuaStack()
				}
				self.stack.push(_getErrObj(r))
			}
		}
	}()

	self.Call(nArgs, nResults)
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
func (self *luaState) CallK() {
	panic("todo: CallK!")
}

// [-(nargs + 1), +(nresults|1), –]
// http://www.lua.org/manual/5.3/manual.html#lua_pcallk
func (self *luaState) PCallK() {
	panic("todo: PCallK!")
}
