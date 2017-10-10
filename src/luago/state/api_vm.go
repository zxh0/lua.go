package state

func (self *luaState) AddPC(n int) {
	self.stack.pc += n
}

func (self *luaState) Instruction() uint32 {
	return self.stack.closure.proto.Code[self.stack.pc]
}

// todo
func (self *luaState) MaxStackSize() int {
	return int(self.stack.closure.proto.MaxStackSize)
}

func (self *luaState) GetConst(idx int) {
	c := self.stack.closure.proto.Constants[idx]
	self.stack.push(c)
}

func (self *luaState) GetRK(rk int) {
	if rk > 0xFF { // constant
		self.GetConst(rk & 0xFF)
	} else { // register
		self.PushValue(rk + 1)
	}
}

func (self *luaState) GetUpvalue2(idx int) {
	upval := self.stack.closure.upvals[idx]
	self.stack.push(*upval)
}

func (self *luaState) SetUpvalue2(idx int) {
	upval := self.stack.closure.upvals[idx]
	*upval = self.stack.pop()
}

func (self *luaState) LoadProto(idx int) {
	proto := self.stack.closure.proto.Protos[idx]
	closure := newLuaClosure(proto)

	// todo
	for i, uvInfo := range proto.Upvalues {
		if uvInfo.Instack == 1 {
			closure.upvals[i] = &(self.stack.slots[uvInfo.Idx])
		} else {
			closure.upvals[i] = self.stack.closure.upvals[uvInfo.Idx]
		}
	}

	self.stack.push(closure)
}

func (self *luaState) LoadVararg(n int) {
	stack := self.stack
	vargs := stack.varargs

	if n < 0 {
		n = len(vargs)
	}

	stack.check(n)
	for i := 0; i < n; i++ {
		if i < len(vargs) {
			stack.push(vargs[i])
		} else {
			stack.push(nil)
		}
	}
}
