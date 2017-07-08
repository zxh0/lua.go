package state

import . "luago/api"
import "luago/vm"

// [-(nargs+1), +nresults, e]
// http://www.lua.org/manual/5.3/manual.html#lua_call
func (self *luaState) Call(nArgs, nResults int) {
	x := self.stack.get(-(nArgs + 1))

	switch f := x.(type) {
	case GoFunction: // todo
		self.callGoClosure(nArgs, nResults, &goClosure{goFunc: f})
	case *goClosure:
		self.callGoClosure(nArgs, nResults, f)
	case *luaClosure:
		self.callLuaClosure(nArgs, nResults, f)
	default:
		panic("not a function!")
	}
}

func (self *luaState) callGoClosure(nArgs, nResults int, f *goClosure) {
	// pop args & func
	callerStack := self.stack
	args := callerStack.popN(nArgs)
	callerStack.pop()

	// create new lua stack
	calleeStack := newLuaStack(nArgs+LUA_MINSTACK, 0, self)
	calleeStack.goCl = f

	// pass args
	calleeStack.check(nArgs)
	calleeStack.pushN(args)

	// call func
	self.pushLuaStack(calleeStack)
	r := f.goFunc(self)
	self.popLuaStack()

	// return results
	if nResults != 0 {
		results := calleeStack.popN(r)
		_pushResults(nResults, results, callerStack)
	}
}

func (self *luaState) callLuaClosure(nArgs, nResults int, f *luaClosure) {
	// pop args & func
	callerStack := self.stack
	args := callerStack.popN(nArgs)
	callerStack.pop()

	// create new lua stack
	nRegs := int(f.proto.MaxStackSize)
	calleeStack := newLuaStack(nRegs+4, nRegs, self)
	calleeStack.luaCl = f

	// pass args
	nParams := int(f.proto.NumParams)
	for i, arg := range args {
		if i < nParams {
			calleeStack.slots[i] = arg
		}
	}
	if f.proto.IsVararg == 1 { // TODO
		if nArgs > nParams {
			calleeStack.xArgs = args[nParams:]
		}
	}

	// call func
	self.pushLuaStack(calleeStack)
	self.runLuaClosure(f)
	self.popLuaStack()

	// return results
	if nResults != 0 {
		results := calleeStack.popAll()
		_pushResults(nResults, results, callerStack)
	}
}

func (self *luaState) runLuaClosure(f *luaClosure) {
	// fmt.Printf("call %s\n", f.toString())
	code := f.proto.Code
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

func _pushResults(nResults int, results []luaValue, stack *luaStack) {
	if nResults > 0 && nResults < len(results) {
		results = results[0:nResults]
	}
	if nResults > 0 {
		stack.check(nResults)
	} else {
		stack.check(len(results))
	}
	stack.pushN(results)
	if nResults > len(results) {
		for i := len(results); i < nResults; i++ {
			stack.push(nil)
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
