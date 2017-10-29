package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

// local namelist [‘=’ explist]
func parseLocalAssignStat(lexer *Lexer) *LocalAssignStat {
	/* keyword local is scanned */
	stat := &LocalAssignStat{}
	stat.NameList = parseNameList(lexer)
	if lexer.LookAhead(1) == TOKEN_OP_ASSIGN {
		lexer.NextToken()
		stat.ExpList = parseExpList(lexer)
	}
	stat.LastLine = lexer.Line()
	return stat
}

// varlist ‘=’ explist |
func parseAssignStat(lexer *Lexer) *AssignStat {
	stat := &AssignStat{}
	stat.VarList = parseVarList(lexer)
	lexer.NextTokenOfKind(TOKEN_OP_ASSIGN)
	stat.ExpList = parseExpList(lexer)
	stat.LastLine = lexer.Line()
	return stat
}

// varlist ::= var {‘,’ var}
func parseVarList(lexer *Lexer) []Exp {
	vars := make([]Exp, 0, 8)
	vars = append(vars, parseVar(lexer))
	for lexer.LookAhead(1) == TOKEN_SEP_COMMA {
		lexer.NextToken()
		vars = append(vars, parseVar(lexer))
	}
	return vars
}

// var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
func parseVar(lexer *Lexer) Exp {
	backup := lexer.Backup()

	exp := parsePrefixExp(lexer)
	switch exp.(type) {
	case *NameExp, *TableAccessExp:
		return exp
	default:
		lexer.Restore(backup)
		lexer.NextTokenOfKind(-1) // trigger error
		panic("unreachable!")
	}
}
