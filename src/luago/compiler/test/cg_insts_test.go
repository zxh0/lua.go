package test

import "fmt"
import "strings"
import "testing"
import "assert"
import "luago/compiler"

func TestLocVar(t *testing.T) {
	testInsts(t, "local a=1; local a=1", "[2/2] loadk(0,-1); loadk(1,-1)")
	testInsts(t, "do local a=1 end; do local a=1 end; local a=1", "[2/3] loadk(0,-1); loadk(0,-1); loadk(0,-1)")
	testInsts(t, "local a=1; do end; local a=1", "[2/2] loadk(0,-1); loadk(1,-1)")
	testInsts(t, "if x then local a else local a end",
		"[2/2] gettabup(0,0,-1); test(0,_,0); jmp(0,2); loadnil(0,0,_); jmp(0,4); loadbool(0,1,0); test(0,_,0); jmp(0,1); loadnil(0,0,_)")
}

func TestReturn(t *testing.T) {
	testInsts(t, "return", "[2/0] return(0,1,_)")
	testInsts(t, "local a,b,c; return a", "[3/3] loadnil(0,2,_); return(0,2,_)")
	testInsts(t, "local a,b,c; return b,c", "[5/3] loadnil(0,2,_); move(3,1,_); move(4,2,_); return(3,3,_)")
	testInsts(t, "return f(),1", "[2/0] gettabup(0,0,-1); call(0,1,2); loadk(1,-2); return(0,3,_)")
	testInsts(t, "return 1,f()", "[2/0] loadk(0,-1); gettabup(1,0,-2); call(1,1,0); return(0,0,_)")
}

func TestConcat(t *testing.T) {
	testInsts(t, "x = a .. b .. c; local d",
		`[4/1]
gettabup(1,0,-2);
gettabup(2,0,-3);
gettabup(3,0,-4);
concat(0,1,3);
settabup(0,-1,0);
loadnil(0,0,_)`)
}

func TestCMP(t *testing.T) {
	testInsts(t, "local a,b; c = a==b", "[3/2] loadnil(0,1,_); eq(1,0,1); jmp(0,1); loadbool(2,0,1); loadbool(2,1,0); settabup(0,-1,2)")
	testInsts(t, "local a,b; c = a<=b", "[3/2] loadnil(0,1,_); le(1,0,1); jmp(0,1); loadbool(2,0,1); loadbool(2,1,0); settabup(0,-1,2)")
	testInsts(t, "local a,b; c = a< b", "[3/2] loadnil(0,1,_); lt(1,0,1); jmp(0,1); loadbool(2,0,1); loadbool(2,1,0); settabup(0,-1,2)")
	testInsts(t, "local a,b; c = a>=b", "[3/2] loadnil(0,1,_); le(1,1,0); jmp(0,1); loadbool(2,0,1); loadbool(2,1,0); settabup(0,-1,2)")
	testInsts(t, "local a,b; c = a> b", "[3/2] loadnil(0,1,_); lt(1,1,0); jmp(0,1); loadbool(2,0,1); loadbool(2,1,0); settabup(0,-1,2)")
}

func TestOP(t *testing.T) {
	testInsts(t, "local a,b; c=-b", "[3/2] loadnil(0,1,_); unm(2,1,_); settabup(0,-1,2)")
	testInsts(t, "local a,b; c=a+b", "[3/2] loadnil(0,1,_); add(2,0,1); settabup(0,-1,2)")
	testInsts(t, "local a,b; c=a+1", "[3/2] loadnil(0,1,_); add(2,0,-2); settabup(0,-1,2)")
	testInsts(t, "x=a+b", "[3/0] gettabup(1,0,-2); gettabup(2,0,-3); add(0,1,2); settabup(0,-1,0)")
	testInsts(t, "x=a+b+c", "[4/0] gettabup(2,0,-2); gettabup(3,0,-3); add(1,2,3); gettabup(2,0,-4); add(0,1,2); settabup(0,-1,0)")
	testInsts(t, "local a,b; c = a and b", "[3/2] loadnil(0,1,_); testset(2,0,0); jmp(0,1); move(2,1,_); settabup(0,-1,2)")
	testInsts(t, "x = a and b and c",
		`[3/0]
gettabup(2,0,-2); testset(1,2,0); jmp(0,2);
gettabup(2,0,-3); move(1,2,_); testset(0,1,0); jmp(0,2);
gettabup(1,0,-4); move(0,1,_);
settabup(0,-1,0)`)
	testInsts(t, "x = a or b or c",
		`[3/0]
gettabup(2,0,-2); testset(1,2,1); jmp(0,2);
gettabup(2,0,-3); move(1,2,_); testset(0,1,1); jmp(0,2);
gettabup(1,0,-4); move(0,1,_);
settabup(0,-1,0)`)
}

