package test

import "testing"

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
