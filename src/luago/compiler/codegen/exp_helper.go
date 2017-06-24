package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

// todo: rename
func (self *codeGen) isLocVar(exp Exp) (int, bool) {
	if nameExp, ok := exp.(*NameExp); ok {
		if slot := self.slotOf(nameExp.Name); slot >= 0 {
			return slot, true
		}
	}
	return -1, false
}

func isNameExp(exp Exp) bool {
	_, ok := exp.(*NameExp)
	return ok
}

// todo: rename
func isExpTrue(exp Exp) bool {
	switch exp.(type) {
	case *TrueExp,
		*IntegerExp, *FloatExp, *StringExp,
		//*TableConstructorExp,
		FuncDefExp:
		return true
	default:
		return false
	}
}

func isTrueAtCompileTime(exp Exp) bool {
	switch exp.(type) {
	case *TrueExp, *IntegerExp, *FloatExp, *StringExp, FuncDefExp:
		return true
	default:
		return false
	}
}

func isVarargOrFuncCallExp(exp Exp) bool {
	switch exp.(type) {
	case *VarargExp, *FuncCallExp:
		return true
	default:
		return false
	}
}

func castNilToFalse(exp Exp) Exp {
	if nilExp, ok := exp.(*NilExp); ok {
		return &FalseExp{nilExp.Line}
	} else {
		return exp
	}
}

func castToRelationalBinopExp(exp Exp) (*BinopExp, bool) {
	if bexp, ok := exp.(*BinopExp); ok {
		switch bexp.Op {
		case TOKEN_OP_EQ, TOKEN_OP_NE,
			TOKEN_OP_LT, TOKEN_OP_LE,
			TOKEN_OP_GT, TOKEN_OP_GE:
			return bexp, true
		}
	}
	return nil, false
}

func castToBinopExp(exp Exp, op int) (*BinopExp, bool) {
	if bexp, ok := exp.(*BinopExp); ok {
		if bexp.Op == op {
			return bexp, true
		}
	}
	return nil, false
}

func castToConcatExp(exp Exp) (*BinopExp, bool) {
	if bexp, ok := exp.(*BinopExp); ok {
		if bexp.Op == TOKEN_OP_CONCAT {
			return bexp, true
		}
	}
	return nil, false
}

// todo
func lineOfExp(exp Exp) int {
	switch x := exp.(type) {
	case *NilExp:
		return x.Line
	case *TrueExp:
		return x.Line
	case *FalseExp:
		return x.Line
	case *IntegerExp:
		return x.Line
	case *FloatExp:
		return x.Line
	case *StringExp:
		return x.Line
	case *VarargExp:
		return x.Line
	case *NameExp:
		return x.Line
	case *FuncDefExp:
		return x.Line
	case *BinopExp: // todo
		return x.Line
	case *UnopExp: // todo
		return x.Line
	case *TableConstructorExp:
		return 0 // todo
	case *BracketsExp:
		return x.Line
	case *FuncCallExp:
		return x.Line
	}

	panic("todo!")
}

// todo
func lastLineOfExp(exp Exp) int {
	switch x := exp.(type) {
	case *NilExp:
		return x.Line
	case *TrueExp:
		return x.Line
	case *FalseExp:
		return x.Line
	case *IntegerExp:
		return x.Line
	case *FloatExp:
		return x.Line
	case *StringExp:
		return x.Line
	case *VarargExp:
		return x.Line
	case *NameExp:
		return x.Line
	case *FuncDefExp:
		return x.LastLine
	case *BinopExp: // todo
		return x.Line
	case *UnopExp: // todo
		return x.Line
	case *TableConstructorExp:
		return x.LastLine
	case *BracketsExp:
		return x.Line
	case *FuncCallExp:
		return x.Line
	}

	panic("todo!")
}