func TestTcExp(t *testing.T) {
	testInsts(t, "local a={1,2}",
		`[3/1]
newtable(0,2,0);
loadk(1,-1);
loadk(2,-2);
setlist(0,2,1)`)
	testInsts(t, "local a={1,f()}",
		`[3/1]
newtable(0,2,0);
loadk(1,-1);
gettabup(2,0,-2);
call(2,1,0);
setlist(0,0,1)`)
	testInsts(t, "local a={1,f(),2}",
		`[4/1]
newtable(0,3,0);
loadk(1,-1);
gettabup(2,0,-2);
call(2,1,2);
loadk(3,-3);
setlist(0,3,1)`)
}

func TestLocalFuncDefStat(t *testing.T) {
	testInsts(t, "local function f() g() end; local a", "[2/2] closure(0,0); loadnil(1,0,_)")
	//testInsts(t, "local function f() return ... end", "!") // cannot use '...' outside a vararg function near '...'
}

func TestFuncCallStat(t *testing.T) {
	testInsts(t, "local f; f()", "[2/1] loadnil(0,0,_); move(1,0,_); call(1,1,1)")
	testInsts(t, "f()", "[2/0] gettabup(0,0,-1); call(0,1,1)")
	testInsts(t, "f(1,2)", "[3/0] gettabup(0,0,-1); loadk(1,-2); loadk(2,-3); call(0,3,1)")
	testInsts(t, "f(1,g(2,h(3)))",
		`[6/0]
gettabup(0,0,-1); loadk(1,-2);
gettabup(2,0,-3); loadk(3,-4);
gettabup(4,0,-5); loadk(5,-6);
call(4,2,0); call(2,0,0); call(0,0,1)`)
	testInsts(t, "obj:f()", "[2/0] gettabup(0,0,-1); self(0,0,-2); call(0,2,1)")
	testInsts(t, "obj:f(...)", "[3/0] gettabup(0,0,-1); self(0,0,-2); vararg(2,0,_); call(0,0,1)")
	//testInsts(t, "local a,b,c; a:f()", "[5/3] loadnil(0,2,_); move(3,0,_); self(3,3,-1); call(3,2,1)")
	//testInsts(t, "a.b.c:f()", "[3/0] gettabup(2,0,-1); gettable(1,2,-2); gettable(0,1,-3); self(0,0,-4); call(0,2,1)")
}

func TestRepeatStat(t *testing.T) {
	testInsts(t, "repeat f() until g()",
		`[2/0]
gettabup(0,0,-1); call(0,1,1);
gettabup(0,0,-2); call(0,1,2);
test(0,_,0); jmp(0,-6)`)
}

func TestWhileStat(t *testing.T) {
	testInsts(t, "while f() do g() end",
		`[2/0]
gettabup(0,0,-1); call(0,1,2);
test(0,_,0); jmp(0,3);
gettabup(0,0,-2); call(0,1,1);
jmp(0,-7)`)
}

func TestIfStat(t *testing.T) {
	testInsts(t, "if a then f() elseif b then g() end",
		`[2/0]
gettabup(0,0,-1); test(0,_,0); jmp(0,3);
gettabup(0,0,-2); call(0,1,1); jmp(0,5);
gettabup(0,0,-3); test(0,_,0); jmp(0,2);
gettabup(0,0,-4); call(0,1,1)`)
}

