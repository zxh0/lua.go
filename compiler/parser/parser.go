package parser

import . "github.com/zxh0/lua.go/compiler/ast"
import . "github.com/zxh0/lua.go/compiler/lexer"

/* recursive descent parser */

func Parse(chunk, chunkName string) *Block {
	lexer := NewLexer(chunk, chunkName)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_EOF)
	return block
}
