package parser

import "testing"
import "assert"
import "luago/compiler/lexer"

func TestExpLiteral(t *testing.T) {
	testExp(t, `nil`)
	testExp(t, `true`)
	testExp(t, `false`)
	testExp(t, `123`)
	testExp(t, `'foo'`)
	testExp(t, `{}`)
	testExp(t, `...`)
}

func TestExpUnOp(t *testing.T) {
	testExp(t, `-128`)
	testExp2(t, `~a`, `~(a)`)
	testExp2(t, `a^b^c`, `(a ^ (b ^ c))`)
	testExp2(t, `#'foo'`, `#('foo')`)
	testExp2(t, `not a`, `not (a)`)
	testExp2(t, `~0xFF`, `-256`)
	testExp2(t, `not true`, `false`)
	testExp2(t, `- - - - - 1`, `-1`)
	testExp2(t, `- - - - - -1`, `1`)
	testExp2(t, `1 | -2`, `-1`)
	testExp2(t, `0xF0F ~ 0xF0`, `4095`)
	testExp2(t, `0xF & 0xFF & 0xF00`, `0`)
	testExp2(t, `4^3^2`, `262144.000000`)
	testExp2(t, `a^3^2`, `(a ^ 9.000000)`)
	testExp2(t, `4^3^a`, `(4 ^ (3 ^ a))`)
	testExp2(t, `1+2+a`, `(3 + a)`)
	testExp2(t, `a+1+2`, `((a + 1) + 2)`)
}

func TestExpBinOp(t *testing.T) {
	//testExp(t, `a + b - c * d / e // f ^ g % h & i ~ j | k`)
	//testExp(t, `k >> l << m .. n < o <= p > q >= r == s ~= t and u or v`)
	testExp2(t, `a * b + c / d`, `((a * b) + (c / d))`)
	testExp2(t, `a + b - c + d`, `(((a + b) - c) + d)`)
	testExp2(t, `a ^ b ^ c ^ d`, `(a ^ (b ^ (c ^ d)))`)
	testExp2(t, `a .. b .. c .. d`, `a .. b .. c .. d`)
	testExp2(t, `(a + b) // (c - d)`, `((a + b) // (c - d))`)
	testExp2(t, `((a ^ b) ^ c) ^ d`, `(((((a ^ b)) ^ c)) ^ d)`)
	testExp2(t, `n - 1`, `(n - 1)`)
	testExp2(t, `n-1`, `(n - 1)`)
	testExp2(t, `a or b or c`, `((a or b) or c)`)
	testExp2(t, `true or false or 2 or nil or "foo"`,
		`((true or 2) or 'foo')`)
	testExp2(t, `true and 1 and "foo" and a`, `a`)
	testExp2(t, `false and true and nil and 0 and a`,
		`((false and nil) and a)`)
	testExp2(t, `true and x and true and x and true`,
		`((x and x) and true)`)
	testExp2(t, `((((a + b))))`, `(a + b)`)
	testExp2(t, `((((a))))`, `a`)
}

func TestExpTC(t *testing.T) {
	testExp(t, `{}`)
	testExp2(t, `{...}`, `{[1]=...,}`)
	testExp2(t, `{f(),}`, `{[1]=f(),}`)
	testExp2(t, `{f(), nil}`, `{[1]=f(),[2]=nil,}`)
	testExp2(t, `{[f(1)] = g, 'x', 'y', x = 1, f(x), [30] = 23, 45}`,
		`{[f(1)]=g,[1]='x',[2]='y',['x']=1,[3]=f(x),[30]=23,[4]=45,}`)
	testExp2(t, `{ [f(1)] = g; "x", "y"; x = 1, f(x), [30] = 23; 45 }`,
		`{[f(1)]=g,[1]='x',[2]='y',['x']=1,[3]=f(x),[30]=23,[4]=45,}`)
}

func TestPrefixExp(t *testing.T) {
	testExp(t, `name`)
	testExp2(t, `(name)`, `name`)
	testExp(t, `name[key]`)
	testExp2(t, `name.field`, `name['field']`)
	testExp2(t, `a.b.c.d.e`, `a['b']['c']['d']['e']`)
	testExp(t, `a[b][c][d][e]`)
	testExp(t, `a[b[c[d[e]]]]`)
}

func TestExpFuncCall(t *testing.T) {
	testExp2(t, `print ''`, `print('')`)
	testExp2(t, `print 'hello'`, `print('hello')`)
	testExp2(t, `print {}`, `print({})`)
	testExp2(t, `print {1}`, `print({[1]=1,})`)
	testExp(t, `f()`)
	testExp(t, `g(f(), x)`)
	testExp(t, `g(x, f())`)
	testExp2(t, `io.read('*n')`, `io['read']('*n')`)
}

