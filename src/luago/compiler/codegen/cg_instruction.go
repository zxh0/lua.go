package codegen

import . "luago/lua/vm"
import . "luago/number" // todo

// r[a], r[a+1], ..., r[a+b] = nil
func (self *cg) loadNil(line, a, n int) {
	self.inst(line, OP_LOADNIL, a, n-1, 0)
}

// r[a] = (bool)b
func (self *cg) loadBool(line, a, b int) {
	self.inst(line, OP_LOADBOOL, a, b, 0)
}

// r[a] = kst[bx]
func (self *cg) loadK(line, a int, k interface{}) {
	idx := self.indexOf(k)
	if idx-0x100 < 0x100 { // todo
		self.inst(line, OP_LOADK, a, idx, 0)
	} else {
		self.inst(line, OP_LOADKX, a, 0, 0)
		self.inst(line, OP_EXTRAARG, idx-0x100, 0, 0)
	}
}

// r[a] = {}
func (self *cg) newTable(line, a, nArr, nRec int) {
	self.inst(line, OP_NEWTABLE, a, INT2FB(nArr), INT2FB(nRec))
}

// r[a][(c-1)*FPF+i] := r[a+i], 1 <= i <= b
func (self *cg) setList(line, a, b, c int) {
	self.inst(line, OP_SETLIST, a, b, c)
}

// r[a] = closure(proto[bx])
func (self *cg) closure(line, a, bx int) {
	self.inst(line, OP_CLOSURE, a, bx, 0)
}

// r[a] = r[b]
func (self *cg) move(line, a, b int) {
	self.inst(line, OP_MOVE, a, b, 0)
}

// r[a] = upval[b]
func (self *cg) getUpval(line, a, b int) {
	self.inst(line, OP_GETUPVAL, a, b, 0)
}

// upval[b] = r[a]
func (self *cg) setUpval(line, a, b int) {
	self.inst(line, OP_SETUPVAL, a, b, 0)
}

// r[a] = upval[b][rk(c)]
func (self *cg) getTabUp(line, a, b, c int) {
	self.inst(line, OP_GETTABUP, a, b, c)
}

// upval[a][rk(b)] = rk(c)
func (self *cg) setTabUp(line, a, b, c int) {
	self.inst(line, OP_SETTABUP, a, b, c)
}

// r[a] := r[b][rk(c)]
func (self *cg) getTable(line, a, b, c int) {
	self.inst(line, OP_GETTABLE, a, b, c)
}

// r[a][rk(b)] = rk(c)
func (self *cg) setTable(line, a, b, c int) {
	self.inst(line, OP_SETTABLE, a, b, c)
}

// r[a], r[a+1], ..., r[a+b-2] = vararg
func (self *cg) vararg(line, a, n int) {
	self.inst(line, OP_VARARG, a, n+1, 0)
}

// r[a], ..., r[a+c-2] = r[a](r[a+1], ..., r[a+b-1])
func (self *cg) call(line, a, nArgs, nRet int) {
	self.inst(line, OP_CALL, a, nArgs+1, nRet+1)
}

// return r[a](r[a+1], ... ,r[a+b-1])
func (self *cg) tailCall(line, a, nArgs int) {
	self.inst(line, OP_TAILCALL, a, nArgs+1, 0)
}

// return r[a], ... ,r[a+b-2]
func (self *cg) _return(line, a, n int) {
	self.inst(line, OP_RETURN, a, n+1, 0)
}

// r[a+1] := r[b]; r[a] := r[b][rk(c)]
func (self *cg) _self(line, a, b, c int) {
	self.inst(line, OP_SELF, a, b, c)
}

// pc+=sBx; if (a) close all upvalues >= r[a - 1]
func (self *cg) jmp(line, sBx int) int {
	return self.inst(line, OP_JMP, 0, sBx, 0) // todo
}

// if not (r[a] <=> c) then pc++
func (self *cg) test(line, a, c int) {
	self.inst(line, OP_TEST, a, 0, c)
}

// if (r[b] <=> c) then r[a] := r[b] else pc++
func (self *cg) testSet(line, a, b, c int) {
	self.inst(line, OP_TESTSET, a, b, c)
}

// r[a] = op r[b]
func (self *cg) unaryOp(line, op, a, b int) {
	self.inst(line, op, a, b, 0)
}
