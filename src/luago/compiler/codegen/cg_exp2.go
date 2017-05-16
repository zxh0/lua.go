package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"
import . "luago/lua/vm"

func (self *cg) testExp(node Exp, lineOfLastJmp int) (pendingJmps []int) {
	if bexp, ok := node.(*BinopExp); ok {
		switch bexp.Op {
		case TOKEN_OP_EQ, TOKEN_OP_NE,
			TOKEN_OP_LT, TOKEN_OP_GT,
			TOKEN_OP_LE, TOKEN_OP_GE:
			self.testRelationalBinopExp(bexp)
			pc := self.jmp(lineOfLastJmp, 0)
			return []int{pc}
		case TOKEN_OP_AND, TOKEN_OP_OR:
			pendingJmps := self.testLogicalBinopExp(bexp, lineOfLastJmp)
			return pendingJmps
		}
	}

	allocator := self.newTmpAllocator(-1)
	a, _ := self.exp2OpArg(node, ARG_REG, allocator)
	allocator.freeAll()

	self.test(lineOfLastJmp, a, 0)
	pc := self.jmp(lineOfLastJmp, 0)
	return []int{pc}
}

func (self *cg) testLogicalBinopExp(exp *BinopExp, lineOfLastJmp int) (pendingJmps []int) {
	list := logicalBinopExpToList(exp)
	for node := list; node != nil; node = node.next {
		node.startPc = self.pc()
		allocator := self.newTmpAllocator(-1)

		if isRelationalBinopExp(node.exp) {
			if node.next != nil {
				if node.op == TOKEN_OP_AND {
					self.testRelationalBinopExp(node.exp.(*BinopExp))
					node.jmpPc = self.jmp(lineOfLastJmp, 0)
				} else {
					self.testRelationalBinopExp222(node.exp.(*BinopExp))
					node.jmpPc = self.jmp(lineOfLastJmp, 0)
				}
			} else {
				self.testRelationalBinopExp(node.exp.(*BinopExp))
				pc := self.jmp(lineOfLastJmp, 0)
				pendingJmps = append(pendingJmps, pc)
			}
		} else {
			b, _ := self.exp2OpArg(node.exp, ARG_REG, allocator)
			allocator.freeAll()
			if node.next != nil {
				c := 1
				if node.op == TOKEN_OP_AND {
					c = 0
				}
				self.test(node.line, b, c)
				node.jmpPc = self.jmp(node.line, 0)
			} else {
				self.test(lineOfLastJmp, b, 0)
				pc := self.jmp(lineOfLastJmp, 0)
				pendingJmps = append(pendingJmps, pc)
			}
		}
	}
	for node := list; node != nil; node = node.next {
		if node.next != nil {
			if node.jmpTo != nil {
				sbx := node.jmpTo.startPc - node.jmpPc
				self.fixSbx(node.jmpPc, sbx)
			} else {
				if node.op == TOKEN_OP_OR {
					sbx := self.pc() - node.jmpPc
					self.fixSbx(node.jmpPc, sbx)
				} else {
					pendingJmps = append(pendingJmps, node.jmpPc)
				}
			}
		}
	}
	return pendingJmps
}

func (self *cg) testRelationalBinopExp(exp *BinopExp) {
	allocator := self.newTmpAllocator(-1)
	rkb, _ := self.exp2OpArg(exp.Exp1, ARG_RK, allocator)
	rkc, _ := self.exp2OpArg(exp.Exp2, ARG_RK, allocator)
	allocator.freeAll()

	switch exp.Op {
	case TOKEN_OP_EQ:
		self.inst(exp.Line, OP_EQ, 0, rkb, rkc)
	case TOKEN_OP_NE:
		self.inst(exp.Line, OP_EQ, 1, rkb, rkc)
	case TOKEN_OP_LT:
		self.inst(exp.Line, OP_LT, 0, rkb, rkc)
	case TOKEN_OP_GT:
		self.inst(exp.Line, OP_LT, 0, rkc, rkb)
	case TOKEN_OP_LE:
		self.inst(exp.Line, OP_LE, 0, rkb, rkc)
	case TOKEN_OP_GE:
		self.inst(exp.Line, OP_LE, 0, rkc, rkb)
	}
}

// todo
func (self *cg) testRelationalBinopExp222(exp *BinopExp) {
	allocator := self.newTmpAllocator(-1)
	rkb, _ := self.exp2OpArg(exp.Exp1, ARG_RK, allocator)
	rkc, _ := self.exp2OpArg(exp.Exp2, ARG_RK, allocator)
	allocator.freeAll()

	switch exp.Op {
	case TOKEN_OP_EQ:
		self.inst(exp.Line, OP_EQ, 1, rkb, rkc)
	case TOKEN_OP_NE:
		self.inst(exp.Line, OP_EQ, 0, rkb, rkc)
	case TOKEN_OP_LT:
		self.inst(exp.Line, OP_LT, 1, rkb, rkc)
	case TOKEN_OP_GT:
		self.inst(exp.Line, OP_LT, 1, rkc, rkb)
	case TOKEN_OP_LE:
		self.inst(exp.Line, OP_LE, 1, rkb, rkc)
	case TOKEN_OP_GE:
		self.inst(exp.Line, OP_LE, 1, rkc, rkb)
	}
}

func isRelationalBinopExp(exp Exp) bool {
	if binopExp, ok := exp.(*BinopExp); ok {
		switch binopExp.Op {
		case TOKEN_OP_EQ, TOKEN_OP_NE,
			TOKEN_OP_LT, TOKEN_OP_LE,
			TOKEN_OP_GT, TOKEN_OP_GE:
			return true
		}
	}
	return false
}
