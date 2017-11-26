package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

var _statEmpty = &EmptyStat{}

/*
stat ::=  ‘;’
	| break
	| ‘::’ Name ‘::’
	| goto Name
	| do block end
	| while exp do block end
	| repeat block until exp
	| if exp then block {elseif exp then block} [else block] end
	| for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
	| for namelist in explist do block end
	| function funcname funcbody
	| local function Name funcbody
	| local namelist [‘=’ explist]
	| varlist ‘=’ explist
	| functioncall
*/
func parseStat(lexer *Lexer) Stat {
	switch lexer.LookAhead(1) {
	case TOKEN_SEP_SEMI:
		return parseEmptyStat(lexer)
	case TOKEN_KW_BREAK:
		return parseBreakStat(lexer)
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

// for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
// for namelist in explist do block end
func parseForStat(lexer *Lexer) Stat {
	if lexer.LookAhead(3) == TOKEN_OP_ASSIGN {
		return parseForNumStat(lexer)
	} else {
		return parseForInStat(lexer)
	}
}

// local function Name funcbody
// local namelist [‘=’ explist]
func parseLocalAssignOrFuncDefStat(lexer *Lexer) Stat {
	if lexer.LookAhead(2) == TOKEN_KW_FUNCTION {
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
		return fc
	} else {
		lexer.Restore(backup)
		return parseAssignStat(lexer)
	}
}

// ;
func parseEmptyStat(lexer *Lexer) *EmptyStat {
	lexer.NextTokenOfKind(TOKEN_SEP_SEMI)
	return _statEmpty
}

// break
func parseBreakStat(lexer *Lexer) *BreakStat {
	lexer.NextTokenOfKind(TOKEN_KW_BREAK)
	return &BreakStat{lexer.Line()}
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
	lexer.NextTokenOfKind(TOKEN_KW_GOTO) // goto
	_, name := lexer.NextIdentifier()    // name
	return &GotoStat{name}
}

// do block end
func parseDoStat(lexer *Lexer) *DoStat {
	lexer.NextTokenOfKind(TOKEN_KW_DO)  // do
	block := parseBlock(lexer)          // block
	lexer.NextTokenOfKind(TOKEN_KW_END) // end
	return &DoStat{block}
}

// while exp do block end
func parseWhileStat(lexer *Lexer) *WhileStat {
	lexer.NextTokenOfKind(TOKEN_KW_WHILE) // while
	exp := parseExp(lexer)                // exp
	lexer.NextTokenOfKind(TOKEN_KW_DO)    // do
	block := parseBlock(lexer)            // block
	lexer.NextTokenOfKind(TOKEN_KW_END)   // end
	return &WhileStat{exp, block}
}

// repeat block until exp
func parseRepeatStat(lexer *Lexer) *RepeatStat {
	lexer.NextTokenOfKind(TOKEN_KW_REPEAT) // repeat
	block := parseBlock(lexer)             // block
	lexer.NextTokenOfKind(TOKEN_KW_UNTIL)  // until
	exp := parseExp(lexer)                 // exp
	return &RepeatStat{block, exp}
}

// if exp then block {elseif exp then block} [else block] end
func parseIfStat(lexer *Lexer) *IfStat {
	exps := make([]Exp, 0, 4)
	blocks := make([]*Block, 0, 4)

	lexer.NextTokenOfKind(TOKEN_KW_IF)         // if
	exps = append(exps, parseExp(lexer))       // exp
	lexer.NextTokenOfKind(TOKEN_KW_THEN)       // then
	blocks = append(blocks, parseBlock(lexer)) // block

	for lexer.LookAhead(1) == TOKEN_KW_ELSEIF {
		lexer.NextToken()                          // elseif
		exps = append(exps, parseExp(lexer))       // exp
		lexer.NextTokenOfKind(TOKEN_KW_THEN)       // then
		blocks = append(blocks, parseBlock(lexer)) // block
	}

	// else block => elseif true then block
	if lexer.LookAhead(1) == TOKEN_KW_ELSE {
		lexer.NextToken()                           // else
		exps = append(exps, &TrueExp{lexer.Line()}) //
		blocks = append(blocks, parseBlock(lexer))  // block
	}

	lexer.NextTokenOfKind(TOKEN_KW_END) // end
	return &IfStat{exps, blocks}
}

// for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
func parseForNumStat(lexer *Lexer) *ForNumStat {
	lineOfFor, _ := lexer.NextTokenOfKind(TOKEN_KW_FOR) // for
	_, varName := lexer.NextIdentifier()                // name
	lexer.NextTokenOfKind(TOKEN_OP_ASSIGN)              // =
	initExp := parseExp(lexer)                          // exp
	lexer.NextTokenOfKind(TOKEN_SEP_COMMA)              // ,
	limitExp := parseExp(lexer)                         // exp

	var stepExp Exp
	if lexer.LookAhead(1) == TOKEN_SEP_COMMA {
		lexer.NextToken()         // ,
		stepExp = parseExp(lexer) // exp
	} else {
		stepExp = &IntegerExp{lexer.Line(), 1}
	}

	lineOfDo, _ := lexer.NextTokenOfKind(TOKEN_KW_DO) // do
	block := parseBlock(lexer)                        // block
	lexer.NextTokenOfKind(TOKEN_KW_END)               // end

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
func parseForInStat(lexer *Lexer) *ForInStat {
	lexer.NextTokenOfKind(TOKEN_KW_FOR)               // for
	nameList := _parseNameList(lexer)                 // namelist
	lexer.NextTokenOfKind(TOKEN_KW_IN)                // in
	expList := parseExpList(lexer)                    // explist
	lineOfDo, _ := lexer.NextTokenOfKind(TOKEN_KW_DO) // do
	block := parseBlock(lexer)                        // block
	lexer.NextTokenOfKind(TOKEN_KW_END)               // end

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
	_, name := lexer.NextIdentifier()
	names = append(names, name)
	for lexer.LookAhead(1) == TOKEN_SEP_COMMA {
		lexer.NextToken()
		_, name := lexer.NextIdentifier()
		names = append(names, name)
	}
	return names
}

// local namelist [‘=’ explist]
func parseLocalAssignStat(lexer *Lexer) *LocalAssignStat {
	lexer.NextTokenOfKind(TOKEN_KW_LOCAL) // local
	nameList := _parseNameList(lexer)     // namelist
	var expList []Exp = nil
	if lexer.LookAhead(1) == TOKEN_OP_ASSIGN {
		lexer.NextToken()             // ==
		expList = parseExpList(lexer) // explist
	}

	return &LocalAssignStat{
		LastLine: lexer.Line(),
		NameList: nameList,
		ExpList:  expList,
	}
}

// varlist ‘=’ explist |
func parseAssignStat(lexer *Lexer) *AssignStat {
	varList := _parseVarList(lexer)        // varlist
	lexer.NextTokenOfKind(TOKEN_OP_ASSIGN) // =
	expList := parseExpList(lexer)         // explist
	return &AssignStat{
		LastLine: lexer.Line(),
		VarList:  varList,
		ExpList:  expList,
	}
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
	lexer.NextTokenOfKind(TOKEN_KW_LOCAL)    // local
	lexer.NextTokenOfKind(TOKEN_KW_FUNCTION) // function
	_, name := lexer.NextIdentifier()        // name
	fdExp := parseFuncDefExp(lexer)          // funcbody
	return &LocalFuncDefStat{name, fdExp}
}

// function funcname funcbody
// funcname ::= Name {‘.’ Name} [‘:’ Name]
// funcbody ::= ‘(’ [parlist] ‘)’ block end
// parlist ::= namelist [‘,’ ‘...’] | ‘...’
// namelist ::= Name {‘,’ Name}
func parseFuncDefStat(lexer *Lexer) *AssignStat {
	lexer.NextTokenOfKind(TOKEN_KW_FUNCTION)
	pexp, hasColon := _parseFuncName(lexer)
	fdExp := parseFuncDefExp(lexer)
	if hasColon { // insert self
		fdExp.ParList = append(fdExp.ParList, "self")
		copy(fdExp.ParList[1:], fdExp.ParList)
		fdExp.ParList[0] = "self"
	}

	return &AssignStat{
		LastLine: fdExp.Line,
		VarList:  []Exp{pexp},
		ExpList:  []Exp{fdExp},
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
