package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"
import "luago/number"

// explist ::= exp {‘,’ exp}
func parseExpList(lexer *Lexer) []Exp {
	exps := make([]Exp, 0, 4)
	exps = append(exps, parseExp(lexer))
	for lexer.LookAhead(1) == TOKEN_SEP_COMMA {
		lexer.NextToken()
		exps = append(exps, parseExp(lexer))
	}
	return exps
}

/*
exp ::=  nil | false | true | Numeral | LiteralString | ‘...’ | functiondef |
	 prefixexp | tableconstructor | exp binop exp | unop exp
*/
func parseExp(lexer *Lexer) Exp {
	return parseExp12(lexer)
}

// x or y
func parseExp12(lexer *Lexer) Exp {
	exp := parseExp11(lexer)
	for lexer.LookAhead(1) == TOKEN_OP_OR {
		line, op, _ := lexer.NextToken()
		lor := &BinopExp{line, op, exp, parseExp11(lexer)}
		exp = optimizeLogicalOr(lor)
	}
	return exp
}

// x and y
func parseExp11(lexer *Lexer) Exp {
	exp := parseExp10(lexer)
	for lexer.LookAhead(1) == TOKEN_OP_AND {
		line, op, _ := lexer.NextToken()
		land := &BinopExp{line, op, exp, parseExp10(lexer)}
		last := lexer.LookAhead(1) != TOKEN_OP_AND // todo
		exp = optimizeLogicalAnd(land, last)
	}
	return exp
}

// compare
func parseExp10(lexer *Lexer) Exp {
	exp := parseExp9(lexer)
	for {
		switch lexer.LookAhead(1) {
		case TOKEN_OP_LT, TOKEN_OP_GT, TOKEN_OP_NE,
			TOKEN_OP_LE, TOKEN_OP_GE, TOKEN_OP_EQ:
			line, op, _ := lexer.NextToken()
			exp = &BinopExp{line, op, exp, parseExp9(lexer)}
		default:
			return exp
		}
	}
	return exp
}

// x | y
func parseExp9(lexer *Lexer) Exp {
	exp := parseExp8(lexer)
	for lexer.LookAhead(1) == TOKEN_OP_BOR {
		line, op, _ := lexer.NextToken()
		bor := &BinopExp{line, op, exp, parseExp8(lexer)}
		exp = optimizeBitwiseBinaryOp(bor)
	}
	return exp
}

// x ~ y
func parseExp8(lexer *Lexer) Exp {
	exp := parseExp7(lexer)
	for lexer.LookAhead(1) == TOKEN_OP_BXOR {
		line, op, _ := lexer.NextToken()
		bxor := &BinopExp{line, op, exp, parseExp7(lexer)}
		exp = optimizeBitwiseBinaryOp(bxor)
	}
	return exp
}

// x & y
func parseExp7(lexer *Lexer) Exp {
	exp := parseExp6(lexer)
	for lexer.LookAhead(1) == TOKEN_OP_BAND {
		line, op, _ := lexer.NextToken()
		band := &BinopExp{line, op, exp, parseExp6(lexer)}
		exp = optimizeBitwiseBinaryOp(band)
	}
	return exp
}

// shift
func parseExp6(lexer *Lexer) Exp {
	exp := parseExp5(lexer)
	for {
		switch lexer.LookAhead(1) {
		case TOKEN_OP_SHL, TOKEN_OP_SHR:
			line, op, _ := lexer.NextToken()
			shx := &BinopExp{line, op, exp, parseExp5(lexer)}
			exp = optimizeBitwiseBinaryOp(shx)
		default:
			return exp
		}
	}
	return exp
}

// a .. b
func parseExp5(lexer *Lexer) Exp {
	line := 0
	exps := make([]Exp, 0, 2)

	exps = append(exps, parseExp4(lexer))
	for lexer.LookAhead(1) == TOKEN_OP_CONCAT {
		line, _, _ = lexer.NextToken()
		exps = append(exps, parseExp4(lexer))
	}

	if len(exps) > 1 {
		return &ConcatExp{line, exps}	
	} else {
		return exps[0]
	}
}

// x +/- y
func parseExp4(lexer *Lexer) Exp {
	exp := parseExp3(lexer)
	for {
		switch lexer.LookAhead(1) {
		case TOKEN_OP_ADD, TOKEN_OP_SUB:
			line, op, _ := lexer.NextToken()
			arith := &BinopExp{line, op, exp, parseExp3(lexer)}
			exp = optimizeArithBinaryOp(arith)
		default:
			return exp
		}
	}
	return exp
}

// *, %, /, //
func parseExp3(lexer *Lexer) Exp {
	exp := parseExp2(lexer)
	for {
		switch lexer.LookAhead(1) {
		case TOKEN_OP_MUL, TOKEN_OP_MOD, TOKEN_OP_DIV, TOKEN_OP_IDIV:
			line, op, _ := lexer.NextToken()
			arith := &BinopExp{line, op, exp, parseExp2(lexer)}
			exp = optimizeArithBinaryOp(arith)
		default:
			return exp
		}
	}
	return exp
}

// unary
func parseExp2(lexer *Lexer) Exp {
	switch lexer.LookAhead(1) {
	case TOKEN_OP_UNM, TOKEN_OP_BNOT, TOKEN_OP_LEN, TOKEN_OP_NOT:
		line, op, _ := lexer.NextToken()
		exp := &UnopExp{line, op, parseExp2(lexer)}
		return optimizeUnaryOp(exp)
	default:
		return parseExp1(lexer)
	}
}

// x ^ y
func parseExp1(lexer *Lexer) Exp { // pow is right associative
	exp := parseExp0(lexer)
	for lexer.LookAhead(1) == TOKEN_OP_POW {
		line, op, _ := lexer.NextToken()
		exp = &BinopExp{line, op, exp, parseExp1(lexer)}
	}
	return optimizePow(exp)
}

func parseExp0(lexer *Lexer) Exp {
	switch lexer.LookAhead(1) {
	case TOKEN_VARARG: // ...
		line, _, _ := lexer.NextToken()
		return &VarargExp{line}
	case TOKEN_KW_NIL: // nil
		line, _, _ := lexer.NextToken()
		return &NilExp{line}
	case TOKEN_KW_TRUE: // true
		line, _, _ := lexer.NextToken()
		return &TrueExp{line}
	case TOKEN_KW_FALSE: // false
		line, _, _ := lexer.NextToken()
		return &FalseExp{line}
	case TOKEN_STRING: // LiteralString
		line, _, token := lexer.NextToken()
		return &StringExp{line, token}
	case TOKEN_NUMBER: // Numeral
		return parseNumberExp(lexer, 1)
	case TOKEN_SEP_LCURLY: // tableconstructor
		return parseTableConstructorExp(lexer)
	case TOKEN_KW_FUNCTION: // functiondef
		lexer.NextToken()
		return parseFuncDefExp(lexer)
	default: // prefixexp
		return parsePrefixExp(lexer)
	}
}

func parseNumberExp(lexer *Lexer, sign int) Exp {
	line, _, token := lexer.NextToken()
	if i, ok := number.ParseInteger(token, 10); ok {
		if sign >= 0 {
			return &IntegerExp{line, i}
		} else {
			return &IntegerExp{line, -i}
		}
	} else if f, ok := number.ParseFloat(token); ok {
		if sign >= 0 {
			return &FloatExp{line, f}
		} else {
			return &FloatExp{line, -f}
		}
	} else { // todo
		panic("not a number: " + token)
	}
}
