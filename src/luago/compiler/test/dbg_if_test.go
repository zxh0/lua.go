package test

import "testing"

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
