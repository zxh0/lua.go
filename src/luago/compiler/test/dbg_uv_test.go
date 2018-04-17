package test

import "testing"

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
