-- http://the-ravi-programming-language.readthedocs.io/en/latest/lua_bytecode_reference.html

--[[ OP_TAILCALL
function y(...) print(...) end
function z1() y(x()) end
x = function() y() end
x = function() z(1,2,3) end
x = function() local p,q,r,s = z(y()) end
x = function() print(string.char(64)) end
--]]

--[[ OP_TAILCALL
function y2() return x('foo', 'bar') end
--]]

--[[ OP_RETURN
function x2(...) return ... end
function x3(...) return 1,... end
--]]

--[[ OP_JMP
function x4() local m, n; return m >= n end
--]]

--[[ OP_VARARG
local a,b,c = ...
local a = function(...) local a,b,c = ... end
local a; a(...)
local a = {...}
local a1 = {1, ...}
local a2 = {1, f()}
do return ... end
--]]

--[[ OP_LOADBOOL
local a,b = true,false
local a = 5 > 2
--]]

--[[ OP_EQ, OP_LT and OP_LE
local x,y; if x ~= y then return "foo" else return "bar" end
if 8 > 9 then return 8 elseif 5 >= 4 then return 5 else return 9 end
local x,y; do return x ~= y end
--]]

--[[ OP_TEST and OP_TESTSET
local a,b,c; c = a and b
local a,b;   a = a and b
local a,b,c; c = a or  b
local a,b;   a = a or  b
local a,b,c; if a > b and a > c then return a end
if Done then return end
if Found and Match then return end
-- local a,b,c; a = a and b or c
--]]

--[[ OP_FORPREP and OP_FORLOOP
local a = 0; for i = 1,100,5 do a = a + i end
-- for i = 10,1,-1 do if i == 5 then break end end
--]]

-- [[ OP_TFORCALL and OP_TFORLOOP
-- for i,v in pairs(t) do print(i,v) end
--]]

--[[ OP_CLOSURE
function x() end; function y() end
local u,v; function p() return v end
local u,v; function p() u=1; local function q() return v end end
local v; local function q() return function() return v end end; do return q(), q() end
--]]

--[[ OP_GETUPVAL and OP_SETUPVAL
local a; function b() a = 1 return a end
--]]

--[[ OP_NEWTABLE
local q = {}
--]]

--[[ OP_SETLIST
local q = {1,2,3,4,5,}
local q = {1,2,3,4,5,6,7,8,9,0,
           1,2,3,4,5,6,7,8,9,0,
           1,2,3,4,5,6,7,8,9,0,
           1,2,3,4,5,6,7,8,9,0,
           1,2,3,4,5,6,7,8,9,0,
           1,2,3,4,5,6}
local q = {a=1,b=2,c=3,d=4,e=5,f=6,g=7,h=8,}
do return {1,2,3,a=1,b=2,c=3,foo()} end
local a; do return {a(), a(), a()} end
--]]

--[[ OP_GETTABLE and OP_SETTABLE 
local p = {}; p[1] = "foo"; return p["bar"]
--]]

--[[ OP_SELF
foo:bar("baz")
foo.bar(foo, "baz")
--]]

--[[ OP_GETTABUP and OP_SETTABUP
u = 40; local b = u
--]]

--[[ OP_CONCAT
local x,y = "foo","bar"; do return x..y..x..y end
local a = "foo".."bar".."baz"
--]]

--[[ OP_LEN
local a,b; a = #b; --a= #"foo"
--]]

--[[ OP_MOVE
local a,b = 10; b = a
--]]

--[[ OP_LOADNIL
-- local a,b,c,d,e = nil,nil,0
--]]

--[[ OP_LOADK
local a,b,c,d = 3,"foo",3,"foo"
--]]

--[[ Binary operators
-- local a,b = 2,4; a = a + 4 * b - a / 2 ^ b % 3
-- local a = 4 + 7 + b; a = b + 4 * 7; a = b + 4 + 7
-- local a = b + (4 + 7)
-- local a = 1 / 0; local b = 1 + "1"
--]]

--[[ Unary operators
-- local p,q = 10,false; q,p = -p,not q
-- local a = - (7 / 4)
--]]
