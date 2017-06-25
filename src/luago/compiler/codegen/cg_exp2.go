package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"
import . "luago/vm"

func (self *codeGen) logicalBinopExp(exp *BinopExp, a int) {
	list := logicalBinopExpToList(exp)
	hasRelationalBinop := false
	lastExpIsRelationalBinop := false
	lineOfLastExp := 0

	for node := list; node != nil; node = node.next {
		node.startPc = self.pc()

		// allocator := self.newTmpAllocator(a)
		allocator := &tmpAllocator{self.scope, a, 0}
		if self.isLocVarSlot(a) && node.next != nil {
			allocator = self.newTmpAllocator(-1)
		}

		if bexp, ok := castToRelationalBinopExp(node.exp); ok {
			hasRelationalBinop = true
			if node.next != nil {
				self.testRelationalBinopExpX(bexp, allocator, 0)
				node.jmpPc = self.emitJmp(node.line, 0)
			} else {
				lastExpIsRelationalBinop = true
				lineOfLastExp = lineOfExp(node.exp)
				self.testRelationalBinopExpX(bexp, allocator, 1)
			}
			continue
		}

		b, _ := self.exp2OpArg(node.exp, ARG_REG, allocator)
		allocator.freeAll()
		if node.next != nil {
			c := 1
			if node.op == TOKEN_OP_AND {
				c = 0
			}
			if b == a {
				self.emitTest(node.line, a, c)
			} else if node.jmpTo != nil {
				self.emitTest(node.line, b, c)
			} else {
				self.emitTestSet(node.line, a, b, c)
			}
			node.jmpPc = self.emitJmp(node.line, 0)
		} else {
			lineOfLastExp = lineOfExp(node.exp)
			if b != a {
				self.emitMove(lineOfLastExp, a, b)
			}
		}
	}
	if hasRelationalBinop {
		if lastExpIsRelationalBinop {
			self.emitJmp(lineOfLastExp, 1)
		} else {
			self.emitJmp(lineOfLastExp, 2)
		}
		self.emitLoadBool(lineOfLastExp, a, 0, 1)
		self.emitLoadBool(lineOfLastExp, a, 1, 0)
	}
	for node := list; node != nil; node = node.next {
		if node.next != nil {
			sbx := 0
			if node.jmpTo != nil {
				sbx = node.jmpTo.startPc - node.jmpPc
			} else {
				sbx = self.pc() - node.jmpPc
			}
			if hasRelationalBinop {
				sbx -= 2
			}
			self.fixSbx(node.jmpPc, sbx)
		}
	}
}

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
