package state

func (state *luaState) AddPC(n int) {
	state.stack.pc += n
}

func (state *luaState) Fetch() uint32 {
	i := state.stack.closure.proto.Code[state.stack.pc]
	state.stack.pc++
	return i
}

func (state *luaState) RegisterCount() int {
	return int(state.stack.closure.proto.MaxStackSize)
}

func (state *luaState) GetConst(idx int) {
	c := state.stack.closure.proto.Constants[idx]
	state.stack.push(c)
}

func (state *luaState) GetRK(rk int) {
	if rk > 0xFF { // constant
		state.GetConst(rk & 0xFF)
	} else { // register
		state.PushValue(rk + 1)
	}
}

func (state *luaState) LoadProto(idx int) {
	stack := state.stack
	subProto := stack.closure.proto.Protos[idx]
	closure := newLuaClosure(subProto)

	for i, uvInfo := range subProto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		if uvInfo.Instack == 1 {
			if stack.openuvs == nil {
				stack.openuvs = map[int]*upvalue{}
			}

			if openuv, found := stack.openuvs[uvIdx]; found {
				closure.upvals[i] = openuv
			} else {
				closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
				stack.openuvs[uvIdx] = closure.upvals[i]
			}
		} else {
			closure.upvals[i] = stack.closure.upvals[uvIdx]
		}
	}

	stack.push(closure)
}

func (state *luaState) CloseUpvalues(a int) {
	for i, openuv := range state.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(state.stack.openuvs, i)
		}
	}
}

func (state *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(state.stack.varargs)
	}

	state.stack.check(n)
	state.stack.pushN(state.stack.varargs, n)
}
