package codegen

import . "luago/compiler/lexer"
import . "luago/lua"
import . "luago/vm"

var arithAndBitwiseBinops = map[int]int{
	TOKEN_OP_ADD:  OP_ADD,
	TOKEN_OP_SUB:  OP_SUB,
	TOKEN_OP_MUL:  OP_MUL,
	TOKEN_OP_MOD:  OP_MOD,
	TOKEN_OP_POW:  OP_POW,
	TOKEN_OP_DIV:  OP_DIV,
	TOKEN_OP_IDIV: OP_IDIV,
	TOKEN_OP_BAND: OP_BAND,
	TOKEN_OP_BOR:  OP_BOR,
	TOKEN_OP_BXOR: OP_BXOR,
	TOKEN_OP_SHL:  OP_SHL,
	TOKEN_OP_SHR:  OP_SHR,
}

// r[a], r[a+1], ..., r[a+b] = nil
func (self *codeGen) emitLoadNil(line, a, n int) {
	self.emit(line, OP_LOADNIL, a, n-1, 0)
}

// r[a] = (bool)b; if (c) pc++
func (self *codeGen) emitLoadBool(line, a, b, c int) {
	self.emit(line, OP_LOADBOOL, a, b, c)
}

// r[a] = kst[bx]
func (self *codeGen) emitLoadK(line, a int, k interface{}) {
	idx := self.indexOfConstant(k)
	if idx-0x100 < 0x100 { // todo
		self.emit(line, OP_LOADK, a, idx, 0)
	} else {
		self.emit(line, OP_LOADKX, a, 0, 0)
		self.emit(line, OP_EXTRAARG, idx-0x100, 0, 0)
	}
}

// r[a] = {}
func (self *codeGen) emitNewTable(line, a, nArr, nRec int) {
	self.emit(line, OP_NEWTABLE, a, INT2FB(nArr), INT2FB(nRec))
}

// r[a][(c-1)*FPF+i] := r[a+i], 1 <= i <= b
func (self *codeGen) emitSetList(line, a, b, c int) {
	self.emit(line, OP_SETLIST, a, b, c)
}

// r[a] = emitClosure(proto[bx])
func (self *codeGen) emitClosure(line, a, bx int) {
	self.emit(line, OP_CLOSURE, a, bx, 0)
}

// r[a] = r[b]
func (self *codeGen) emitMove(line, a, b int) {
	self.emit(line, OP_MOVE, a, b, 0)
}

// r[a] = upval[b]
func (self *codeGen) emitGetUpval(line, a, b int) {
	self.emit(line, OP_GETUPVAL, a, b, 0)
}

// upval[b] = r[a]
func (self *codeGen) emitSetUpval(line, a, b int) {
	self.emit(line, OP_SETUPVAL, a, b, 0)
}

// r[a] = upval[b][rk(c)]
func (self *codeGen) emitGetTabUp(line, a, b, c int) {
	self.emit(line, OP_GETTABUP, a, b, c)
}

// upval[a][rk(b)] = rk(c)
func (self *codeGen) emitSetTabUp(line, a, b, c int) {
	self.emit(line, OP_SETTABUP, a, b, c)
}

// r[a] := r[b][rk(c)]
func (self *codeGen) emitGetTable(line, a, b, c int) {
	self.emit(line, OP_GETTABLE, a, b, c)
}

// r[a][rk(b)] = rk(c)
func (self *codeGen) emitSetTable(line, a, b, c int) {
	self.emit(line, OP_SETTABLE, a, b, c)
}

// r[a], r[a+1], ..., r[a+b-2] = vararg
func (self *codeGen) emitVararg(line, a, n int) {
	self.emit(line, OP_VARARG, a, n+1, 0)
}

// r[a], ..., r[a+c-2] = r[a](r[a+1], ..., r[a+b-1])
func (self *codeGen) emitCall(line, a, nArgs, nRet int) {
	self.emit(line, OP_CALL, a, nArgs+1, nRet+1)
}

// return r[a](r[a+1], ... ,r[a+b-1])
func (self *codeGen) emitTailCall(line, a, nArgs int) {
	self.emit(line, OP_TAILCALL, a, nArgs+1, 0)
}

// return r[a], ... ,r[a+b-2]
func (self *codeGen) emitReturn(line, a, n int) {
	self.emit(line, OP_RETURN, a, n+1, 0)
}

// r[a+1] := r[b]; r[a] := r[b][rk(c)]
func (self *codeGen) emitSelf(line, a, b, c int) {
	self.emit(line, OP_SELF, a, b, c)
}

// pc+=sBx; if (a) close all upvalues >= r[a - 1]
func (self *codeGen) emitJmp(line, sBx int) int {
	return self.emit(line, OP_JMP, 0, sBx, 0) // todo: a?
}

// if not (r[a] <=> c) then pc++
func (self *codeGen) emitTest(line, a, c int) {
	self.emit(line, OP_TEST, a, 0, c)
}

// if (r[b] <=> c) then r[a] := r[b] else pc++
func (self *codeGen) emitTestSet(line, a, b, c int) {
	self.emit(line, OP_TESTSET, a, b, c)
}

func (self *codeGen) emitForPrep(line, a, sBx int) int {
	return self.emit(line, OP_FORPREP, a, sBx, 0)
}

func (self *codeGen) emitForLoop(line, a, sBx int) int {
	return self.emit(line, OP_FORLOOP, a, sBx, 0)
}

func (self *codeGen) emitTForCall(line, a, c int) {
	self.emit(line, OP_TFORCALL, a, 0, c)
}

func (self *codeGen) emitTForLoop(line, a, sBx int) {
	self.emit(line, OP_TFORLOOP, a, sBx, 0)
}

// r[a] = op r[b]
func (self *codeGen) emitUnaryOp(line, op, a, b int) {
	switch op {
	case TOKEN_OP_NOT:
		self.emit(line, OP_NOT, a, b, 0)
	case TOKEN_OP_BNOT:
		self.emit(line, OP_BNOT, a, b, 0)
	case TOKEN_OP_LEN:
		self.emit(line, OP_LEN, a, b, 0)
	case TOKEN_OP_UNM:
		self.emit(line, OP_UNM, a, b, 0)
	}
}

// r[a] = rk[b] op rk[c]
// arith & bitwise & relational
func (self *codeGen) emitBinaryOp(line, op, a, b, c int) {
	if opcode, found := arithAndBitwiseBinops[op]; found {
		self.emit(line, opcode, a, b, c)
	} else {
		switch op {
		case TOKEN_OP_EQ:
			self.emit(line, OP_EQ, 1, b, c)
		case TOKEN_OP_NE:
			self.emit(line, OP_EQ, 0, b, c)
		case TOKEN_OP_LT:
			self.emit(line, OP_LT, 1, b, c)
		case TOKEN_OP_GT:
			self.emit(line, OP_LT, 1, c, b)
		case TOKEN_OP_LE:
			self.emit(line, OP_LE, 1, b, c)
		case TOKEN_OP_GE:
			self.emit(line, OP_LE, 1, c, b)
		}
		self.emitLoadBool(line, a, 1, 1)
		self.emitLoadBool(line, a, 0, 0)
	}
}
