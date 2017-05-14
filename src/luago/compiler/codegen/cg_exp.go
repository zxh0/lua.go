package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"
import . "luago/lua/vm"

// kind of operands
const (
	ARG_CONST  = 1 // const index
	ARG_REG    = 2 // register index
	ARG_UPVAL  = 4 // upvalue index
	ARG_GLOBAL = 8 // ?
	ARG_RK     = ARG_REG | ARG_CONST
	ARG_RU     = ARG_REG | ARG_UPVAL
	// ARG_LOCAL ?
	// ARG_TMP ?
)

// todo: rename to evalExp()?
func (self *cg) exp(node Exp, a, n int) {
	switch exp := node.(type) {
	case *NilExp:
		self.loadNil(exp.Line, a, n)
	case *FalseExp:
		self.loadBool(exp.Line, a, 0, 0)
	case *TrueExp:
		self.loadBool(exp.Line, a, 1, 0)
	case *IntegerExp:
		self.loadK(exp.Line, a, exp.Val)
	case *FloatExp:
		self.loadK(exp.Line, a, exp.Val)
	case *StringExp:
		self.loadK(exp.Line, a, exp.Str)
	case *VarargExp:
		self.vararg(exp.Line, a, n)
	case *ParensExp:
		self.exp(exp.Exp, a, 1)
	case *NameExp:
		self.nameExp(exp, a)
	case *TableConstructorExp:
		self.tcExp(exp, a)
	case *FuncDefExp:
		self.funcDefExp(exp, a)
	case *FuncCallExp:
		self.funcCallExp(exp, a, n)
	case *BracketsExp:
		self.bracketsExp(exp, a)
	case *UnopExp:
		self.unopExp(exp, a)
	case *BinopExp:
		self.binopExp(exp, a)
	}
}

func (self *cg) tcExp(exp *TableConstructorExp, a int) {
	nExps := len(exp.KeyExps)
	nArr := exp.NArr
	nRec := nExps - nArr
	lastExpIsVarargOrFuncCall := nExps > 0 &&
		isVarargOrFuncCallExp(exp.ValExps[nExps-1])

	if lastExpIsVarargOrFuncCall {
		self.newTable(exp.Line, a, nArr-1, nRec)
	} else {
		self.newTable(exp.Line, a, nArr, nRec)
	}

	for i, keyExp := range exp.KeyExps {
		valExp := exp.ValExps[i]

		if nArr > 0 {
			if idx, ok := keyExp.(int); ok {
				tmp := self.allocTmp()
				if i == nExps-1 && lastExpIsVarargOrFuncCall {
					self.exp(valExp, tmp, -1)
				} else {
					self.exp(valExp, tmp, 1)
				}

				if idx%50 == 0 {
					self.freeTmps(50)
					line := lineOfExp(valExp)
					if i == nExps-1 && lastExpIsVarargOrFuncCall {
						self.setList(line, a, 0, idx/50) // todo
					} else {
						self.setList(line, a, 50, idx/50) // todo
					}
				}

				continue
			}
		}

		nTmps := 0
		iKey, tKey := self.toOpArg(keyExp)
		if tKey != ARG_CONST && tKey != ARG_REG {
			iKey = self.allocTmp()
			nTmps++
			self.exp(keyExp, iKey, 1)
		}

		iVal, tVal := self.toOpArg(valExp)
		if tVal != ARG_CONST && tVal != ARG_REG {
			iVal = self.allocTmp()
			nTmps++
			self.exp(valExp, iVal, 1)
		}
		self.freeTmps(nTmps)
		self.setTable(lineOfExp(valExp), a, iKey, iVal)
	}

	if nArr > 0 {
		self.freeTmps(nArr)
		if lastExpIsVarargOrFuncCall {
			self.setList(exp.LastLine, a, 0, 1) // todo
		} else {
			self.setList(exp.LastLine, a, nArr%50, nArr/50+1) // todo
		}
	}
}

// f[a] := function(args) body end
func (self *cg) funcDefExp(exp *FuncDefExp, a int) {
	bx := self.genSubProto(exp)
	self.closure(exp.LastLine, a, bx)
}

// r[a] := f(args)
func (self *cg) funcCallExp(exp *FuncCallExp, a, n int) {
	nArgs := self.prepFuncCall(exp, a)
	self.call(exp.Line, a, nArgs, n)
}

// return f(args)
func (self *cg) tailCallExp(exp *FuncCallExp, a int) {
	nArgs := self.prepFuncCall(exp, a)
	self.tailCall(exp.Line, a, nArgs)
}

func (self *cg) prepFuncCall(exp *FuncCallExp, a int) int {
	nArgs := len(exp.Args)
	lastArgIsVarargOrFuncCall := false

	self.exp(exp.PrefixExp, a, 1)
	if exp.MethodName != "" {
		self.allocTmp()
		idx := self.indexOf(exp.MethodName)
		self._self(exp.Line, a, a, idx)
	}
	for i, arg := range exp.Args {
		tmp := self.allocTmp()
		if i == nArgs-1 && isVarargOrFuncCallExp(arg) {
			lastArgIsVarargOrFuncCall = true
			self.exp(arg, tmp, -1)
		} else {
			self.exp(arg, tmp, 1)
		}
	}
	self.freeTmps(nArgs)

	if lastArgIsVarargOrFuncCall {
		nArgs = -1
	}
	if exp.MethodName != "" {
		self.freeTmp()
		nArgs++
	}

	return nArgs
}

