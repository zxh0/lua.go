package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

// prefixexp ::= var | functioncall | ‘(’ exp ‘)’
// var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
// functioncall ::=  prefixexp args | prefixexp ‘:’ Name args

/*
prefixexp ::= Name |
              ‘(’ exp ‘)’ |
              prefixexp ‘[’ exp ‘]’ |
              prefixexp ‘.’ Name |
              prefixexp ‘:’ Name args |
              prefixexp args
*/
func parsePrefixExp(lexer *Lexer) Exp {
	var exp Exp

	if lexer.LookAhead(1) == TOKEN_IDENTIFIER { // Name
		line, name := lexer.NextIdentifier()
		exp = &NameExp{line, name}
	} else { // ‘(’ exp ‘)’
		exp = parseParensExp(lexer)
	}

	for {
		switch lexer.LookAhead(1) {
		case TOKEN_SEP_LBRACK: // prefixexp ‘[’ exp ‘]’
			lexer.NextToken() // TOKEN_SEP_LBRACK
			idx := parseExp(lexer)
			lexer.NextTokenOfKind(TOKEN_SEP_RBRACK)
			exp = &TableAccessExp{lexer.Line(), exp, idx}
		case TOKEN_SEP_DOT: // prefixexp ‘.’ Name
			lexer.NextToken() // TOKEN_SEP_DOT
			line, name := lexer.NextIdentifier()
			idx := &StringExp{line, name}
			exp = &TableAccessExp{line, exp, idx}
		case TOKEN_SEP_COLON, // prefixexp ‘:’ Name args
			TOKEN_SEP_LPAREN, TOKEN_SEP_LCURLY, TOKEN_STRING: // prefixexp args
			exp = _finishFuncCallExp(lexer, exp)
		default:
			return exp
		}
	}

	return exp
}

func parseParensExp(lexer *Lexer) Exp {
	lexer.NextTokenOfKind(TOKEN_SEP_LPAREN)
	exp := parseExp(lexer)
	lexer.NextTokenOfKind(TOKEN_SEP_RPAREN)

	switch x := exp.(type) {
	case *BinopExp:
		if x.Op == TOKEN_OP_POW || x.Op == TOKEN_OP_CONCAT {
			return &ParensExp{exp}
		}
	case *VarargExp, *FuncCallExp:
		return &ParensExp{exp}
	}

	// no need to keep parens
	return exp
}

// functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
func _finishFuncCallExp(lexer *Lexer, prefixExp Exp) *FuncCallExp {
	fc := &FuncCallExp{PrefixExp: prefixExp}
	fc.NameExp = _parseNameExp(lexer)
	fc.Line = lexer.Line() // todo
	fc.Args = _parseArgs(lexer)
	fc.LastLine = lexer.Line()
	return fc
}

func _parseNameExp(lexer *Lexer) *StringExp {
	if lexer.LookAhead(1) == TOKEN_SEP_COLON {
		lexer.NextToken()
		line, name := lexer.NextIdentifier()
		return &StringExp{line, name}
	}
	return nil
}

// args ::=  ‘(’ [explist] ‘)’ | tableconstructor | LiteralString
func _parseArgs(lexer *Lexer) []Exp {
	var args []Exp = nil

	switch lexer.LookAhead(1) {
	case TOKEN_SEP_LPAREN: // ‘(’ [explist] ‘)’
		lexer.NextToken() // TOKEN_SEP_LPAREN
		if lexer.LookAhead(1) != TOKEN_SEP_RPAREN {
			args = parseExpList(lexer)
		}
		lexer.NextTokenOfKind(TOKEN_SEP_RPAREN)
	case TOKEN_SEP_LCURLY: // ‘{’ [fieldlist] ‘}’
		args = []Exp{parseTableConstructorExp(lexer)}
	default: // LiteralString
		line, str := lexer.NextTokenOfKind(TOKEN_STRING)
		args = []Exp{&StringExp{line, str}}
	}

	return args
}
