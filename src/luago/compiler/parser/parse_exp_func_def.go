package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

// functiondef ::= function funcbody
// funcbody ::= ‘(’ [parlist] ‘)’ block end
func parseFuncDefExp(lexer *Lexer) *FuncDefExp {
	/* keyword function is scanned */
	line := lexer.Line()
	lexer.NextTokenOfKind(TOKEN_SEP_LPAREN)
	names, isVararg := parseParList(lexer)
	lexer.NextTokenOfKind(TOKEN_SEP_RPAREN)

	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_KW_END)
	lastLine := lexer.Line()

	return &FuncDefExp{line, lastLine, names, isVararg, true, block}
}

// parlist ::= namelist [‘,’ ‘...’] | ‘...’
func parseParList(lexer *Lexer) (names []string, isVararg bool) {
	names = make([]string, 0, 8)
	isVararg = false

	for {
		switch lexer.LookAhead(1) {
		case TOKEN_IDENTIFIER:
			_, name := lexer.NextIdentifier()
			names = append(names, name)
		case TOKEN_VARARG:
			lexer.NextToken()
			isVararg = true
			return
		}

		if lexer.LookAhead(1) == TOKEN_SEP_COMMA {
			lexer.NextToken()
		} else {
			break
		}
	}

	return
}
