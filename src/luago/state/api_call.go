package state

// import "fmt"
import . "luago/api"
import "luago/binchunk"
import "luago/compiler"
import "luago/vm"

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_dump
func (self *luaState) Dump(strip bool) []byte {
	panic("todo!")
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_load
func (self *luaState) Load(chunk []byte, chunkName, mode string) ThreadStatus {
	var proto *binchunk.Prototype
	if binchunk.IsBinaryChunk(chunk) {
		proto = binchunk.Undump(chunk)
	} else {
		proto = compiler.Compile(chunkName, string(chunk))
	}

	c := newLuaClosure(proto)
	if len(proto.Upvalues) > 0 { // todo
		env := self.registry.get(LUA_RIDX_GLOBALS)
		c.upvals[0] = &env
	}
	self.stack.push(c)
	return LUA_OK
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
	calleeStack := newLuaStack(nArgs+LUA_MINSTACK, self)
	calleeStack.closure = c

	// pass args, pop func
	callerStack := self.stack
	if nArgs > 0 {
		args := callerStack.popN(nArgs)
		calleeStack.pushN(args, nArgs)
	}
	callerStack.pop()

	// run closure
	self.pushLuaStack(calleeStack)
	r := c.goFunc(self)
	self.popLuaStack()

	// return results
	if nResults != 0 {
		results := calleeStack.popN(r)
		callerStack.pushN(results, nResults)
	}
}

func (self *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	// create new lua stack
	nRegs := int(c.proto.MaxStackSize)
	calleeStack := newLuaStack(nRegs+LUA_MINSTACK, self)
	calleeStack.closure = c

	// pass args, pop func
	callerStack := self.stack
	if nArgs > 0 {
		args := callerStack.popN(nArgs)
		calleeStack.pushN(args, nArgs)

		nParams := int(c.proto.NumParams)
		isVararg := c.proto.IsVararg == 1
		if nArgs > nParams && isVararg {
			calleeStack.varargs = args[nParams:]
		}
	}
	callerStack.pop()
	calleeStack.top = nRegs

	// run closure
	self.pushLuaStack(calleeStack)
	self.runLuaClosure()
	self.popLuaStack()

	// return results
	if nResults != 0 {
		results := calleeStack.popN(calleeStack.top - nRegs)
		callerStack.pushN(results, nResults)
	}
}

func (self *luaState) runLuaClosure() {
	// fmt.Printf("call %s\n", c.toString())
	code := self.stack.closure.proto.Code
	for {
		pc := self.stack.pc
		inst := vm.Instruction(code[pc])
		self.stack.pc++

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
				status = LUA_ERRRUN
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
