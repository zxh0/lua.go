package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

// tableconstructor ::= ‘{’ [fieldlist] ‘}’
func parseTableConstructorExp(lexer *Lexer) *TableConstructorExp {
	tc := &TableConstructorExp{
		NArr:    0,
		KeyExps: make([]Exp, 0, 8),
		ValExps: make([]Exp, 0, 8),
	}

	tc.Line = lexer.Line()
	lexer.NextTokenOfKind(TOKEN_SEP_LCURLY)
	if lexer.LookAhead(1) != TOKEN_SEP_RCURLY {
		parseFieldList(lexer, tc)
	}
	lexer.NextTokenOfKind(TOKEN_SEP_RCURLY)
	tc.LastLine = lexer.Line()

	return tc
}

// fieldlist ::= field {fieldsep field} [fieldsep]
func parseFieldList(lexer *Lexer, tc *TableConstructorExp) {
	parseField(lexer, tc)

	for isFieldSep(lexer.LookAhead(1)) {
		lexer.NextToken()
		if lexer.LookAhead(1) == TOKEN_SEP_RCURLY {
			break
		} else {
			parseField(lexer, tc)
		}
	}
}

// fieldsep ::= ‘,’ | ‘;’
func isFieldSep(tokenKind int) bool {
	return tokenKind == TOKEN_SEP_COMMA || tokenKind == TOKEN_SEP_SEMI
}

// field ::= ‘[’ exp ‘]’ ‘=’ exp | Name ‘=’ exp | exp
func parseField(lexer *Lexer, tc *TableConstructorExp) {
	var k, v Exp

	switch lexer.LookAhead(1) {
	case TOKEN_SEP_LBRACK: // [exp]=exp
		lexer.NextToken() // TOKEN_SEP_LBRACK
		k = parseExp(lexer)
		lexer.NextTokenOfKind(TOKEN_SEP_RBRACK)
		lexer.NextTokenOfKind(TOKEN_ASSIGN)
		v = parseExp(lexer)
	case TOKEN_IDENTIFIER:
		if lexer.LookAhead(2) == TOKEN_ASSIGN { // name=exp
			line, name := lexer.NextIdentifier()
			k = &StringExp{line, name}
			lexer.NextToken() // TOKEN_ASSIGN
			v = parseExp(lexer)
		} else { // name
			tc.NArr++
			k = tc.NArr // todo
			v = parseExp(lexer)
		}
	default: // exp
		tc.NArr++
		k = tc.NArr // todo
		v = parseExp(lexer)
	}

	tc.KeyExps = append(tc.KeyExps, k)
	tc.ValExps = append(tc.ValExps, v)
}
