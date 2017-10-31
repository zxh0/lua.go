package parser

import . "luago/compiler/ast"
import . "luago/compiler/lexer"

// block ::= {stat} [retstat]
func parseBlock(lexer *Lexer) *Block {
	return &Block{
		Stats:    parseStats(lexer),
		RetExps:  parseRetExps(lexer),
		LastLine: lexer.Line(),
	}
}

func parseStats(lexer *Lexer) []Stat {
	stats := make([]Stat, 0, 8)
	for !_isReturnOrBlockEnd(lexer.LookAhead(1)) {
		stat := parseStat(lexer)
		if _, ok := stat.(*EmptyStat); !ok {
			stats = append(stats, stat)
		}
	}
	return stats
}

func _isReturnOrBlockEnd(tokenKind int) bool {
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
func parseRetExps(lexer *Lexer) []Exp {
	if lexer.LookAhead(1) != TOKEN_KW_RETURN {
		return nil
	}

	lexer.NextTokenOfKind(TOKEN_KW_RETURN)
	switch lexer.LookAhead(1) {
	case TOKEN_KW_END, TOKEN_EOF,
		TOKEN_KW_ELSE, TOKEN_KW_ELSEIF, TOKEN_KW_UNTIL:
		return make([]Exp, 0)
	case TOKEN_SEP_SEMI:
		lexer.NextToken()
		return make([]Exp, 0)
	default:
		exps := parseExpList(lexer)
		if lexer.LookAhead(1) == TOKEN_SEP_SEMI {
			lexer.NextToken()
		}
		return exps
	}
}
