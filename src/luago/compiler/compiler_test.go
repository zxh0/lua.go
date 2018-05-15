package compiler

import "fmt"
import "strings"
import "testing"
import "github.com/stretchr/testify/assert"
import . "luago/vm"

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

func TestLocalVarDeclStat(t *testing.T) {
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

func TestIfDbg(t *testing.T) {
	testDbg(t,
		`
local a,b,c
if
a
then -- 5
local 
x
f
(
) -- 10
elseif 
x
then 
local 
x -- 15
b
(
)
elseif 
g -- 20
(
)
then 
x
( -- 25
) 
end
`,
		`[5/5]
[ 2]loadnil(0,2,_);
[ 4]test(0,_,0);
[ 4]jmp(0,4);
[ 7]loadnil(3,0,_);
[ 8]gettabup(4,0,-1);
[ 8]call(4,1,1);
[10]jmp(0,13);
[12]gettabup(3,0,-2);
[12]test(3,_,0);
[12]jmp(0,4);
[15]loadnil(3,0,_);
[16]move(4,1,_);
[16]call(4,1,1);
[18]jmp(0,6);
[20]gettabup(3,0,-3);
[20]call(3,1,2);
[22]test(3,_,0);
[22]jmp(0,2);
[24]gettabup(3,0,-2);
[24]call(3,1,1);
[27]return(0,1,_);
@a[2,22]
@b[2,22]
@c[2,22]
@x[5,7]
@x[12,14]`)
}

func TestRepeatDbg(t *testing.T) {
	testDbg(t,
		`
repeat
local 
a
, -- 5
b;
c
=
f
( -- 10
a
,
b
)
until -- 15
a
`,
		`[5/2]
[ 6]loadnil(0,1,_);
[ 9]gettabup(2,0,-2);
[11]move(3,0,_);
[13]move(4,1,_);
[ 9]call(2,3,2);
[14]settabup(0,-1,2);
[16]test(0,_,0);
[16]jmp(0,-8);
[16]return(0,1,_);
@a[2,9]
@b[2,9]`)
}

func TestWhileDbg(t *testing.T) {
	testDbg(t,
		`
while
a
do
local -- 5
b
,
c;
f
( -- 10
a
,
b
,
c -- 15
)
end
`,
		`[6/2]
[ 3]gettabup(0,0,-1);
[ 3]test(0,_,0);
[ 3]jmp(0,7);
[ 8]loadnil(0,1,_);
[ 9]gettabup(2,0,-2);
[11]gettabup(3,0,-1);
[13]move(4,0,_);
[15]move(5,1,_);
[ 9]call(2,4,1);
[16]jmp(0,-10);
[17]return(0,1,_);
@b[5,10]
@c[5,10]`)
}

func TestForNumDbg(t *testing.T) {
	testDbg(t,
		`
for
i
=
1 -- 5
,
100
,
2
do -- 10
f
(
i
)
end -- 15
`,
		`[6/4]
[ 5]loadk(0,-1);
[ 7]loadk(1,-2);
[ 9]loadk(2,-3);
[10]forprep(0,3);
[11]gettabup(4,0,-4);
[13]move(5,3,_);
[11]call(4,2,1);
[ 2]forloop(0,-4);
[15]return(0,1,_);
@(for index)[4,9]
@(for limit)[4,9]
@(for step)[4,9]
@i[5,8]`)
}

func TestForInDbg(t *testing.T) {
	testDbg(t,
		`
for
k
,
v -- 5
in
pairs
(
t
) -- 10
do
f
(
k
, -- 15
v
)
end
`,
		`[8/5]
[ 7]gettabup(0,0,-1);
[ 9]gettabup(1,0,-2);
[ 7]call(0,2,4);
[11]jmp(0,4);
[12]gettabup(5,0,-3);
[14]move(6,3,_);
[16]move(7,4,_);
[12]call(5,3,1);
[ 7]tforcall(0,_,2);
[ 7]tforloop(2,-6);
[18]return(0,1,_);
@(for generator)[4,11]
@(for state)[4,11]
@(for control)[4,11]
@k[5,9]
@v[5,9]`)
}

func TestK(t *testing.T) {
	chunk := `
t = {0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,
20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,
40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,
60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,
80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,
100,101,102,103,104,105,106,107,108,109,110,111,112,113,114,115,116,117,118,119,
120,121,122,123,124,125,126,127,128,129,130,131,132,133,134,135,136,137,138,139,
140,141,142,143,144,145,146,147,148,149,150,151,152,153,154,155,156,157,158,159,
160,161,162,163,164,165,166,167,168,169,170,171,172,173,174,175,176,177,178,179,
180,181,182,183,184,185,186,187,188,189,190,191,192,193,194,195,196,197,198,199,
200,201,202,203,204,205,206,207,208,209,210,211,212,213,214,215,216,217,218,219,
220,221,222,223,224,225,226,227,228,229,230,231,232,233,234,235,236,237,238,239,
240,241,242,243,244,245,246,247,248,249,250,251,252,253,254,255}
print(t[100], t[200], t[255])
`
	print(chunk)
	//testInsts(t, chunk, ``)
}

func TestDoUV(t *testing.T) {
	testDbg(t,
		`
do
  local a,b,c
  function f() print(b) end
end
`,
		`[4/3]
[ 3]loadnil(0,2,_);
[ 4]closure(3,0);
[ 4]settabup(0,-1,3);
[ 4]jmp(1,0);
[ 5]return(0,1,_);
@a[2,5]
@b[2,5]
@c[2,5]`)
}

func TestForNumUV(t *testing.T) {
	testDbg(t,
		`
local a,b,c
for i=1,100 do
  local function f() print(i) end
  break
  do do break end end
end
`,
		`[8/8]
[ 2]loadnil(0,2,_);
[ 3]loadk(3,-1);
[ 3]loadk(4,-2);
[ 3]loadk(5,-1);
[ 3]forprep(3,4);
[ 4]closure(7,0);
[ 5]jmp(7,3);
[ 6]jmp(7,2);
[ 6]jmp(7,0);
[ 3]forloop(3,-5);
[ 7]return(0,1,_);
@a[2,12]
@b[2,12]
@c[2,12]
@(for index)[5,11]
@(for limit)[5,11]
@(for step)[5,11]
@i[6,10]
@f[7,10]`)
}

func TestForInUV(t *testing.T) {
	testDbg(t,
		`
local a,b,c
for k,v in pairs(t) do
  local function f() print(k,v) end
end
`,
		`[9/9]
[ 2]loadnil(0,2,_);
[ 3]gettabup(3,0,-1);
[ 3]gettabup(4,0,-2);
[ 3]call(3,2,4);
[ 3]jmp(0,2);
[ 4]closure(8,0);
[ 4]jmp(7,0);
[ 3]tforcall(3,_,2);
[ 3]tforloop(5,-4);
[ 5]return(0,1,_);
@a[2,11]
@b[2,11]
@c[2,11]
@(for generator)[5,10]
@(for state)[5,10]
@(for control)[5,10]
@k[6,8]
@v[6,8]
@f[7,8]`)
}

func TestWhileUV(t *testing.T) {
	testDbg(t,
		`
while x do
  local a,b,c
  local function f() print(b) end
end
`,
		`[4/4]
[ 2]gettabup(0,0,-1);
[ 2]test(0,_,0);
[ 2]jmp(0,4);
[ 3]loadnil(0,2,_);
[ 4]closure(3,0);
[ 4]jmp(1,0);
[ 4]jmp(0,-7);
[ 5]return(0,1,_);
@a[5,7]
@b[5,7]
@c[5,7]
@f[6,7]`)
}

func TestRepeatUV(t *testing.T) {
	testDbg(t,
		`
repeat
  local a,b,c
  local function f() print(b) end
until x
`,
		`[5/4]
[ 3]loadnil(0,2,_);
[ 4]closure(3,0);
[ 5]gettabup(4,0,-1);
[ 5]test(4,_,0);
[ 5]jmp(1,-5);
[ 5]jmp(1,0);
[ 5]return(0,1,_);
@a[2,7]
@b[2,7]
@c[2,7]
@f[3,7]`)
}

func TestIfUV(t *testing.T) {
	testDbg(t,
		`
if x then
  local a,b,c
  local function f() print(b) end
elseif y then 
  local a,b,c
  local function f() print(b) end
else
  local a,b,c
  local function f() print(b) end
end
`,
		`[4/12]
[ 2]gettabup(0,0,-1);
[ 2]test(0,_,0);
[ 2]jmp(0,4);
[ 3]loadnil(0,2,_);
[ 4]closure(3,0);
[ 4]jmp(1,0);
[ 4]jmp(0,13);
[ 5]gettabup(0,0,-2);
[ 5]test(0,_,0);
[ 5]jmp(0,4);
[ 6]loadnil(0,2,_);
[ 7]closure(3,1);
[ 7]jmp(1,0);
[ 7]jmp(0,6);
[ 8]loadbool(0,1,0);
[ 8]test(0,_,0);
[ 8]jmp(0,3);
[ 9]loadnil(0,2,_);
[10]closure(3,2);
[10]jmp(1,0);
[11]return(0,1,_);
@a[5,7]
@b[5,7]
@c[5,7]
@f[6,7]
@a[12,14]
@b[12,14]
@c[12,14]
@f[13,14]
@a[19,21]
@b[19,21]
@c[19,21]
@f[20,21]`)
}

func TestGoto(t *testing.T) {
	testDbg(t,
		`
local a=1
::L1::
local b=2
goto L1`,
`[2/2]
[ 2]loadk(0,-1);
[ 4]loadk(1,-2);
[ 5]jmp(2,-2);
[ 5]return(0,1,_);
@a[2,5]
@b[3,5]`)
}

func testInsts(t *testing.T, chunk, expected string) {
	insts := compile(chunk, false)
	expected = strings.Replace(expected, "\n", " ", -1)
	expected += "; return(0,1,_);"
	assert.Equal(t, expected, insts)
}

func testDbg(t *testing.T, chunk, expected string) {
	insts := compile(chunk, true)
	expected = strings.Replace(expected, "\n", " ", -1)
	assert.Equal(t, expected, insts)
}

func compile(chunk string, dbgFlag bool) string {
	proto := Compile("src", chunk)

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



func instToStr(_i uint32) string {
	i := Instruction(_i)
	opName := strings.ToLower(i.OpName())
	opName = strings.TrimSpace(opName)

	switch i.OpMode() {
	case IABC:
		a, b, c := i.ABC()
		return fmt.Sprintf("%s(%d,%s,%s)", opName, a,
			argBCToStr(b, i.BMode()),
			argBCToStr(c, i.CMode()))
	case IABx:
		a, bx := i.ABx()
		return fmt.Sprintf("%s(%d,%s)", opName, a,
			argBxToStr(bx, i.BMode()))
	case IAsBx:
		a, sbx := i.AsBx()
		return fmt.Sprintf("%s(%d,%d)", opName, a, sbx)
	case IAx:
		ax := i.Ax()
		return fmt.Sprintf("%s(%d)", opName, -1-ax)
	default:
		panic("unreachable!")
	}
}

func argBCToStr(arg int, mode byte) string {
	if mode == OpArgN {
		return "_"
	}
	if arg > 0xFF {
		return fmt.Sprintf("%d", -1-arg&0xFF)
	}
	return fmt.Sprintf("%d", arg)
}

func argBxToStr(bx int, mode byte) string {
	if mode == OpArgK {
		return fmt.Sprintf("%d", -1-bx)
	}
	if mode == OpArgU {
		return fmt.Sprintf("%d", bx)
	}
	return "_"
}
