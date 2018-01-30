package parser

import "fmt"
import . "luago/compiler/ast"
import . "luago/compiler/lexer"

func blockToString(block *Block) string {
	str := ""
	if len(block.Stats) > 0 {
		for i, stat := range block.Stats {
			str += statToString(stat)
			if i < len(block.Stats)-1 {
				str += " "
			}
		}
	}
	if block.RetExps != nil {
		str += "return"
		for _, exp := range block.RetExps {
			str += " " + expToString(exp)
		}
	}
	return str
}

func statToString(stat Stat) string {
	switch x := stat.(type) {
	case *EmptyStat:
		return ";"
	case *BreakStat:
		return "break"
	case *LabelStat:
		return "::" + x.Name + "::"
	case *GotoStat:
		return "goto " + x.Name
	case *DoStat:
		return "do " + blockToString(x.Block) + " end"
	case *WhileStat:
		return "while " + expToString(x.Exp) +
			" do " + blockToString(x.Block) + " end"
	case *RepeatStat:
		return "repeat " + blockToString(x.Block) +
			" until " + expToString(x.Exp)
	case *FuncCallStat:
		return funcCallToString(x)
	case *IfStat:
		return ifStatToString(x)
	case *ForNumStat:
		return forNumStatToString(x)
	case *ForInStat:
		return forInStatToString(x)
	case *AssignStat:
		return assignStatToString(x)
	case *LocalAssignStat:
		return localAssignStatToString(x)
	case *LocalFuncDefStat:
		return "local " + funcDefExpToString(x.Exp, x.Name)
	}
	panic("todo!")
}

func ifStatToString(stat *IfStat) string {
	str := "if " + expToString(stat.Exps[0]) +
		" then " + blockToString(stat.Blocks[0])
	for i := 1; i < len(stat.Exps); i++ {
		str += " elseif " + expToString(stat.Exps[i])
		str += " then " + blockToString(stat.Blocks[i])
	}
	str += " end"
	return str
}

func forNumStatToString(stat *ForNumStat) string {
	str := "for " + stat.VarName + " = " + expToString(stat.InitExp) +
		", " + expToString(stat.LimitExp)
	if stat.StepExp != nil {
		str += ", " + expToString(stat.StepExp)
	}
	str += " do " + blockToString(stat.Block) + " end"
	return str
}

func forInStatToString(stat *ForInStat) string {
	panic("todo!")
}

func assignStatToString(stat *AssignStat) string {
	str := ""
	for _, pexp := range stat.VarList {
		str += expToString(pexp)
	}
	str += " = "
	for _, exp := range stat.ExpList {
		str += expToString(exp)
	}
	return str
}

func localAssignStatToString(stat *LocalAssignStat) string {
	str := "local "
	for i, name := range stat.NameList {
		str += name
		if i < len(stat.NameList)-1 {
			str += ", "
		}
	}
	str += " = "
	for i, exp := range stat.ExpList {
		str += expToString(exp)
		if i < len(stat.ExpList)-1 {
			str += ", "
		}
	}
	return str
}

func expToString(exp Exp) string {
	switch x := exp.(type) {
	case *NilExp:
		return "nil"
	case *TrueExp:
		return "true"
	case *FalseExp:
		return "false"
	case *VarargExp:
		return "..."
	case *IntegerExp:
		return fmt.Sprintf("%d", x.Val)
	case *FloatExp:
		return fmt.Sprintf("%f", x.Val)
	case *StringExp:
		return "'" + x.Str + "'"
	case *ConcatExp:
		str := ""
		for _, exp := range x.Exps {
			str += " .. "
			str += expToString(exp)
		}
		return str[4:]
	case *UnopExp:
		return unopToString(x.Op) + "(" + expToString(x.Exp) + ")"
	case *BinopExp:
		return "(" + expToString(x.Exp1) + binopToString(x.Op) + expToString(x.Exp2) + ")"
	case *TableConstructorExp:
		return tcExpToString(x)
	case *FuncDefExp:
		return funcDefExpToString(x, "")
	case *NameExp:
		return x.Name
	case *ParensExp:
		return "(" + expToString(x.Exp) + ")"
	case *TableAccessExp:
		return expToString(x.PrefixExp) + "[" + expToString(x.KeyExp) + "]"
	case *FuncCallExp:
		return funcCallToString(x)
	case int: // table/list index
		return fmt.Sprintf("%d", x)
	default:
		panic("unreachable!")
	}
}

func funcCallToString(fc *FuncCallExp) string {
	str := expToString(fc.PrefixExp)
	if fc.NameExp != nil {
		str = str + ":" + fc.NameExp.Str
	}
	str += "("
	for i, exp := range fc.Args {
		str += expToString(exp)
		if i < len(fc.Args)-1 {
			str += ", "
		}
	}
	str += ")"
	return str
}

func tcExpToString(tc *TableConstructorExp) string {
	str := "{"
	for i, k := range tc.KeyExps {
		v := tc.ValExps[i]
		if k != nil {
			str += "[" + expToString(k) + "]" + "="
		}
		str += expToString(v) + ","
	}
	str += "}"
	return str
}

func funcDefExpToString(fd *FuncDefExp, name string) string {
	str := "function"
	if name != "" {
		str += " " + name
	}
	str += "("
	for i, name := range fd.ParList {
		str += name
		if i < len(fd.ParList)-1 {
			str += ", "
		}
	}
	if fd.IsVararg {
		if len(fd.ParList) > 0 {
			str += ", ..."
		} else {
			str += "..."
		}
	}
	str += ") end" // todo
	return str
}

func unopToString(op int) string {
	switch op {
	case TOKEN_OP_UNM:
		return "-"
	case TOKEN_OP_LEN:
		return "#"
	case TOKEN_OP_BNOT:
		return "~"
	case TOKEN_OP_NOT:
		return "not "
	default:
		panic("unreachable!")
	}
}

func binopToString(op int) string {
	switch op {
	case TOKEN_OP_ADD:
		return " + "
	case TOKEN_OP_SUB:
		return " - "
	case TOKEN_OP_MUL:
		return " * "
	case TOKEN_OP_DIV:
		return " / "
	case TOKEN_OP_IDIV:
		return " // "
	case TOKEN_OP_POW:
		return " ^ "
	case TOKEN_OP_MOD:
		return " % "
	case TOKEN_OP_BAND:
		return " & "
	case TOKEN_OP_BXOR:
		return " ~ "
	case TOKEN_OP_BOR:
		return " | "
	case TOKEN_OP_SHR:
		return " >> "
	case TOKEN_OP_SHL:
		return " << "
	case TOKEN_OP_CONCAT:
		return " .. "
	case TOKEN_OP_LT:
		return " < "
	case TOKEN_OP_LE:
		return " <= "
	case TOKEN_OP_GT:
		return " > "
	case TOKEN_OP_GE:
		return " >= "
	case TOKEN_OP_EQ:
		return " == "
	case TOKEN_OP_NE:
		return " ~= "
	case TOKEN_OP_AND:
		return " and "
	case TOKEN_OP_OR:
		return " or "
	default:
		panic("unreachable!")
	}
}
