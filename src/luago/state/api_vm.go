package state

func (self *luaState) AddPC(n int) {
	self.stack.pc += n
}

// todo
func (self *luaState) MaxStackSize() int {
	return int(self.stack.luaCl.proto.MaxStackSize)
}

func (self *luaState) GetConst(index int) {
	c := self.stack.luaCl.proto.Constants[index]
	self.stack.push(c)
}

func (self *luaState) GetRK(rk int) {
	if rk > 0xFF { // constant
		self.GetConst(rk & 0xFF)
	} else { // register
		self.PushValue(rk + 1)
	}
}

func (self *luaState) GetUpvalue2(index int) {
	upval := self.stack.luaCl.upvals[index]
	self.stack.push(*upval)
}

func (self *luaState) SetUpvalue2(index int) {
	upval := self.stack.luaCl.upvals[index]
	*upval = self.stack.pop()
}

func (self *luaState) LoadProto(index int) {
	proto := self.stack.luaCl.proto.Protos[index]
	closure := newLuaClosure(proto)

	// todo
	for i, uvInfo := range proto.Upvalues {
		if uvInfo.Instack == 1 {
			closure.upvals[i] = &(self.stack.slots[uvInfo.Idx])
		} else {
			closure.upvals[i] = self.stack.luaCl.upvals[uvInfo.Idx]
		}
	}

	self.stack.push(closure)
}

func (self *luaState) LoadVararg(n int) {
	stack := self.stack
	xArgs := stack.xArgs

	if n < 0 {
		n = len(xArgs)
	}

	stack.check(n)
	for i := 0; i < n; i++ {
		if i < len(xArgs) {
			stack.push(xArgs[i])
		} else {
			stack.push(nil)
		}
	}
}
