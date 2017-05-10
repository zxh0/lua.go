package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"
import . "luago/lua/vm"

// kind of operands
const (
	ARG_CONST = iota
	ARG_REG
	ARG_UPVAL
	ARG_GLOBAL
)

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
	tmpAlloc := self.newTmpAllocator(a)

	b, kindB := self.toOpArg(exp.PrefixExp)
	if kindB != ARG_REG && kindB != ARG_UPVAL {
		b = tmpAlloc.allocTmp()
		self.exp(exp.PrefixExp, b, 1)
	}

	c, kindC := self.toOpArg(exp.KeyExp)
	if kindC != ARG_REG && kindC != ARG_CONST {
		c = tmpAlloc.allocTmp()
		self.exp(exp.KeyExp, c, 1)
	}

	tmpAlloc.freeAll()
	if kindB == ARG_UPVAL {
		self.getTabUp(exp.Line, a, b, c)
	} else {
		self.getTable(exp.Line, a, b, c)
	}
}

// r[a] := op exp
func (self *cg) unopExp(exp *UnopExp, a int) {
	tmpAlloc := self.newTmpAllocator(a)

	b := self.toREG(exp.Exp)
	if b < 0 {
		b = tmpAlloc.allocTmp()
		self.exp(exp.Exp, b, 1)
		tmpAlloc.freeTmp()
	}

	switch exp.Op {
	case TOKEN_OP_NOT:
		self.inst(exp.Line, OP_NOT, a, b, 0)
	case TOKEN_OP_BNOT:
		self.inst(exp.Line, OP_BNOT, a, b, 0)
	case TOKEN_OP_LEN:
		self.inst(exp.Line, OP_LEN, a, b, 0)
	case TOKEN_OP_UNM:
		self.inst(exp.Line, OP_UNM, a, b, 0)
	}
}

// r[a] := exp1 op exp2
func (self *cg) binopExp(exp *BinopExp, a int) {
	switch exp.Op {
	case TOKEN_OP_CONCAT:
		self.concatExp(exp, a)
	case TOKEN_OP_OR:
		self.logicalOrExp(exp, a)
	case TOKEN_OP_AND:
		self.logicalAndExp(exp, a)
	default:
		if _, ok := arithAndBitwiseBinops[exp.Op]; ok {
			self.arithAndBitwiseBinopExp(exp, a)
		} else {
			self.relationalBinopExp(exp, a)
		}
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

// r[a] := exp1 op exp2
func (self *cg) arithAndBitwiseBinopExp(exp *BinopExp, a int) {
	allocator := self.newTmpAllocator(a)

	rkb := self.toRK(exp.Exp1)
	if rkb < 0 {
		rkb = allocator.allocTmp()
		self.exp(exp.Exp1, rkb, 1)
	}

	rkc := self.toRK(exp.Exp2)
	if rkc < 0 {
		rkc = allocator.allocTmp()
		self.exp(exp.Exp2, rkc, 1)
	}

	allocator.freeAll()
	opcode := arithAndBitwiseBinops[exp.Op]
	self.inst(exp.Line, opcode, a, rkb, rkc)
}

// r[a] := exp1 op exp2
func (self *cg) relationalBinopExp(exp *BinopExp, a int) {
	allocator := self.newTmpAllocator(a)

	rkb := self.toRK(exp.Exp1)
	if rkb < 0 {
		rkb = allocator.allocTmp()
		self.exp(exp.Exp1, rkb, 1)
	}

	rkc := self.toRK(exp.Exp2)
	if rkc < 0 {
		rkc = allocator.allocTmp()
		self.exp(exp.Exp2, rkc, 1)
	}

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
	self.jmp(exp.Line, 1)
	self.loadBool(exp.Line, a, 0, 1)
	self.loadBool(exp.Line, a, 1, 0)
}

func (self *cg) logicalOrExp(exp *BinopExp, a int) {
	self.logicalExp(exp, a, 1)
}

func (self *cg) logicalAndExp(exp *BinopExp, a int) {
	self.logicalExp(exp, a, 0)
}

func (self *cg) logicalExp(exp *BinopExp, a, c int) {
	jmps := make([]int, 0, 4)

	for {
		if slot, ok := self.isLocVar(exp.Exp1); ok {
			if slot == a {
				self.test(exp.Line, a, c)
			} else {
				self.testSet(exp.Line, a, slot, c)
			}
			// } else if andExp, ok := castToLogicalAndExp(exp.Exp1); ok {
			// 	self.logicalAndExp(andExp, a)
			// } else if orExp, ok := castToLogicalOrExp(exp.Exp1); ok {
			// 	self.logicalOrExp(orExp, a)
		} else {
			tmp := self.allocTmp()
			self.exp(exp.Exp1, tmp, 1)
			self.freeTmp()
			self.testSet(exp.Line, a, tmp, c)
		}
		jmps = append(jmps, self.jmp(exp.Line, 0))

		if exp2, ok := exp.Exp2.(*BinopExp); ok && exp2.Op == exp.Op {
			exp = exp2
		} else {
			self.exp(exp.Exp2, a, 1)
			break
		}
	}

	jmpTo := self.pc() - 1 // todo
	for _, jmpFrom := range jmps {
		self.fixSbx(jmpFrom, jmpTo-jmpFrom-1)
	}
}

func (self *cg) isLocVar(exp Exp) (int, bool) {
	if nameExp, ok := exp.(*NameExp); ok {
		if slot := self.slotOf(nameExp.Name); slot >= 0 {
			return slot, true
		}
	}
	return -1, false
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
	case *NameExp:
		if slot := self.slotOf(x.Name); slot >= 0 {
			return slot, ARG_REG
		} else if idx := self.lookupUpval(x.Name); idx >= 0 {
			return idx, ARG_UPVAL
		} else {
			self.indexOf(x.Name)
			return -1, ARG_GLOBAL
		}
	}
	return -1, -1
}

func (self *cg) toREG(exp Exp) int {
	if nameExp, ok := exp.(*NameExp); ok {
		return self.slotOf(nameExp.Name)
	}
	return -1
}

func (self *cg) toRK(exp Exp) int {
	switch x := exp.(type) {
	case *NilExp:
		return self.indexOf(nil)
	case *FalseExp:
		return self.indexOf(false)
	case *TrueExp:
		return self.indexOf(true)
	case *IntegerExp:
		return self.indexOf(x.Val)
	case *FloatExp:
		return self.indexOf(x.Val)
	case *StringExp:
		return self.indexOf(x.Str)
	case *NameExp:
		return self.slotOf(x.Name)
	default:
		return -1
	}
}
