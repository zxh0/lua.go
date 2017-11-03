package codegen

import . "luago/compiler/lexer"
import "luago/number"
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

func (self *codeGen) pc() int {
	return len(self.insts) - 1
}

func (self *codeGen) emitABC(line, opcode, a, b, c int) {
	i := b<<23 | c<<14 | a<<6 | opcode
	self.insts = append(self.insts, uint32(i))
	self.lines = append(self.lines, uint32(line))
}

func (self *codeGen) emitABx(line, opcode, a, bx int) {
	i := bx<<14 | a<<6 | opcode
	self.insts = append(self.insts, uint32(i))
	self.lines = append(self.lines, uint32(line))
}

func (self *codeGen) emitAsBx(line, opcode, a, b, c int) {
	i := (b+MAXARG_sBx)<<14 | a<<6 | opcode
	self.insts = append(self.insts, uint32(i))
	self.lines = append(self.lines, uint32(line))
}

func (self *codeGen) emitAx(line, opcode, ax int) {
	i := ax<<6 | opcode
	self.insts = append(self.insts, uint32(i))
	self.lines = append(self.lines, uint32(line))
}

func (self *codeGen) fixSbx(pc, sBx int) {
	i := self.insts[pc]
	i = i << 18 >> 18                  // clear sBx
	i = i | uint32(sBx+MAXARG_sBx)<<14 // reset sBx
	self.insts[pc] = i
}

// r[a] = r[b]
func (self *codeGen) emitMove(line, a, b int) {
	self.emitABC(line, OP_MOVE, a, b, 0)
}

// r[a], r[a+1], ..., r[a+b] = nil
func (self *codeGen) emitLoadNil(line, a, n int) {
	self.emitABC(line, OP_LOADNIL, a, n-1, 0)
}

// r[a] = (bool)b; if (c) pc++
func (self *codeGen) emitLoadBool(line, a, b, c int) {
	self.emitABC(line, OP_LOADBOOL, a, b, c)
}

// r[a] = kst[bx]
func (self *codeGen) emitLoadK(line, a int, k interface{}) {
	idx := self.indexOfConstant(k)
	if idx < (1 << 18) {
		self.emitABx(line, OP_LOADK, a, idx)
	} else {
		self.emitABx(line, OP_LOADKX, a, 0)
		self.emitAx(line, OP_EXTRAARG, idx)
	}
}

// r[a], r[a+1], ..., r[a+b-2] = vararg
func (self *codeGen) emitVararg(line, a, n int) {
	self.emitABC(line, OP_VARARG, a, n+1, 0)
}

// r[a] = emitClosure(proto[bx])
func (self *codeGen) emitClosure(line, a, bx int) {
	self.emitABx(line, OP_CLOSURE, a, bx)
}

// r[a] = {}
func (self *codeGen) emitNewTable(line, a, nArr, nRec int) {
	self.emitABC(line, OP_NEWTABLE,
		a, number.Int2fb(nArr), number.Int2fb(nRec))
}

// r[a][(c-1)*FPF+i] := r[a+i], 1 <= i <= b
func (self *codeGen) emitSetList(line, a, b, c int) {
	self.emitABC(line, OP_SETLIST, a, b, c)
}

// r[a] := r[b][rk(c)]
func (self *codeGen) emitGetTable(line, a, b, c int) {
	self.emitABC(line, OP_GETTABLE, a, b, c)
}

// r[a][rk(b)] = rk(c)
func (self *codeGen) emitSetTable(line, a, b, c int) {
	self.emitABC(line, OP_SETTABLE, a, b, c)
}

// r[a] = upval[b]
func (self *codeGen) emitGetUpval(line, a, b int) {
	self.emitABC(line, OP_GETUPVAL, a, b, 0)
}

// upval[b] = r[a]
func (self *codeGen) emitSetUpval(line, a, b int) {
	self.emitABC(line, OP_SETUPVAL, a, b, 0)
}

// r[a] = upval[b][rk(c)]
func (self *codeGen) emitGetTabUp(line, a, b, c int) {
	self.emitABC(line, OP_GETTABUP, a, b, c)
}

