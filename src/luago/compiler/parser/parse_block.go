package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

// block ::= {stat} [retstat]
func parseBlock(lexer *Lexer) *Block {
	stats := parseStats(lexer)
	retStat := parseRetStat(lexer)
	lastLine := lexer.Line()
	return &Block{lastLine, stats, retStat}
}

func parseStats(lexer *Lexer) []Stat {
	stats := make([]Stat, 0, 8)
	for !isReturnOrBlockEnd(lexer.LookAhead(1)) {
		stat := parseStat(lexer)
		if _, ok := stat.(*EmptyStat); !ok {
			stats = append(stats, stat)
		}
	}
	return stats
}

func isReturnOrBlockEnd(tokenKind int) bool {
	switch tokenKind {
	case TOKEN_KW_RETURN, TOKEN_KW_END, TOKEN_EOF,
		TOKEN_KW_ELSE, TOKEN_KW_ELSEIF, TOKEN_KW_UNTIL:
		return true
	default:
		return false
	}
}

// retstat ::= return [explist] [‘;’]
// explist ::= exp {‘,’ exp}
func parseRetStat(lexer *Lexer) *RetStat {
	if lexer.LookAhead(1) != TOKEN_KW_RETURN {
		return nil
	}

	line, _ := lexer.NextTokenOfKind(TOKEN_KW_RETURN)
	switch lexer.LookAhead(1) {
	case TOKEN_KW_END, TOKEN_EOF,
		TOKEN_KW_ELSE, TOKEN_KW_ELSEIF, TOKEN_KW_UNTIL:
		return &RetStat{line, line, nil}
	case TOKEN_SEP_SEMI:
		lexer.NextToken()
		return &RetStat{line, line, nil}
	default:
		exps := parseExpList(lexer)
		if lexer.LookAhead(1) == TOKEN_SEP_SEMI {
			lexer.NextToken()
		}
		lastLine := lexer.Line()
		return &RetStat{line, lastLine, exps}
	}
}
