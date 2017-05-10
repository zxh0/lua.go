package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"
import . "luago/lua/vm"

func (self *cg) testExp(node Exp, a int) {
	if bexp, ok := node.(*BinopExp); ok {
		switch bexp.Op {
		case TOKEN_OP_EQ, TOKEN_OP_NE,
			TOKEN_OP_LT, TOKEN_OP_GT,
			TOKEN_OP_LE, TOKEN_OP_GE:
			self.testRelationalBinopExp(bexp, a)
			return
		}
	}

	self.exp(node, a, 1)
}

func (self *cg) testRelationalBinopExp(exp *BinopExp, a int) {
	iX, tX := self.toOpArg(exp.Exp1)
	if tX != ARG_REG && tX != ARG_CONST {
		iX = a
		self.exp(exp.Exp1, iX, 1)
	}

	iY, tY := self.toOpArg(exp.Exp2)
	if tY != ARG_REG && tY != ARG_CONST {
		if tX != ARG_REG && tX != ARG_CONST {
			iY = self.allocTmp()
			self.exp(exp.Exp2, iY, 1)
			self.freeTmp()
		} else {
			iY = a
			self.exp(exp.Exp2, iY, 1)
		}
	}

	switch exp.Op {
	case TOKEN_OP_EQ:
		self.inst(exp.Line, OP_EQ, 0, iX, iY)
	case TOKEN_OP_NE:
		self.inst(exp.Line, OP_EQ, 1, iX, iY)
	case TOKEN_OP_LT:
		self.inst(exp.Line, OP_LT, 0, iX, iY)
	case TOKEN_OP_GT:
		self.inst(exp.Line, OP_LT, 0, iY, iX)
	case TOKEN_OP_LE:
		self.inst(exp.Line, OP_LE, 0, iX, iY)
	case TOKEN_OP_GE:
		self.inst(exp.Line, OP_LE, 0, iY, iX)
	}
}

func (self *cg) testLogicalAndExp(exp *BinopExp, lastJmpLine int) []int {
	jmps := make([]int, 0, 4)

	for {
		if slot, ok := self.isLocVar(exp.Exp1); ok {
			self.test(exp.Line, slot, 0)
		} else {
			if !isRelationalBinopExp(exp.Exp1) {
				tmp := self.allocTmp()
				self.testExp(exp.Exp1, tmp)
				self.freeTmp()
				self.test(exp.Line, tmp, 0)
			} else {
				self.testExp(exp.Exp1, 0) // todo
			}
		}
		jmps = append(jmps, self.jmp(exp.Line, 0))

		if exp2, ok := exp.Exp2.(*BinopExp); ok && exp2.Op == exp.Op {
			exp = exp2
		} else {
			break
		}
	}

	if slot, ok := self.isLocVar(exp.Exp2); ok {
		self.test(lastJmpLine, slot, 0)
	} else {

		if !isRelationalBinopExp(exp.Exp2) {
			tmp := self.allocTmp()
			self.testExp(exp.Exp2, tmp)
			self.freeTmp()
			self.test(lastJmpLine, tmp, 0)
		} else {
			self.testExp(exp.Exp2, 0) // todo
		}
	}
	jmps = append(jmps, self.jmp(lastJmpLine, 0))

	return jmps
}

func (self *cg) testLogicalOrExp(exp *BinopExp, lastJmpLine int) int {
	jmps := make([]int, 0, 4)

	for {
		if slot, ok := self.isLocVar(exp.Exp1); ok {
			self.test(exp.Line, slot, 1)
		} else {
			tmp := self.allocTmp()
			self.testExp(exp.Exp1, tmp)
			self.freeTmp()
			self.test(exp.Line, tmp, 1)
		}
		jmps = append(jmps, self.jmp(exp.Line, 0))

		if exp2, ok := exp.Exp2.(*BinopExp); ok && exp2.Op == exp.Op {
			exp = exp2
		} else {
			break
		}
	}

	if slot, ok := self.isLocVar(exp.Exp2); ok {
		self.test(lastJmpLine, slot, 0)
	} else {
		tmp := self.allocTmp()
		self.testExp(exp.Exp2, tmp)
		self.freeTmp()
		self.test(lastJmpLine, tmp, 0)
	}
	lastJmp := self.jmp(lastJmpLine, 0)

	for _, pc := range jmps {
		self.fixSbx(pc, self.pc0()-pc)
	}

	return lastJmp
}
