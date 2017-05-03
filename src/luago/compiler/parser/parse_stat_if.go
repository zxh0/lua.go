package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

// if exp then block {elseif exp then block} [else block] end
func parseIfStat(lexer *Lexer) *IfStat {
	stat := &IfStat{
		Lines:  make([]int, 0, 8),
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

	line, _ := lexer.NextTokenOfKind(TOKEN_KW_THEN)
	stat.Lines = append(stat.Lines, line)
	stat.Blocks = append(stat.Blocks, parseBlock(lexer))
}

// {elseif exp then block}
func _parseElseIf(lexer *Lexer, stat *IfStat) {
	for lexer.LookAhead(1) == TOKEN_KW_ELSEIF {
		lexer.NextTokenOfKind(TOKEN_KW_ELSEIF)
		stat.Exps = append(stat.Exps, parseExp(lexer))

		line, _ := lexer.NextTokenOfKind(TOKEN_KW_THEN)
		stat.Lines = append(stat.Lines, line)
		stat.Blocks = append(stat.Blocks, parseBlock(lexer))
	}
}

// [else block] end
func _parseElse(lexer *Lexer, stat *IfStat) {
	if lexer.LookAhead(1) == TOKEN_KW_ELSE {
		line, _ := lexer.NextTokenOfKind(TOKEN_KW_ELSE)
		stat.Lines = append(stat.Lines, line)

		// else block => elseif true then block
		stat.Exps = append(stat.Exps, &TrueExp{line})
		stat.Blocks = append(stat.Blocks, parseBlock(lexer))
	}

	line, _ := lexer.NextTokenOfKind(TOKEN_KW_END)
	stat.Lines = append(stat.Lines, line)
}
