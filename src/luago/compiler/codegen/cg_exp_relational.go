package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"
import . "luago/vm"

func (self *codeGen) testRelationalBinopExpX(exp *BinopExp, allocator *tmpAllocator, a int) {
	rkb, _ := self.exp2OpArg(exp.Exp1, ARG_RK, allocator)
	rkc, _ := self.exp2OpArg(exp.Exp2, ARG_RK, allocator)

	switch exp.Op {
	case TOKEN_OP_EQ:
		self.emit(exp.Line, OP_EQ, a, rkb, rkc)
	case TOKEN_OP_NE:
		self.emit(exp.Line, OP_EQ, 1-a, rkb, rkc)
	case TOKEN_OP_LT:
		self.emit(exp.Line, OP_LT, a, rkb, rkc)
	case TOKEN_OP_GT:
		self.emit(exp.Line, OP_LT, a, rkc, rkb)
	case TOKEN_OP_LE:
		self.emit(exp.Line, OP_LE, a, rkb, rkc)
	case TOKEN_OP_GE:
		self.emit(exp.Line, OP_LE, a, rkc, rkb)
	}
}

func (self *codeGen) testRelationalBinopExp(exp *BinopExp, a int) {
	allocator := self.newTmpAllocator(-1)
	rkb, _ := self.exp2OpArg(exp.Exp1, ARG_RK, allocator)
	rkc, _ := self.exp2OpArg(exp.Exp2, ARG_RK, allocator)
	allocator.freeAll()

	switch exp.Op {
	case TOKEN_OP_EQ:
		self.emit(exp.Line, OP_EQ, a, rkb, rkc)
	case TOKEN_OP_NE:
		self.emit(exp.Line, OP_EQ, 1-a, rkb, rkc)
	case TOKEN_OP_LT:
		self.emit(exp.Line, OP_LT, a, rkb, rkc)
	case TOKEN_OP_GT:
		self.emit(exp.Line, OP_LT, a, rkc, rkb)
	case TOKEN_OP_LE:
		self.emit(exp.Line, OP_LE, a, rkb, rkc)
	case TOKEN_OP_GE:
		self.emit(exp.Line, OP_LE, a, rkc, rkb)
	}
}
