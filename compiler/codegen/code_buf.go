package codegen

import (
	. "github.com/zxh0/lua.go/compiler/lexer"
	. "github.com/zxh0/lua.go/vm"
)

type codeBuf struct {
	insts    []uint32
	lineNums []uint32
}

func (cb *codeBuf) pc() int {
	return len(cb.insts) - 1
}

func (cb *codeBuf) fixSbx(pc, sBx int) {
	if sBx > 0 && sBx > MAXARG_sBx || sBx < 0 && -sBx > MAXARG_sBx {
		panic("control structure too long")
	}

	i := cb.insts[pc]
	i = i << 18 >> 18                  // clear sBx
	i = i | uint32(sBx+MAXARG_sBx)<<14 // reset sBx
	cb.insts[pc] = i
}

func (cb *codeBuf) emitABC(line, opcode, a, b, c int) {
	i := b<<23 | c<<14 | a<<6 | opcode
	cb.insts = append(cb.insts, uint32(i))
	cb.lineNums = append(cb.lineNums, uint32(line))
}

func (cb *codeBuf) emitABx(line, opcode, a, bx int) {
	i := bx<<14 | a<<6 | opcode
	cb.insts = append(cb.insts, uint32(i))
	cb.lineNums = append(cb.lineNums, uint32(line))
}

func (cb *codeBuf) emitAsBx(line, opcode, a, sBx int) {
	i := (sBx+MAXARG_sBx)<<14 | a<<6 | opcode
	cb.insts = append(cb.insts, uint32(i))
	cb.lineNums = append(cb.lineNums, uint32(line))
}

func (cb *codeBuf) emitAx(line, opcode, ax int) {
	i := ax<<6 | opcode
	cb.insts = append(cb.insts, uint32(i))
	cb.lineNums = append(cb.lineNums, uint32(line))
}

// r[a] = r[b]
func (cb *codeBuf) emitMove(line, a, b int) {
	cb.emitABC(line, OP_MOVE, a, b, 0)
}

// r[a], r[a+1], ..., r[a+b] = nil
func (cb *codeBuf) emitLoadNil(line, a, n int) {
	cb.emitABC(line, OP_LOADNIL, a, n-1, 0)
}

// r[a] = (bool)b; if (c) pc++
func (cb *codeBuf) emitLoadBool(line, a, b, c int) {
	//cb.emitABC(line, OP_LOADBOOL, a, b, c)
	panic("TODO")
}

// r[a] = kst[bx]
func (cb *codeBuf) emitLoadK(line, a int, idx int) {
	if idx < (1 << 18) {
		cb.emitABx(line, OP_LOADK, a, idx)
	} else {
		cb.emitABx(line, OP_LOADKX, a, 0)
		cb.emitAx(line, OP_EXTRAARG, idx)
	}
}

// r[a], r[a+1], ..., r[a+b-2] = vararg
func (cb *codeBuf) emitVararg(line, a, n int) {
	cb.emitABC(line, OP_VARARG, a, n+1, 0)
}

// r[a] = emitClosure(proto[bx])
func (cb *codeBuf) emitClosure(line, a, bx int) {
	cb.emitABx(line, OP_CLOSURE, a, bx)
}

// r[a] = {}
func (cb *codeBuf) emitNewTable(line, a, nArr, nRec int) {
	cb.emitABC(line, OP_NEWTABLE,
		a, Int2fb(nArr), Int2fb(nRec))
}

// r[a][(c-1)*FPF+i] := r[a+i], 1 <= i <= b
func (cb *codeBuf) emitSetList(line, a, b, c int) {
	cb.emitABC(line, OP_SETLIST, a, b, c)
}

// r[a] := r[b][rk(c)]
func (cb *codeBuf) emitGetTable(line, a, b, c int) {
	cb.emitABC(line, OP_GETTABLE, a, b, c)
}

// r[a][rk(b)] = rk(c)
func (cb *codeBuf) emitSetTable(line, a, b, c int) {
	cb.emitABC(line, OP_SETTABLE, a, b, c)
}

// r[a] = upval[b]
func (cb *codeBuf) emitGetUpval(line, a, b int) {
	cb.emitABC(line, OP_GETUPVAL, a, b, 0)
}

// upval[b] = r[a]
func (cb *codeBuf) emitSetUpval(line, a, b int) {
	cb.emitABC(line, OP_SETUPVAL, a, b, 0)
}