// r[a] := name
func (self *cg) nameExp(exp *NameExp, a int) {
	if slot := self.slotOf(exp.Name); slot >= 0 {
		self.move(exp.Line, a, slot)
	} else if idx := self.lookupUpval(exp.Name); idx >= 0 {
		self.getUpval(exp.Line, a, idx)
	} else { // x => _ENV['x']
		bracketsExp := &BracketsExp{
			Line:      exp.Line,
			PrefixExp: &NameExp{exp.Line, "_ENV"},
			KeyExp:    &StringExp{exp.Line, exp.Name},
		}
		self.bracketsExp(bracketsExp, a)
	}
}

// r[a] := prefix[key]
func (self *cg) bracketsExp(exp *BracketsExp, a int) {
	allocator := self.newTmpAllocator(a)
	b, kindB := self.exp2OpArg(exp.PrefixExp, ARG_RU, allocator)
	c, _ := self.exp2OpArg(exp.KeyExp, ARG_RK, allocator)
	allocator.freeAll()

	if kindB == ARG_UPVAL {
		self.getTabUp(exp.Line, a, b, c)
	} else {
		self.getTable(exp.Line, a, b, c)
	}
}

// r[a] := op exp
func (self *cg) unopExp(exp *UnopExp, a int) {
	allocator := self.newTmpAllocator(a)
	b, _ := self.exp2OpArg(exp.Exp, ARG_REG, allocator)
	self.unaryOp(exp.Line, exp.Op, a, b)
	allocator.freeAll()
}

// r[a] := exp1 op exp2
func (self *cg) binopExp(exp *BinopExp, a int) {
	switch exp.Op {
	case TOKEN_OP_CONCAT:
		self.concatExp(exp, a)
	case TOKEN_OP_OR, TOKEN_OP_AND:
		self.logicalBinopExp(exp, a)
	default:
		allocator := self.newTmpAllocator(a)
		rkb, _ := self.exp2OpArg(exp.Exp1, ARG_RK, allocator)
		rkc, _ := self.exp2OpArg(exp.Exp2, ARG_RK, allocator)
		self.binaryOp(exp.Line, exp.Op, a, rkb, rkc)
		allocator.freeAll()
	}
}

// r[a] := exp1 .. exp2
func (self *cg) concatExp(exp *BinopExp, a int) {
	allocator := self.newTmpAllocator(a)
	line, b, c := exp.Line, -1, -1

	for {
		tmp := allocator.allocTmp()
		self.exp(exp.Exp1, tmp, 1)
		if b < 0 {
			b = tmp
			c = b
		} else {
			c++
		}

		if exp2, ok := castToConcatExp(exp.Exp2); ok {
			exp = exp2
		} else {
			tmp := allocator.allocTmp()
			self.exp(exp.Exp2, tmp, 1)
			c++
			break
		}
	}

	allocator.freeAll()
	self.inst(line, OP_CONCAT, a, b, c)
}

func (self *cg) logicalBinopExp(exp *BinopExp, a int) {
	list := logicalBinopExpToList(exp)
	for node := list; node != nil; node = node.next {
		node.startPc = self.pc()
		allocator := self.newTmpAllocator(a)
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
		} else if b != a {
			self.move(lineOfExp(node.exp), a, b)
		}
	}
	for node := list; node != nil; node = node.next {
		if node.next != nil {
			sbx := 0
			if node.jmpTo != nil {
				sbx = node.jmpTo.startPc - node.jmpPc
			} else {
				sbx = self.pc() - node.jmpPc
			}
			self.fixSbx(node.jmpPc, sbx)
		}
	}
}

// todo: rename
func (self *cg) toOperand0(exp Exp) (int, int) {
	switch x := exp.(type) {
	case *NilExp, *FalseExp, *TrueExp,
		*IntegerExp, *FloatExp, *StringExp,
		*VarargExp:
		return -1, ARG_CONST
	case *NameExp:
		if slot := self.slotOf(x.Name); slot >= 0 {
			return slot, ARG_REG
		} else if idx := self.lookupUpval(x.Name); idx >= 0 {
			return idx, ARG_UPVAL
		} else {
			return self.indexOf(x.Name), ARG_GLOBAL
		}
	}
	return -1, -1
}

func (self *cg) toOpArg(exp Exp) (int, int) {
	return self._toOpArg(exp, ARG_CONST|ARG_REG|ARG_UPVAL)
}

func (self *cg) exp2OpArg(exp Exp, argKinds int,
	allocator *tmpAllocator) (arg, argKind int) {

	arg, argKind = self._toOpArg(exp, argKinds)
	if arg < 0 {
		argKind = ARG_REG
		arg = allocator.allocTmp()
		self.exp(exp, arg, 1)
	}
	return
}

// todo: rename
func (self *cg) _toOpArg(exp Exp, argKinds int) (arg, argKind int) {
	if argKinds&ARG_CONST > 0 {
		switch x := exp.(type) {
		case *NilExp:
			return self.indexOf(nil), ARG_CONST
		case *FalseExp:
			return self.indexOf(false), ARG_CONST
		case *TrueExp:
			return self.indexOf(true), ARG_CONST
		case *IntegerExp:
			return self.indexOf(x.Val), ARG_CONST
		case *FloatExp:
			return self.indexOf(x.Val), ARG_CONST
		case *StringExp:
			return self.indexOf(x.Str), ARG_CONST
		}
	}
	if argKinds&ARG_REG > 0 {
		if nameExp, ok := exp.(*NameExp); ok {
			if slot := self.slotOf(nameExp.Name); slot >= 0 {
				return slot, ARG_REG
			}
		}
	}
	if argKinds&ARG_UPVAL > 0 {
		if nameExp, ok := exp.(*NameExp); ok {
			if idx := self.lookupUpval(nameExp.Name); idx >= 0 {
				return idx, ARG_UPVAL
			}
		}
	}
	return -1, 0 // todo
}
