package test

import "testing"

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
