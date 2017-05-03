package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

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
func parseLocalFuncDefStat(lexer *Lexer) *LocalAssignStat {
	/* keyword local is passed */
	lexer.NextTokenOfKind(TOKEN_KW_FUNCTION)
	_, name := lexer.NextIdentifier()
	funcDef := parseFuncDefExp(lexer)
	funcDef.IsAno = false

	return &LocalAssignStat{
		LastLine: funcDef.Line,
		NameList: []string{name},
		ExpList:  []Exp{funcDef},
	}
}

// function funcname funcbody
// funcname ::= Name {‘.’ Name} [‘:’ Name]
// funcbody ::= ‘(’ [parlist] ‘)’ block end
// parlist ::= namelist [‘,’ ‘...’] | ‘...’
// namelist ::= Name {‘,’ Name}
func parseFuncDefStat(lexer *Lexer) *AssignStat {
	lexer.NextTokenOfKind(TOKEN_KW_FUNCTION)
	pexp, hasColon := parseFuncName(lexer)
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
func parseFuncName(lexer *Lexer) (Exp, bool) {
	var exp Exp
	hasColon := false

	line, name := lexer.NextIdentifier()
	exp = &NameExp{line, name}

	for lexer.LookAhead(1) == TOKEN_SEP_DOT {
		lexer.NextToken()
		line, name := lexer.NextIdentifier()
		idx := &StringExp{line, name}
		exp = &BracketsExp{line, exp, idx}
	}
	if lexer.LookAhead(1) == TOKEN_SEP_COLON {
		lexer.NextToken()
		line, name := lexer.NextIdentifier()
		idx := &StringExp{line, name}
		exp = &BracketsExp{line, exp, idx}
		hasColon = true
	}

	return exp, hasColon
}
