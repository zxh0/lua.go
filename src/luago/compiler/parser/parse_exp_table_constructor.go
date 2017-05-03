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

	lexer.NextTokenOfKind(TOKEN_SEP_LCURLY)
	tc.Line = lexer.Line()
	if lexer.LookAhead(1) != TOKEN_SEP_RCURLY {
		parseFieldList(lexer, tc)
	}
	lexer.NextTokenOfKind(TOKEN_SEP_RCURLY)
	tc.LastLine = lexer.Line()

	return tc
}

// fieldlist ::= field {fieldsep field} [fieldsep]
func parseFieldList(lexer *Lexer, tc *TableConstructorExp) {
	// field
	parseField(lexer, tc)
	// {fieldsep field}
	for isFieldSep(lexer.LookAhead(1)) {
		lexer.NextToken()
		parseField(lexer, tc)
	}
	// [fieldsep]
	if isFieldSep(lexer.LookAhead(1)) {
		lexer.NextToken()
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
	case TOKEN_SEP_LBRACK:
		lexer.NextToken() // TOKEN_SEP_LBRACK
		k = parseExp(lexer)
		lexer.NextTokenOfKind(TOKEN_SEP_RBRACK)
		lexer.NextTokenOfKind(TOKEN_ASSIGN)
		v = parseExp(lexer)
	case TOKEN_IDENTIFIER:
		if lexer.LookAhead(2) == TOKEN_ASSIGN {
			line, name := lexer.NextIdentifier()
			k = &StringExp{line, name}
			lexer.NextToken() // TOKEN_ASSIGN
			v = parseExp(lexer)
		} else {
			tc.NArr++
			k = &IntegerExp{lexer.Line(), int64(tc.NArr)}
			v = parseExp(lexer)
		}
	default:
		tc.NArr++
		k = &IntegerExp{lexer.Line(), int64(tc.NArr)}
		v = parseExp(lexer)
	}

	tc.KeyExps = append(tc.KeyExps, k)
	tc.ValExps = append(tc.ValExps, v)
}
