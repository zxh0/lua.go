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
	case *FuncCallExp:
		return x.Line
	case *TableConstructorExp:
		return x.Line
	case *UnopExp:
		return x.Line
	case *TableAccessExp:
		return lineOfExp(x.PrefixExp)
	case *ConcatExp:
		return lineOfExp(x.Exps[0])
	case *BinopExp:
		return lineOfExp(x.Exp1)
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
	case *FuncCallExp:
		return x.LastLine
	case *TableConstructorExp:
		return x.LastLine
	case *TableAccessExp:
		return x.LastLine
	case *ConcatExp:
		return lastLineOfExp(x.Exps[len(x.Exps)-1])
	case *BinopExp:
		return lastLineOfExp(x.Exp2)
	case *UnopExp:
		return lastLineOfExp(x.Exp)
	}

	panic("todo!")
}
