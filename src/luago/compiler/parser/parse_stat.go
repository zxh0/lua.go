package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

var _statEmpty = &EmptyStat{}

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
		line, _, _ := lexer.NextToken()
		return &BreakStat{line}
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
		return parseForStat(lexer)
	case TOKEN_KW_FUNCTION:
		return parseFuncDefStat(lexer)
	case TOKEN_KW_LOCAL:
		return parseLocalAssignOrFuncDefStat(lexer)
	default:
		return parseAssignOrFuncCallStat(lexer)
	}
}

// label
func parseLabelStat(lexer *Lexer) *LabelStat {
	lexer.NextTokenOfKind(TOKEN_SEP_LABEL)
	_, name := lexer.NextIdentifier()
	lexer.NextTokenOfKind(TOKEN_SEP_LABEL)
	return &LabelStat{name}
}

// goto Name
func parseGotoStat(lexer *Lexer) *GotoStat {
	lexer.NextTokenOfKind(TOKEN_KW_GOTO)
	_, name := lexer.NextIdentifier()
	return &GotoStat{name}
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
	lexer.NextTokenOfKind(TOKEN_KW_WHILE)
	exp := parseExp(lexer)
	lexer.NextTokenOfKind(TOKEN_KW_DO)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_KW_END)
	return &WhileStat{exp, block}
}

// repeat block until exp
func parseRepeatStat(lexer *Lexer) *RepeatStat {
	lexer.NextTokenOfKind(TOKEN_KW_REPEAT)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_KW_UNTIL)
	exp := parseExp(lexer)
	return &RepeatStat{block, exp}
}

// for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
// for namelist in explist do block end
func parseForStat(lexer *Lexer) Stat {
	line, _ := lexer.NextTokenOfKind(TOKEN_KW_FOR)
	if lexer.LookAhead(2) == TOKEN_OP_ASSIGN {
		return parseForNumStat(lexer, line)
	} else {
		return parseForInStat(lexer, line)
	}
}

// for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
func parseForNumStat(lexer *Lexer, lineOfFor int) *ForNumStat {
	_, varName := lexer.NextIdentifier()
	lexer.NextTokenOfKind(TOKEN_OP_ASSIGN)
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
	nameList := _parseNameList(lexer)

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

// namelist ::= Name {‘,’ Name}
func _parseNameList(lexer *Lexer) []string {
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

// if exp then block {elseif exp then block} [else block] end
func parseIfStat(lexer *Lexer) *IfStat {
	stat := &IfStat{
		Exps:   make([]Exp, 0, 8),
		Blocks: make([]*Block, 0, 8),
	}
	_parseIf(lexer, stat)
	_parseElseIf(lexer, stat)
	_parseElse(lexer, stat)
	return stat
}

// if exp then block
func _parseIf(lexer *Lexer, stat *IfStat) {
	lexer.NextTokenOfKind(TOKEN_KW_IF)
	stat.Exps = append(stat.Exps, parseExp(lexer))

	lexer.NextTokenOfKind(TOKEN_KW_THEN)
	stat.Blocks = append(stat.Blocks, parseBlock(lexer))
}

// {elseif exp then block}
func _parseElseIf(lexer *Lexer, stat *IfStat) {
	for lexer.LookAhead(1) == TOKEN_KW_ELSEIF {
		lexer.NextTokenOfKind(TOKEN_KW_ELSEIF)
		stat.Exps = append(stat.Exps, parseExp(lexer))

		lexer.NextTokenOfKind(TOKEN_KW_THEN)
		stat.Blocks = append(stat.Blocks, parseBlock(lexer))
	}
}

// [else block] end
func _parseElse(lexer *Lexer, stat *IfStat) {
	if lexer.LookAhead(1) == TOKEN_KW_ELSE {
		line, _ := lexer.NextTokenOfKind(TOKEN_KW_ELSE)

		// else block => elseif true then block
		stat.Exps = append(stat.Exps, &TrueExp{line})
		stat.Blocks = append(stat.Blocks, parseBlock(lexer))
	}

	lexer.NextTokenOfKind(TOKEN_KW_END)
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

// local namelist [‘=’ explist]
func parseLocalAssignStat(lexer *Lexer) *LocalAssignStat {
	/* keyword local is scanned */
	stat := &LocalAssignStat{}
	stat.NameList = _parseNameList(lexer)
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
	stat.VarList = _parseVarList(lexer)
	lexer.NextTokenOfKind(TOKEN_OP_ASSIGN)
	stat.ExpList = parseExpList(lexer)
	stat.LastLine = lexer.Line()
	return stat
}

// varlist ::= var {‘,’ var}
func _parseVarList(lexer *Lexer) []Exp {
	vars := make([]Exp, 0, 8)
	vars = append(vars, _parseVar(lexer))
	for lexer.LookAhead(1) == TOKEN_SEP_COMMA {
		lexer.NextToken()
		vars = append(vars, _parseVar(lexer))
	}
	return vars
}

// var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
func _parseVar(lexer *Lexer) Exp {
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

/*
http://www.lua.org/manual/5.3/manual.html#3.4.11

function f() end          =>  f = function() end
function t.a.b.c.f() end  =>  t.a.b.c.f = function() end
function t.a.b.c:f() end  =>  t.a.b.c.f = function(self) end
local function f() end    =>  local f; f = function() end

The statement `local function f () body end`
translates to `local f; f = function () body end`
not to `local f = function () body end`
(This only makes a difference when the body of the function
 contains references to f.)
*/

// local function Name funcbody
func parseLocalFuncDefStat(lexer *Lexer) *LocalFuncDefStat {
	/* keyword local is passed */
	lexer.NextTokenOfKind(TOKEN_KW_FUNCTION)
	_, name := lexer.NextIdentifier()
	fdExp := parseFuncDefExp(lexer)

	return &LocalFuncDefStat{
		Name: name,
		Exp:  fdExp,
	}
}

// function funcname funcbody
// funcname ::= Name {‘.’ Name} [‘:’ Name]
// funcbody ::= ‘(’ [parlist] ‘)’ block end
// parlist ::= namelist [‘,’ ‘...’] | ‘...’
// namelist ::= Name {‘,’ Name}
func parseFuncDefStat(lexer *Lexer) *AssignStat {
	lexer.NextTokenOfKind(TOKEN_KW_FUNCTION)
	pexp, hasColon := _parseFuncName(lexer)
	funcDef := parseFuncDefExp(lexer)
	if hasColon { // insert self
		funcDef.ParList = append(funcDef.ParList, "self")
		copy(funcDef.ParList[1:], funcDef.ParList)
		funcDef.ParList[0] = "self"
	}

	return &AssignStat{
		LastLine: funcDef.Line,
		VarList:  []Exp{pexp},
		ExpList:  []Exp{funcDef},
	}
}

// funcname ::= Name {‘.’ Name} [‘:’ Name]
func _parseFuncName(lexer *Lexer) (Exp, bool) {
	var exp Exp
	hasColon := false

	line, name := lexer.NextIdentifier()
	exp = &NameExp{line, name}

	for lexer.LookAhead(1) == TOKEN_SEP_DOT {
		lexer.NextToken()
		line, name := lexer.NextIdentifier()
		idx := &StringExp{line, name}
		exp = &TableAccessExp{line, exp, idx}
	}
	if lexer.LookAhead(1) == TOKEN_SEP_COLON {
		lexer.NextToken()
		line, name := lexer.NextIdentifier()
		idx := &StringExp{line, name}
		exp = &TableAccessExp{line, exp, idx}
		hasColon = true
	}

	return exp, hasColon
}
