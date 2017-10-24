package test

import "fmt"
import "testing"
import "assert"
import "luago/compiler"

func TestLocalAssign(t *testing.T) {
	testInsts(t, "local a", "[2/1] loadnil(0,0,_)")
	testInsts(t, "local a=nil", "[2/1] loadnil(0,0,_)")
	testInsts(t, "local a=true", "[2/1] loadbool(0,1,0)")
	testInsts(t, "local a=false", "[2/1] loadbool(0,0,0)")
	testInsts(t, "local a=1", "[2/1] loadk(0,-1)")
	testInsts(t, "local a='foo'", "[2/1] loadk(0,-1)")
	testInsts(t, "local a,b,c=1,2,3", "[3/3] loadk(0,-1); loadk(1,-2); loadk(2,-3)")
	testInsts(t, "local a,b,c=f()", "[3/3] gettabup(0,0,-1); call(0,1,4)")
	testInsts(t, "local a=1,nil", "[2/1] loadk(0,-1)")
	testInsts(t, "local a=1,f()", "[2/1] loadk(0,-1); gettabup(1,0,-2); call(1,1,1)")
	testInsts(t, "local a,b,c", "[3/3] loadnil(0,2,_)")
}

func TestAssign(t *testing.T) {
	testInsts(t, "local a; a=nil", "[2/1] loadnil(0,0,_); loadnil(1,0,_); move(0,1,_)")
	testInsts(t, "local a; a=1", "[2/1] loadnil(0,0,_); loadk(1,-1); move(0,1,_)")
	testInsts(t, "local a; a=f()", "[2/1] loadnil(0,0,_); gettabup(1,0,-1); call(1,1,2); move(0,1,_)")
	testInsts(t, "local a; a=1,f()", "[3/1] loadnil(0,0,_); loadk(1,-1); gettabup(2,0,-2); call(2,1,1); move(0,1,_)")
	testInsts(t, "local a; a=f(),1", "[3/1] loadnil(0,0,_); gettabup(1,0,-1); call(1,1,2); loadk(2,-2); move(0,1,_)")
	testInsts(t, "local a; a[1]=2", "[4/1] loadnil(0,0,_); move(1,0,_); loadk(2,-1); loadk(3,-2); settable(1,2,3)")
	testInsts(t, "a=nil", "[2/0] loadnil(0,0,_); settabup(0,-1,0)")
	testInsts(t, "a=1", "[2/0] loadk(0,-1); settabup(0,-2,0)")
}

func testInsts(t *testing.T, chunk, expected string) {
	insts := compile(chunk)
	assert.StringEqual(t, insts, expected + "; return(0,1,_)")
}

func compile(chunk string) string {
	proto := compiler.Compile("src", chunk)

	s := fmt.Sprintf("[%d/%d] ", proto.MaxStackSize, len(proto.LocVars))
	for i, inst := range proto.Code {
		s += instToStr(inst)
		if i < len(proto.Code)-1 {
			s += "; "
		}
	}

	return s
}
