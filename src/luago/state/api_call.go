package state

import . "luago/api"
import "luago/binchunk"
import "luago/compiler"
import "luago/vm"

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_load
func (self *luaState) Load(chunk []byte, chunkName, mode string) ThreadStatus {
	var proto *binchunk.Prototype
	if binchunk.IsBinaryChunk(chunk) {
		proto = binchunk.Undump(chunk)
	} else {
		proto = compiler.Compile(chunkName, string(chunk))
	}

	cl := newLuaClosure(proto)
	if len(proto.Upvalues) > 0 { // todo
		env := self.registry.get(LUA_RIDX_GLOBALS)
		cl.upvals[0] = &env
	}
	self.stack.push(cl)
	return LUA_OK
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_dump
func (self *luaState) Dump() {
	panic("todo!")
}

// [-(nargs+1), +nresults, e]
// http://www.lua.org/manual/5.3/manual.html#lua_call
func (self *luaState) Call(nArgs, nResults int) {
	val := self.stack.get(-(nArgs + 1))

	switch x := val.(type) {
	case *luaClosure:
		self.callLuaClosure(nArgs, nResults, x)
	case *goClosure:
		self.callGoClosure(nArgs, nResults, x)
	case GoFunction: // todo
		self.callGoClosure(nArgs, nResults, &goClosure{goFunc: x})
	default:
		panic("not a function!")
	}
}

func (self *luaState) callGoClosure(nArgs, nResults int, c *goClosure) {
	// create new lua stack
	calleeStack := newLuaStack(nArgs+LUA_MINSTACK, self)
	calleeStack.goCl = c

	// pass args, pop func
	callerStack := self.stack
	if nArgs > 0 {
		args := callerStack.popN(nArgs)
		calleeStack.pushN(args)
	}
	callerStack.pop()

	// call func
	self.pushLuaStack(calleeStack)
	r := c.goFunc(self)
	self.popLuaStack()

	// return results
	if nResults != 0 {
		results := calleeStack.popN(r)
		_getResults(callerStack, results, nResults)
	}
}

func (self *luaState) callLuaClosure(nArgs, nResults int, c *luaClosure) {
	// create new lua stack
	nRegs := int(c.proto.MaxStackSize)
	calleeStack := newLuaStack(nRegs+LUA_MINSTACK, self)
	calleeStack.top = nRegs
	calleeStack.luaCl = c

	// pass args, pop func
	callerStack := self.stack
	if nArgs > 0 {
		args := callerStack.popN(nArgs)
		_passArgs(calleeStack, args, c)
	}
	callerStack.pop()

	// call func
	self.pushLuaStack(calleeStack)
	self.runLuaClosure(c)
	self.popLuaStack()

	// return results
	if nResults != 0 {
		results := calleeStack.popN(calleeStack.top - nRegs)
		_getResults(callerStack, results, nResults)
	}
}

func (self *luaState) runLuaClosure(c *luaClosure) {
	// fmt.Printf("call %s\n", c.toString())
	code := c.proto.Code
	for {
		pc := self.stack.pc
		inst := vm.Instruction(code[pc])
		self.stack.pc++

		inst.Execute(self)

		// indent := fmt.Sprintf("%%%ds", ls.callDepth*2)
		// fmt.Printf(indent+"[%02d: %s] => %s\n",
		// 	"", pc+1, inst.OpName(), ls)

		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

func _passArgs(calleeStack *luaStack, args []luaValue, c *luaClosure) {
	nParams := int(c.proto.NumParams)
	for i, arg := range args {
		if i < nParams {
			calleeStack.slots[i] = arg
		}
	}
	if len(args) > nParams && c.proto.IsVararg == 1 {
		calleeStack.xArgs = args[nParams:]
	}
}

func _getResults(callerStack *luaStack, results []luaValue, nResults int) {
	if nResults < 0 || nResults == len(results) {
		callerStack.pushN(results)
	} else if nResults < len(results) {
		callerStack.pushN(results[0:nResults])
	} else { // nResults > len(results)
		callerStack.pushN(results)
		for i := len(results); i < nResults; i++ {
			callerStack.push(nil)
		}
	}
}

// Calls a function in protected mode.
// http://www.lua.org/manual/5.3/manual.html#lua_pcall
func (self *luaState) PCall(nArgs, nResults, msgh int) (status ThreadStatus) {
	callDepth := self.callDepth

	// catch error
	defer func() {
		if r := recover(); r != nil { // todo
			if msgh < 0 {
				panic(r)
			} else if msgh > 0 {
				panic("todo: msgh > 0")
			} else {
				for self.callDepth > callDepth {
					self.popLuaStack()
				}
				self.stack.push(_getErrMsg(r))
				status = LUA_ERRRUN
			}
		}
	}()

	self.Call(nArgs, nResults)
	status = LUA_OK
	return
}

func _getErrMsg(err interface{}) string {
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
