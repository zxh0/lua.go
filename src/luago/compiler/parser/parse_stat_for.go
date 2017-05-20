package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

// for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
// for namelist in explist do block end
func parseStatFor(lexer *Lexer) Stat {
	line, _ := lexer.NextTokenOfKind(TOKEN_KW_FOR)
	if lexer.LookAhead(2) == TOKEN_ASSIGN {
		return parseForNumStat(lexer, line)
	} else {
		return parseForInStat(lexer, line)
	}
}

// for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
func parseForNumStat(lexer *Lexer, lineOfFor int) *ForNumStat {
	_, varName := lexer.NextIdentifier()
	lexer.NextTokenOfKind(TOKEN_ASSIGN)
	initExp := parseExp(lexer)

	lexer.NextTokenOfKind(TOKEN_SEP_COMMA)
	limitExp := parseExp(lexer)

	var stepExp Exp
	if lexer.LookAhead(1) == TOKEN_SEP_COMMA {
		lexer.NextToken()
		stepExp = parseExp(lexer)
	} else {
		stepExp = &IntegerExp{lexer.Line(), 1}
	}

	lineOfDo, _ := lexer.NextTokenOfKind(TOKEN_KW_DO)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_KW_END)

	return &ForNumStat{
		LineOfFor: lineOfFor,
		LineOfDo:  lineOfDo,
		VarName:   varName,
		InitExp:   initExp,
		LimitExp:  limitExp,
		StepExp:   stepExp,
		Block:     block,
	}
}

// for namelist in explist do block end
// namelist ::= Name {‘,’ Name}
// explist ::= exp {‘,’ exp}
func parseForInStat(lexer *Lexer, line int) *ForInStat {
	nameList := parseNameList(lexer)

	lexer.NextTokenOfKind(TOKEN_KW_IN)
	expList := parseExpList(lexer)

	lineOfDo, _ := lexer.NextTokenOfKind(TOKEN_KW_DO)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_KW_END)

	return &ForInStat{
		LineOfDo: lineOfDo,
		NameList: nameList,
		ExpList:  expList,
		Block:    block,
	}
}
