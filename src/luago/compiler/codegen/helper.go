package codegen

import . "luago/compiler/ast"

// todo: rename
func isExpTrue(exp Exp) bool {
	switch exp.(type) {
	case *TrueExp,
		*IntegerExp, *FloatExp, *StringExp,
		*TableConstructorExp, FuncDefExp:
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