func TestExpPiLCh03(t *testing.T) {
	testExp2(t, `x ^ 0.5`, `(x ^ 0.500000)`)
	testExp2(t, `x ^ (-1 / 3)`, `(x ^ -0.333333)`)
	testExp2(t, `a+i < b/2+1`, `((a + i) < ((b / 2) + 1))`)
	testExp2(t, `5+x^2*8`, `(5 + ((x ^ 2) * 8))`)
	testExp2(t, `a<y and y<=z`, `((a < y) and (y <= z))`)
	testExp2(t, `x^y^z`, `(x ^ (y ^ z))`)
	testExp2(t, `-x^2`, `-((x ^ 2))`)
	testExp2(t, `-1^a`, `-((1 ^ a))`)
	// testExp(t, `-a^b + not a^b`)
}

func TestStat0(t *testing.T) {
	testStat(t, `;`)
	testStat(t, `break`)
	testStat(t, `::label::`)
	testStat(t, `goto label`)
	testStat(t, `do  end`)
	testStat2(t, `do ; end`, `do  end`)
	testStat(t, `do return end`)
	testStat2(t, `while true do ; end`, `while true do  end`)
	testStat2(t, `repeat ; until true`, `repeat  until true`)
	testStat2(t, `for v = 1, 100, 1 do ; end`, `for v = 1, 100, 1 do  end`)
	//testStat(t, `for var_1, ···, var_n in explist do block end`)
	//testStat(t, `function foo() end`)
	//testStat(t, `local function foo() end`)
	//testStat(t, `local a = 1`)
}

func TestStatIf(t *testing.T) {
	testStat2(t, `if true then ; end`, `if true then  end`)
	testStat2(t, `if a then ; else ; end`,
		`if a then  elseif true then  end`)
	testStat2(t, `if a then ; elseif b then ; else ; end`,
		`if a then  elseif b then  elseif true then  end`)
}

func TestStatFuncCall(t *testing.T) {
	testStat(t, `print()`)
	testStat(t, `print(i)`)
	testStat(t, `print('hello, world!')`)
	testStat2(t, `fact(n-1)`, `fact((n - 1))`)
	testStat2(t, `assert((4 and 5) == 5)`, `assert((5 == 5))`)
	testStat(t, `obj:f()`)
}

func TestStatAssign(t *testing.T) {
	testStat2(t, `f().a = 1`, `f()['a'] = 1`)
	testStat2(t, `a = io.read('*n')`, `a = io['read']('*n')`)
	testStat(t, `local tolerance = 10`) // todo
	testStat(t, `local f = function() end`)
}

func TestBlock(t *testing.T) {
	testBlock(t, `return`)
	testBlock(t, `return 1`)
	testBlock2(t, `return n * fact(n - 1)`, `return (n * fact((n - 1)))`)
}

func TestFuncDef(t *testing.T) {
	testBlock2(t, `function f() end`, `f = function() end`)
	testBlock2(t, `function f(a) end`, `f = function(a) end`)
	testBlock2(t, `function f(a,b,...) end`, `f = function(a, b, ...) end`)
	testBlock2(t, `function f(...) end`, `f = function(...) end`)
	testBlock2(t, `function t.a.b.c.f() end`, `t['a']['b']['c']['f'] = function() end`)
	testBlock2(t, `function t.a.b.c:f() end`, `t['a']['b']['c']['f'] = function(self) end`)
	testBlock2(t, `local function f(a) end`, `local function f(a) end`)
}

func testExp(t *testing.T, str string) {
	exp := parseExp(lexer.NewLexer("", str))
	_str := expToString(exp)
	if _str != str {
		t.Errorf(_str)
	}
}

func testExp2(t *testing.T, str, str2 string) {
	exp := parseExp(lexer.NewLexer("", str))
	_str := expToString(exp)
	if _str != str2 {
		t.Errorf(_str)
	}
}

func testStat(t *testing.T, src string) {
	stat := parseStat(lexer.NewLexer("", src))
	str := statToString(stat)
	assert.StringEqual(t, str, src)
}

func testStat2(t *testing.T, src, src2 string) {
	stat := parseStat(lexer.NewLexer("", src))
	str := statToString(stat)
	assert.StringEqual(t, str, src2)
}

func testBlock(t *testing.T, str string) {
	block := parseBlock(lexer.NewLexer("", str))
	_str := blockToString(block)
	if _str != str {
		t.Errorf(_str)
	}
}

func testBlock2(t *testing.T, str, str2 string) {
	block := parseBlock(lexer.NewLexer("", str))
	_str := blockToString(block)
	if _str != str2 {
		t.Errorf(_str)
	}
}
