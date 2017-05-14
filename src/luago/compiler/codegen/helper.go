package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

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

func isVarargOrFuncCallExp(exp Exp) bool {
	switch exp.(type) {
	case *VarargExp, *FuncCallExp:
		return true
	default:
		return false
	}
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
	case *VarargExp:
		return x.Line
	case *IntegerExp:
		return x.Line
	case *FuncDefExp:
		return x.Line
	case *FloatExp:
		return x.Line
	case *StringExp:
		return x.Line
	case *BinopExp: // todo
		return x.Line
	case *UnopExp: // todo
		return x.Line
	case *TableConstructorExp:
		return 0 // todo
	case *NameExp:
		return x.Line
	case *BracketsExp:
		return x.Line
	case *FuncCallExp:
		return x.Line
	}

	panic("todo!")
}
