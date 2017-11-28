package test

import "testing"

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