func TestForNumStat(t *testing.T) {
	testInsts(t, "for i=1,100,2 do f() end; local a",
		`[5/5]
loadk(0,-1);
loadk(1,-2);
loadk(2,-3);
forprep(0,2);
gettabup(4,0,-4);
call(4,1,1);
forloop(0,-3);
loadnil(0,0,_)`)
}

func TestForInStat(t *testing.T) {
	testInsts(t, "for k,v in pairs(t) do print(k,v) end; local a",
		`[8/6]
gettabup(0,0,-1);
gettabup(1,0,-2);
call(0,2,4);
jmp(0,4);
gettabup(5,0,-3);
move(6,3,_);
move(7,4,_);
call(5,3,1);
tforcall(0,_,2);
tforloop(2,-6);
loadnil(0,0,_)`)
}

func TestBreakStat(t *testing.T) {
	//testInsts(t, "do break end", "")
	testInsts(t,
		`
while x do 
  break; 
  do 
    break 
    do
      break
    end
  end 
end
`,
		`[2/0]
gettabup(0,0,-1);
test(0,_,0);
jmp(0,4);
jmp(0,3);
jmp(0,2);
jmp(0,1);
jmp(0,-7)`)
}

func TestLocalAssignStat(t *testing.T) {
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
	testInsts(t, "local a=t[100]", "[2/1] gettabup(1,0,-1); gettable(0,1,-2)")
	testInsts(t, "t={}; local a=t[100]", "[2/1] newtable(0,0,0); settabup(0,-1,0); gettabup(1,0,-1); gettable(0,1,-2)")
}

func TestAssignStat(t *testing.T) {
	testInsts(t, "local a; a=nil", "[2/1] loadnil(0,0,_); loadnil(1,0,_); move(0,1,_)")
	testInsts(t, "local a; a=1", "[2/1] loadnil(0,0,_); loadk(1,-1); move(0,1,_)")
	testInsts(t, "local a; a=f()", "[2/1] loadnil(0,0,_); gettabup(1,0,-1); call(1,1,2); move(0,1,_)")
	testInsts(t, "local a; a=1,f()", "[3/1] loadnil(0,0,_); loadk(1,-1); gettabup(2,0,-2); call(2,1,1); move(0,1,_)")
	testInsts(t, "local a; a=f(),1", "[3/1] loadnil(0,0,_); gettabup(1,0,-1); call(1,1,2); loadk(2,-2); move(0,1,_)")
	testInsts(t, "local a; a[1]=2", "[4/1] loadnil(0,0,_); move(1,0,_); loadk(2,-1); loadk(3,-2); settable(1,2,3)")
	testInsts(t, "a=nil", "[2/0] loadnil(0,0,_); settabup(0,-1,0)")
	testInsts(t, "a=1", "[2/0] loadk(0,-2); settabup(0,-1,0)")
	//testInsts(t, "local a; a=a+1", "")
}

func testInsts(t *testing.T, chunk, expected string) {
	insts := compile(chunk, false)
	expected = strings.Replace(expected, "\n", " ", -1)
	expected += "; return(0,1,_);"
	assert.StringEqual(t, insts, expected)
}

func testDbg(t *testing.T, chunk, expected string) {
	insts := compile(chunk, true)
	expected = strings.Replace(expected, "\n", " ", -1)
	assert.StringEqual(t, insts, expected)
}

func compile(chunk string, dbgFlag bool) string {
	proto := compiler.Compile("src", chunk)

	s := fmt.Sprintf("[%d/%d] ", proto.MaxStackSize, len(proto.LocVars))
	for i, inst := range proto.Code {
		if dbgFlag {
			s += fmt.Sprintf("[%2d]", proto.LineInfo[i])
		}
		s += instToStr(inst)
		s += "; "
	}
	if dbgFlag {
		for _, locVar := range proto.LocVars {
			s += fmt.Sprintf("@%s[%d,%d] ",
				locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
		}
	}

	return strings.TrimSpace(s)
}
