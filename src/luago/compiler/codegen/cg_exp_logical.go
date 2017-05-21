package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

func (self *cg) logicalBinopExp(exp *BinopExp, a int) {
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

		if isRelationalBinopExp(node.exp) {
			hasRelationalBinop = true
			if node.next != nil {
				self.testRelationalBinopExpX(node.exp.(*BinopExp), allocator, 0)
				node.jmpPc = self.jmp(node.line, 0)
			} else {
				lastExpIsRelationalBinop = true
				lineOfLastExp = lineOfExp(node.exp)
				self.testRelationalBinopExpX(node.exp.(*BinopExp), allocator, 1)
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
				self.test(node.line, a, c)
			} else if node.jmpTo != nil {
				self.test(node.line, b, c)
			} else {
				self.testSet(node.line, a, b, c)
			}
			node.jmpPc = self.jmp(node.line, 0)
		} else {
			lineOfLastExp = lineOfExp(node.exp)
			if b != a {
				self.move(lineOfLastExp, a, b)
			}
		}
	}
	if hasRelationalBinop {
		if lastExpIsRelationalBinop {
			self.jmp(lineOfLastExp, 1)
		} else {
			self.jmp(lineOfLastExp, 2)
		}
		self.loadBool(lineOfLastExp, a, 0, 1)
		self.loadBool(lineOfLastExp, a, 1, 0)
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

func (self *cg) testLogicalBinopExp(exp *BinopExp, lineOfLastJmp int) (pendingJmps []int) {
	list := logicalBinopExpToList(exp)
	for node := list; node != nil; node = node.next {
		node.startPc = self.pc()
		allocator := self.newTmpAllocator(-1)

		if isRelationalBinopExp(node.exp) {
			if node.next != nil {
				if node.op == TOKEN_OP_AND {
					self.testRelationalBinopExp(node.exp.(*BinopExp), 0)
					node.jmpPc = self.jmp(lineOfLastJmp, 0)
				} else {
					self.testRelationalBinopExp(node.exp.(*BinopExp), 1)
					node.jmpPc = self.jmp(lineOfLastJmp, 0)
				}
			} else {
				self.testRelationalBinopExp(node.exp.(*BinopExp), 0)
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