// r[a] = upval[b][rk(c)]
func (cb *codeBuf) emitGetTabUp(line, a, b, c int) {
	cb.emitABC(line, OP_GETTABUP, a, b, c)
}

// upval[a][rk(b)] = rk(c)
func (cb *codeBuf) emitSetTabUp(line, a, b, c int) {
	cb.emitABC(line, OP_SETTABUP, a, b, c)
}

// r[a], ..., r[a+c-2] = r[a](r[a+1], ..., r[a+b-1])
func (cb *codeBuf) emitCall(line, a, nArgs, nRet int) {
	cb.emitABC(line, OP_CALL, a, nArgs+1, nRet+1)
}

// return r[a](r[a+1], ... ,r[a+b-1])
func (cb *codeBuf) emitTailCall(line, a, nArgs int) {
	cb.emitABC(line, OP_TAILCALL, a, nArgs+1, 0)
}

// return r[a], ... ,r[a+b-2]
func (cb *codeBuf) emitReturn(line, a, n int) {
	cb.emitABC(line, OP_RETURN, a, n+1, 0)
}

// r[a+1] := r[b]; r[a] := r[b][rk(c)]
func (cb *codeBuf) emitSelf(line, a, b, c int) {
	cb.emitABC(line, OP_SELF, a, b, c)
}

// pc+=sBx; if (a) close all upvalues >= r[a - 1]
func (cb *codeBuf) emitJmp(line, a, sBx int) int {
	cb.emitAsBx(line, OP_JMP, a, sBx)
	return len(cb.insts) - 1
}

// if not (r[a] <=> c) then pc++
func (cb *codeBuf) emitTest(line, a, c int) {
	cb.emitABC(line, OP_TEST, a, 0, c)
}

// if (r[b] <=> c) then r[a] := r[b] else pc++
func (cb *codeBuf) emitTestSet(line, a, b, c int) {
	cb.emitABC(line, OP_TESTSET, a, b, c)
}

func (cb *codeBuf) emitForPrep(line, a, sBx int) int {
	cb.emitAsBx(line, OP_FORPREP, a, sBx)
	return len(cb.insts) - 1
}

func (cb *codeBuf) emitForLoop(line, a, sBx int) int {
	cb.emitAsBx(line, OP_FORLOOP, a, sBx)
	return len(cb.insts) - 1
}

func (cb *codeBuf) emitTForCall(line, a, c int) {
	cb.emitABC(line, OP_TFORCALL, a, 0, c)
}

func (cb *codeBuf) emitTForLoop(line, a, sBx int) {
	cb.emitAsBx(line, OP_TFORLOOP, a, sBx)
}

// r[a] = op r[b]
func (cb *codeBuf) emitUnaryOp(line, op, a, b int) {
	switch op {
	case TOKEN_OP_NOT:
		cb.emitABC(line, OP_NOT, a, b, 0)
	case TOKEN_OP_BNOT:
		cb.emitABC(line, OP_BNOT, a, b, 0)
	case TOKEN_OP_LEN:
		cb.emitABC(line, OP_LEN, a, b, 0)
	case TOKEN_OP_UNM:
		cb.emitABC(line, OP_UNM, a, b, 0)
	}
}

// r[a] = rk[b] op rk[c]
// arith & bitwise & relational
func (cb *codeBuf) emitBinaryOp(line, op, a, b, c int) {
	if opcode, found := arithAndBitwiseBinops[op]; found {
		cb.emitABC(line, opcode, a, b, c)
	} else {
		switch op {
		case TOKEN_OP_EQ:
			cb.emitABC(line, OP_EQ, 1, b, c)
		case TOKEN_OP_NE:
			cb.emitABC(line, OP_EQ, 0, b, c)
		case TOKEN_OP_LT:
			cb.emitABC(line, OP_LT, 1, b, c)
		case TOKEN_OP_GT:
			cb.emitABC(line, OP_LT, 1, c, b)
		case TOKEN_OP_LE:
			cb.emitABC(line, OP_LE, 1, b, c)
		case TOKEN_OP_GE:
			cb.emitABC(line, OP_LE, 1, c, b)
		}
		cb.emitJmp(line, 0, 1)
		cb.emitLoadBool(line, a, 0, 1)
		cb.emitLoadBool(line, a, 1, 0)
	}
}
