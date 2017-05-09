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

func (self *cg) testExp(node Exp, a int) {
	switch exp := node.(type) {
	case *NilExp:
		self.loadNil(exp.Line, a, 1)
	case *FalseExp:
		self.loadBool(exp.Line, a, 0)
	case *TrueExp:
		self.loadBool(exp.Line, a, 1)
	case *IntegerExp:
		self.loadK(exp.Line, a, exp.Val)
	case *FloatExp:
		self.loadK(exp.Line, a, exp.Val)
	case *StringExp:
		self.loadK(exp.Line, a, exp.Str)
	case *TableConstructorExp:
		self.tcExp(exp, a)
	case *FuncDefExp:
		self.funcDefExp(exp, a)
	case *VarargExp:
		self.vararg(exp.Line, a, 1)
	case *FuncCallExp:
		self.funcCallExp(exp, a, 1)
	case *NameExp:
		self.nameExp(exp, a)
	case *ParensExp:
		self.exp(exp.Exp, a, 1)
	case *BracketsExp:
		self.bracketsExp(exp, a)
	case *BinopExp:
		self.binopExp(exp, a, 0)
	case *UnopExp:
		self.unopExp(exp, a)
	default:
		panic("todo!")
	}
}

// todo: rename to evalExp()?
func (self *cg) exp(node Exp, a, n int) {
	switch exp := node.(type) {
	case *NilExp:
		self.loadNil(exp.Line, a, n)
	case *FalseExp:
		self.loadBool(exp.Line, a, 0)
	case *TrueExp:
		self.loadBool(exp.Line, a, 1)
	case *IntegerExp:
		self.loadK(exp.Line, a, exp.Val)
	case *FloatExp:
		self.loadK(exp.Line, a, exp.Val)
	case *StringExp:
		self.loadK(exp.Line, a, exp.Str)
	case *TableConstructorExp:
		self.tcExp(exp, a)
	case *FuncDefExp:
		self.funcDefExp(exp, a)
	case *VarargExp:
		self.vararg(exp.Line, a, n)
	case *FuncCallExp:
		self.funcCallExp(exp, a, n)
	case *NameExp:
		self.nameExp(exp, a)
	case *ParensExp:
		self.exp(exp.Exp, a, 1)
	case *BracketsExp:
		self.bracketsExp(exp, a)
	case *BinopExp:
		self.binopExp(exp, a, n)
	case *UnopExp:
		self.unopExp(exp, a)
	default:
		panic("todo!")
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

func (self *cg) funcDefExp(exp *FuncDefExp, a int) {
	bx := self.genSubProto(exp)
	self.closure(exp.LastLine, a, bx)
}

func (self *cg) funcCallExp(exp *FuncCallExp, a, n int) {
	self._funcCallExp(exp, a, n, false)
}

func (self *cg) tailCallExp(exp *FuncCallExp, a int) {
	self._funcCallExp(exp, a, 0, true)
}

// todo: rename?
func (self *cg) _funcCallExp(exp *FuncCallExp, a, n int, tailCall bool) {
	nArgs := len(exp.Args)
	lastArgIsVarargOrFuncCall := false

	self.loadFunc(exp.PrefixExp, a)
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
	if tailCall {
		self.tailCall(exp.Line, a, nArgs)
	} else {
		self.call(exp.Line, a, nArgs, n)
	}
}

func (self *cg) loadFunc(pexp Exp, tmp int) {
	if nameExp, ok := pexp.(*NameExp); ok {
		if slot := self.slotOf(nameExp.Name); slot >= 0 {
			self.move(nameExp.Line, tmp, slot)
			return
		}
	}
	self.exp(pexp, tmp, 1)
}

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

// prefix[key]
func (self *cg) bracketsExp(exp *BracketsExp, a int) {
	iTab, tTab := self.toOpArg(exp.PrefixExp)
	iKey, tKey := self.toOpArg(exp.KeyExp)

	if tTab != ARG_REG && tTab != ARG_UPVAL {
		iTab = a
		self.exp(exp.PrefixExp, iTab, 1)
	}
	if tKey != ARG_REG && tKey != ARG_CONST {
		if tTab != ARG_REG && tTab != ARG_UPVAL {
			iKey = self.allocTmp()
			self.exp(exp.KeyExp, iKey, 1)
			self.freeTmp()
		} else {
			iKey = a
			self.exp(exp.KeyExp, iKey, 1)
		}
	}

	if tTab == ARG_UPVAL {
		self.getTabUp(exp.Line, a, iTab, iKey)
	} else {
		self.getTable(exp.Line, a, iTab, iKey)
	}
}

func (self *cg) unopExp(exp *UnopExp, a int) {
	b := self.toREG(exp.Exp)
	if b < 0 {
		b = self.allocTmp()
		self.exp(exp.Exp, b, 1)
		self.freeTmp()
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

func (self *cg) binopExp(exp *BinopExp, a, n int) {
	switch exp.Op {
	case TOKEN_OP_POW:
		self.powExp(exp, a)
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
			self.relationalBinopExp(exp, a, n)
		}
	}
}

func (self *cg) powExp(exp *BinopExp, a int) {
	aIsFree := self.isTmpVar(a)
	operands := make([]int, 0, 4) // todo: rename

	nTmps := 0
	for {
		if rk := self.toRK(exp.Exp1); rk >= 0 {
			operands = append(operands, rk, 0, exp.Line)
		} else {
			var tmp int
			if aIsFree {
				tmp = a
				aIsFree = false
			} else {
				tmp = self.allocTmp()
				nTmps++
			}
			self.exp(exp.Exp1, tmp, 1)
			operands = append(operands, tmp, 1, exp.Line)
		}

		if exp.Exp2 == nil {
			break
		} else if exp2, ok := castToPowExp(exp.Exp2); ok {
			exp = exp2
		} else {
			exp = &BinopExp{Exp1: exp.Exp2, Exp2: nil}
		}
	}

	for i := len(operands) - 2; i > 3; i -= 3 {
		yIsTmp := operands[i] == 1
		iy := operands[i-1]
		line := operands[i-2]
		xIsTmp := operands[i-3] == 1
		ix := operands[i-4]
		if xIsTmp {
			self.inst(line, OP_POW, ix, ix, iy)
		} else if yIsTmp {
			self.inst(line, OP_POW, iy, ix, iy)
			operands[i-3] = 1
			operands[i-4] = iy
		} else {
			var tmp int
			if aIsFree {
				tmp = a
				aIsFree = false
			} else {
				tmp = self.allocTmp()
				nTmps++
			}
			self.inst(line, OP_POW, tmp, ix, iy)
			operands[i-3] = 1
			operands[i-4] = tmp
		}
	}
	self.fixA(self.pc0(), a)
	self.freeTmps(nTmps)
}

func (self *cg) concatExp(exp *BinopExp, a int) {
	nTmps := -1
	for {
		if nTmps < 0 {
			nTmps = 0
			self.exp(exp.Exp1, a, 1)
		} else {
			tmp := self.allocTmp()
			nTmps++
			self.exp(exp.Exp1, tmp, 1)
		}

		exp2, ok := exp.Exp2.(*BinopExp)
		if ok && exp2.Op == TOKEN_OP_CONCAT {
			exp = exp2
		} else {
			tmp := self.allocTmp()
			nTmps++
			self.exp(exp.Exp2, tmp, 1)
			break
		}
	}

	self.freeTmps(nTmps)
	self.inst(exp.Line, OP_CONCAT, a, a, a+nTmps)
}

func (self *cg) arithAndBitwiseBinopExp(exp *BinopExp, a int) {
	var rkb, rkc int
	nTmps := 0

	// calc rkb
	if rkb = self.toRK(exp.Exp1); rkb < 0 {
		if self.isTmpVar(a) {
			rkb = a
		} else {
			rkb = self.allocTmp()
			nTmps++
		}
		self.exp(exp.Exp1, rkb, 1)
	}

	prec := exp.Prec
	for {
		rhs := exp.Exp2
		line := exp.Line
		opcode := arithAndBitwiseBinops[exp.Op]

		if exp2, ok := castToBinopExp(exp.Exp2, prec); ok {
			rhs = exp2.Exp1
			exp = exp2
		}

		// calc rkc
		if rkc = self.toRK(rhs); rkc < 0 {
			if !self.isTmpVar(rkb) && self.isTmpVar(a) {
				rkc = a
				self.exp(rhs, rkc, 1)
			} else {
				rkc = self.allocTmp()
				self.exp(rhs, rkc, 1)
			}
		}
		
		if !self.isTmpVar(rkb) && !self.isTmpVar(rkc) {
			if self.isTmpVar(a) {
				self.inst(line, opcode, a, rkb, rkc)
				rkb = a
			} else {
				tmp := self.allocTmp()
				nTmps++
				self.inst(line, opcode, tmp, rkb, rkc)
				rkb = tmp
			}
		} else if !self.isTmpVar(rkb) {
			self.inst(line, opcode, rkc, rkb, rkc)
			rkb = rkc
		} else {
			self.inst(line, opcode, rkb, rkb, rkc)
			self.freeTmp()
		}

		if rhs == exp.Exp2 {
			break
		}
	}

	self.freeTmps(nTmps)
	self.fixA(self.pc0(), a)
}

// todo: name
func (self *cg) relationalBinopExp(exp *BinopExp, a, n int) {
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

	if n == 0 {
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
	} else {
		switch exp.Op {
		case TOKEN_OP_EQ:
			self.inst(exp.Line, OP_EQ, 1, iX, iY)
		case TOKEN_OP_NE:
			self.inst(exp.Line, OP_EQ, 0, iX, iY)
		case TOKEN_OP_LT:
			self.inst(exp.Line, OP_LT, 1, iX, iY)
		case TOKEN_OP_GT:
			self.inst(exp.Line, OP_LT, 1, iY, iX)
		case TOKEN_OP_LE:
			self.inst(exp.Line, OP_LE, 1, iX, iY)
		case TOKEN_OP_GE:
			self.inst(exp.Line, OP_LE, 1, iY, iX)
		}
		self.inst(exp.Line, OP_JMP, 0, 1, 0)
		self.inst(exp.Line, OP_LOADBOOL, a, 0, 1)
		self.inst(exp.Line, OP_LOADBOOL, a, 1, 0)
	}
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

func (self *cg) isLocVar(exp Exp) (int, bool) {
	if nameExp, ok := exp.(*NameExp); ok {
		if slot := self.slotOf(nameExp.Name); slot >= 0 {
			return slot, true
		}
	}
	return -1, false
}

// func (self *cg) isUpval(exp Exp) (int, bool) {
// 	if nameExp, ok := exp.(*NameExp); ok {
// 		if idx := self.lookupUpval(nameExp.Name); idx >= 0 {
// 			return idx, true
// 		}
// 	}
// 	return -1, false
// }

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
