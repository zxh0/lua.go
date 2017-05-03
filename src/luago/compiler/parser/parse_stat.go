package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

var _statEmpty = &EmptyStat{}
var _statBreak = &BreakStat{}

/*
stat ::=  ‘;’ |
	 varlist ‘=’ explist |
	 functioncall |
	 label |
	 break |
	 goto Name |
	 do block end |
	 while exp do block end |
	 repeat block until exp |
	 if exp then block {elseif exp then block} [else block] end |
	 for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end |
	 for namelist in explist do block end |
	 function funcname funcbody |
	 local function Name funcbody |
	 local namelist [‘=’ explist]
*/
func parseStat(lexer *Lexer) Stat {
	switch lexer.LookAhead(1) {
	case TOKEN_SEP_SEMI:
		lexer.NextToken()
		return _statEmpty
	case TOKEN_KW_BREAK:
		lexer.NextToken()
		return _statBreak
	case TOKEN_SEP_LABEL:
		return parseLabelStat(lexer)
	case TOKEN_KW_GOTO:
		return parseGotoStat(lexer)
	case TOKEN_KW_DO:
		return parseDoStat(lexer)
	case TOKEN_KW_WHILE:
		return parseWhileStat(lexer)
	case TOKEN_KW_REPEAT:
		return parseRepeatStat(lexer)
	case TOKEN_KW_IF:
		return parseIfStat(lexer)
	case TOKEN_KW_FOR:
		return parseStatFor(lexer)
	case TOKEN_KW_FUNCTION:
		return parseFuncDefStat(lexer)
	case TOKEN_KW_LOCAL:
		return parseLocalAssignOrFuncDefStat(lexer)
	default:
		return parseAssignOrFuncCallStat(lexer)
	}
}

// label
func parseLabelStat(lexer *Lexer) LabelStat {
	lexer.NextTokenOfKind(TOKEN_SEP_LABEL)
	_, name := lexer.NextIdentifier()
	lexer.NextTokenOfKind(TOKEN_SEP_LABEL)
	return LabelStat(name)
}

// goto Name
func parseGotoStat(lexer *Lexer) GotoStat {
	lexer.NextTokenOfKind(TOKEN_KW_GOTO)
	_, name := lexer.NextIdentifier()
	return GotoStat(name)
}

// do block end
func parseDoStat(lexer *Lexer) DoStat {
	lexer.NextTokenOfKind(TOKEN_KW_DO)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_KW_END)
	return DoStat(block)
}

// while exp do block end
func parseWhileStat(lexer *Lexer) *WhileStat {
	line, _ := lexer.NextTokenOfKind(TOKEN_KW_WHILE)
	exp := parseExp(lexer)
	lexer.NextTokenOfKind(TOKEN_KW_DO)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_KW_END)
	return &WhileStat{line, exp, block}
}

// repeat block until exp
func parseRepeatStat(lexer *Lexer) *RepeatStat {
	line, _ := lexer.NextTokenOfKind(TOKEN_KW_REPEAT)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_KW_UNTIL)
	exp := parseExp(lexer)
	return &RepeatStat{line, block, exp}
}

// local function Name funcbody
// local namelist [‘=’ explist]
func parseLocalAssignOrFuncDefStat(lexer *Lexer) Stat {
	lexer.NextToken()
	if lexer.LookAhead(1) == TOKEN_KW_FUNCTION {
		return parseLocalFuncDefStat(lexer)
	} else {
		return parseLocalAssignStat(lexer)
	}
}

// varlist ‘=’ explist
// functioncall
func parseAssignOrFuncCallStat(lexer *Lexer) Stat {
	backup := lexer.Backup()
	prefixExp := parsePrefixExp(lexer)
	if fc, ok := prefixExp.(*FuncCallExp); ok {
		return FuncCallStat(fc)
	} else {
		lexer.Restore(backup)
		return parseAssignStat(lexer)
	}
}

// namelist ::= Name {‘,’ Name}
func parseNameList(lexer *Lexer) []string {
	names := make([]string, 0, 4)
	_, token := lexer.NextIdentifier()
	names = append(names, token)
	for lexer.LookAhead(1) == TOKEN_SEP_COMMA {
		lexer.NextToken()
		_, token := lexer.NextIdentifier()
		names = append(names, token)
	}
	return names
}
