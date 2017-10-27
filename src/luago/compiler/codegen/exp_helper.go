package codegen

import . "luago/compiler/ast"

func isVarargOrFuncCallExp(exp Exp) bool {
	switch exp.(type) {
	case *VarargExp, *FuncCallExp:
		return true
	default:
		return false
	}
}

func removeTailNils(exps []Exp) []Exp {
	for {
		if n := len(exps); n > 0 {
			if _, ok := exps[n-1].(*NilExp); ok {
				exps = exps[:n-1]
				continue
			}
		}
		break
	}
	return exps
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
	case *ConcatExp: // todo
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
	case *ConcatExp: // todo
		return x.Line
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
