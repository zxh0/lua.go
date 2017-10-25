package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"
import . "luago/vm"

// kind of operands
const (
	ARG_CONST  = 1 // const index
	ARG_REG    = 2 // register index
	ARG_UPVAL  = 4 // upvalue index
	ARG_GLOBAL = 8 // ?
	ARG_RK     = ARG_REG | ARG_CONST
	ARG_RU     = ARG_REG | ARG_UPVAL
	ARG_RUK    = ARG_REG | ARG_UPVAL | ARG_CONST
	//ARG_RUG    = ARG_REG | ARG_UPVAL | ARG_GLOBAL
	// ARG_LOCAL ?
	// ARG_TMP ?
)

// todo: rename to evalExp()?
func (self *codeGen) cgExp(node Exp, a, n int) {
	switch exp := node.(type) {
	case *NilExp:
		self.emitLoadNil(exp.Line, a, n)
	case *FalseExp:
		self.emitLoadBool(exp.Line, a, 0, 0)
	case *TrueExp:
		self.emitLoadBool(exp.Line, a, 1, 0)
	case *IntegerExp:
		self.emitLoadK(exp.Line, a, exp.Val)
	case *FloatExp:
		self.emitLoadK(exp.Line, a, exp.Val)
	case *StringExp:
		self.emitLoadK(exp.Line, a, exp.Str)
	case *VarargExp:
		self.emitVararg(exp.Line, a, n)
	case *ParensExp:
		self.cgExp(exp.Exp, a, 1)
	case *NameExp:
		self.cgNameExp(exp, a)
	case *TableConstructorExp:
		self.cgTableConstructorExp(exp, a)
	case *FuncDefExp:
		self.cgFuncDefExp(exp, a)
	case *FuncCallExp:
		self.cgFuncCallExp(exp, a, n)
	case *BracketsExp:
		self.cgBracketsExp(exp, a)
	case *ConcatExp:
		self.cgConcatExp(exp, a)
	case *UnopExp:
		self.cgUnopExp(exp, a)
	case *BinopExp:
		self.cgBinopExp(exp, a)
	}
}

func (self *codeGen) cgTableConstructorExp(exp *TableConstructorExp, a int) {
	nArr := exp.NArr
	nExps := len(exp.KeyExps)
	lastExpIsVarargOrFuncCall := nExps > 0 &&
		isVarargOrFuncCallExp(exp.ValExps[nExps-1])

	self.emitNewTable(exp.Line, a, nArr, nExps - nArr)

	for i, keyExp := range exp.KeyExps {
		valExp := exp.ValExps[i]

		if nArr > 0 {
			if idx, ok := keyExp.(int); ok {
				_a := self.allocReg()
				if i == nExps-1 && lastExpIsVarargOrFuncCall {
					self.cgExp(valExp, _a, -1)
				} else {
					self.cgExp(valExp, _a, 1)
				}

				if idx%50 == 0 || idx == nArr { // LFIELDS_PER_FLUSH
					if idx%50 == 0 {
						self.freeRegs(50)
					} else {
						self.freeRegs(idx%50)
					}
					line := lastLineOfExp(valExp)
					if i == nExps-1 && lastExpIsVarargOrFuncCall {
						self.emitSetList(line, a, 0, idx/50 + 1)
					} else {
						self.emitSetList(line, a, idx%50, idx/50 + 1)
					}
				}

				continue
			}
		}

		b := self.allocReg()
		self.cgExp(keyExp, b, 1)
		c := self.allocReg()
		self.cgExp(valExp, c, 1)
		self.freeRegs(2)

		line := lastLineOfExp(valExp)
		self.emitSetTable(line, a, b, c)
	}
}

// f[a] := function(args) body end
func (self *codeGen) cgFuncDefExp(exp *FuncDefExp, a int) {
	bx := self.genSubProto(exp)
	self.emitClosure(exp.LastLine, a, bx)
}

// r[a] := f(args)
func (self *codeGen) cgFuncCallExp(exp *FuncCallExp, a, n int) {
	nArgs := self.prepFuncCall(exp, a)
	self.emitCall(exp.Line, a, nArgs, n)
}

// return f(args)
func (self *codeGen) cgTailCallExp(exp *FuncCallExp, a int) {
	nArgs := self.prepFuncCall(exp, a)
	self.emitTailCall(exp.Line, a, nArgs)
}

