package compiler

import "luago/binchunk"
import "luago/compiler/ast"
import "luago/compiler/codegen"
import "luago/compiler/parser"

func Compile(source, chunk string) *binchunk.FuncProto {
	block := parser.Parse(source, chunk)
	funcBody := &ast.FuncDefExp{
		LastLine: block.LastLine,
		IsVararg: true,
		Block:    block,
	}
	proto := codegen.GenProto(funcBody)
	setSource(proto, source) // todo
	return proto
}

// todo
func setSource(proto *binchunk.FuncProto, source string) {
	proto.Source = "@" + source
	for _, f := range proto.Protos {
		setSource(f, source)
	}
}