// upval[a][rk(b)] = rk(c)
func (self *codeGen) emitSetTabUp(line, a, b, c int) {
	self.emitABC(line, OP_SETTABUP, a, b, c)
}

// r[a], ..., r[a+c-2] = r[a](r[a+1], ..., r[a+b-1])
func (self *codeGen) emitCall(line, a, nArgs, nRet int) {
	self.emitABC(line, OP_CALL, a, nArgs+1, nRet+1)
}

// return r[a](r[a+1], ... ,r[a+b-1])
func (self *codeGen) emitTailCall(line, a, nArgs int) {
	self.emitABC(line, OP_TAILCALL, a, nArgs+1, 0)
}

// return r[a], ... ,r[a+b-2]
func (self *codeGen) emitReturn(line, a, n int) {
	self.emitABC(line, OP_RETURN, a, n+1, 0)
}

// r[a+1] := r[b]; r[a] := r[b][rk(c)]
func (self *codeGen) emitSelf(line, a, b, c int) {
	self.emitABC(line, OP_SELF, a, b, c)
}

// pc+=sBx; if (a) close all upvalues >= r[a - 1]
func (self *codeGen) emitJmp(line, sBx int) int {
	self.emitAsBx(line, OP_JMP, 0, sBx, 0) // todo: a?
	return len(self.insts) - 1
}

// if not (r[a] <=> c) then pc++
func (self *codeGen) emitTest(line, a, c int) {
	self.emitABC(line, OP_TEST, a, 0, c)
}

// if (r[b] <=> c) then r[a] := r[b] else pc++
func (self *codeGen) emitTestSet(line, a, b, c int) {
	self.emitABC(line, OP_TESTSET, a, b, c)
}

func (self *codeGen) emitForPrep(line, a, sBx int) int {
	self.emitAsBx(line, OP_FORPREP, a, sBx, 0)
	return len(self.insts) - 1
}

func (self *codeGen) emitForLoop(line, a, sBx int) int {
	self.emitAsBx(line, OP_FORLOOP, a, sBx, 0)
	return len(self.insts) - 1
}

func (self *codeGen) emitTForCall(line, a, c int) {
	self.emitABC(line, OP_TFORCALL, a, 0, c)
}

func (self *codeGen) emitTForLoop(line, a, sBx int) {
	self.emitAsBx(line, OP_TFORLOOP, a, sBx, 0)
}

// r[a] = op r[b]
func (self *codeGen) emitUnaryOp(line, op, a, b int) {
	switch op {
	case TOKEN_OP_NOT:
		self.emitABC(line, OP_NOT, a, b, 0)
	case TOKEN_OP_BNOT:
		self.emitABC(line, OP_BNOT, a, b, 0)
	case TOKEN_OP_LEN:
		self.emitABC(line, OP_LEN, a, b, 0)
	case TOKEN_OP_UNM:
		self.emitABC(line, OP_UNM, a, b, 0)
	}
}

// r[a] = rk[b] op rk[c]
// arith & bitwise & relational
func (self *codeGen) emitBinaryOp(line, op, a, b, c int) {
	if opcode, found := arithAndBitwiseBinops[op]; found {
		self.emitABC(line, opcode, a, b, c)
	} else {
		switch op {
		case TOKEN_OP_EQ:
			self.emitABC(line, OP_EQ, 1, b, c)
		case TOKEN_OP_NE:
			self.emitABC(line, OP_EQ, 0, b, c)
		case TOKEN_OP_LT:
			self.emitABC(line, OP_LT, 1, b, c)
		case TOKEN_OP_GT:
			self.emitABC(line, OP_LT, 1, c, b)
		case TOKEN_OP_LE:
			self.emitABC(line, OP_LE, 1, b, c)
		case TOKEN_OP_GE:
			self.emitABC(line, OP_LE, 1, c, b)
		}
		self.emitJmp(line, 1)
		self.emitLoadBool(line, a, 0, 1)
		self.emitLoadBool(line, a, 1, 0)
	}
}