func (self *codeGen) prepFuncCall(exp *FuncCallExp, a int) int {
	nArgs := len(exp.Args)
	lastArgIsVarargOrFuncCall := false

	self.cgExp(exp.PrefixExp, a, 1)
	if exp.MethodName != "" {
		self.allocReg()
		idx := self.indexOfConstant(exp.MethodName)
		self.emitSelf(exp.Line, a, a, idx)
	}
	for i, arg := range exp.Args {
		tmp := self.allocReg()
		if i == nArgs-1 && isVarargOrFuncCallExp(arg) {
			lastArgIsVarargOrFuncCall = true
			self.cgExp(arg, tmp, -1)
		} else {
			self.cgExp(arg, tmp, 1)
		}
	}
	self.freeRegs(nArgs)

	if lastArgIsVarargOrFuncCall {
		nArgs = -1
	}
	if exp.MethodName != "" {
		self.freeReg()
		nArgs++
	}

	return nArgs
}

// r[a] := name
func (self *codeGen) cgNameExp(exp *NameExp, a int) {
	if slot := self.slotOf(exp.Name); slot >= 0 {
		self.emitMove(exp.Line, a, slot)
	} else if idx := self.lookupUpval(exp.Name); idx >= 0 {
		self.emitGetUpval(exp.Line, a, idx)
	} else { // x => _ENV['x']
		bracketsExp := &BracketsExp{
			Line:      exp.Line,
			PrefixExp: &NameExp{exp.Line, "_ENV"},
			KeyExp:    &StringExp{exp.Line, exp.Name},
		}
		self.cgBracketsExp(bracketsExp, a)
	}
}

// r[a] := prefix[key]
func (self *codeGen) cgBracketsExp(exp *BracketsExp, a int) {
	oldRegs := self.usedRegs()
	b, kindB := self.expToOpArg(exp.PrefixExp, ARG_RU)
	c, _ := self.expToOpArg(exp.KeyExp, ARG_RK)
	self.freeRegs(self.usedRegs() - oldRegs)

	if kindB == ARG_UPVAL {
		self.emitGetTabUp(exp.Line, a, b, c)
	} else {
		self.emitGetTable(exp.Line, a, b, c)
	}
}

// r[a] := op exp
func (self *codeGen) cgUnopExp(exp *UnopExp, a int) {
	oldRegs := self.usedRegs()
	b, _ := self.expToOpArg(exp.Exp, ARG_REG)
	self.emitUnaryOp(exp.Line, exp.Op, a, b)
	self.freeRegs(self.usedRegs() - oldRegs)
}

// r[a] := exp1 op exp2
func (self *codeGen) cgBinopExp(exp *BinopExp, a int) {
	switch exp.Op {
	case TOKEN_OP_AND, TOKEN_OP_OR:
		oldRegs := self.usedRegs()

		b, _ := self.expToOpArg(exp.Exp1, ARG_REG)
		self.freeRegs(self.usedRegs() - oldRegs)
		if exp.Op == TOKEN_OP_AND {
			self.emitTestSet(exp.Line, a, b, 0)
		} else {
			self.emitTestSet(exp.Line, a, b, 1)
		}
		pcOfJmp := self.emitJmp(exp.Line, 0)

		b, _ = self.expToOpArg(exp.Exp2, ARG_REG)
		self.freeRegs(self.usedRegs() - oldRegs)
		self.emitMove(exp.Line, a, b)		
		self.fixSbx(pcOfJmp, self.pc()-pcOfJmp)
	default:
		oldRegs := self.usedRegs()
		b, _ := self.expToOpArg(exp.Exp1, ARG_RK)
		c, _ := self.expToOpArg(exp.Exp2, ARG_RK)
		self.emitBinaryOp(exp.Line, exp.Op, a, b, c)
		self.freeRegs(self.usedRegs() - oldRegs)
	}
}

// r[a] := exp1 .. exp2
func (self *codeGen) cgConcatExp(exp *ConcatExp, a int) {
	for _, subExp := range exp.Exps {
		a := self.allocReg()
		self.cgExp(subExp, a, 1)
	}

	c := self.usedRegs() - 1
	b := c - len(exp.Exps) + 1
	self.freeRegs(c - b + 1)
	self.emit(exp.Line, OP_CONCAT, a, b, c)
}

func (self *codeGen) expToOpArg(exp Exp, argKinds int) (arg, argKind int) {
	if argKinds&ARG_CONST > 0 {
		switch x := exp.(type) {
		case *NilExp:
			return self.indexOfConstant(nil), ARG_CONST
		case *FalseExp:
			return self.indexOfConstant(false), ARG_CONST
		case *TrueExp:
			return self.indexOfConstant(true), ARG_CONST
		case *IntegerExp:
			return self.indexOfConstant(x.Val), ARG_CONST
		case *FloatExp:
			return self.indexOfConstant(x.Val), ARG_CONST
		case *StringExp:
			return self.indexOfConstant(x.Str), ARG_CONST
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
	a := self.allocReg()
	self.cgExp(exp, a, 1)
	return a, ARG_REG
}
