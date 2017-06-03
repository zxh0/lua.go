s, e = string.find("hello Lua users", "Lua")
assert(s == 7 and e == 9)

function foo0 () end                 -- returns no results
function foo1 () return "a" end      -- returns 1 result
function foo2 () return "a", "b" end -- returns 2 results

x,y = foo2();      assert(x == "a" and y == "b")
x = foo2();        assert(x == "a")
x,y,z = 10,foo2(); assert(x == 10 and y == "a" and z == "b")

x,y = foo0();   assert(x == nil and y == nil)
x,y = foo1();   assert(x == "a" and y == nil)
x,y,z = foo2(); assert(x == "a" and y == "b" and z == nil)

x,y = foo2(), 20;     assert(x == "a" and y == 20)
x,y = foo0(), 20, 30; assert(x == nil and y == 20)

-- print(foo0()) --> (no results)
-- print(foo1()) --> a
-- print(foo2()) --> a b
-- print(foo2(), 1) --> a 1
-- print(foo2() .. "x") --> ax (see next)

t = {foo0()} -- t = {} (an empty table)
assert(#t == 0)
t = {foo1()} -- t = {"a"}
assert(#t == 1 and t[1] == "a")
t = {foo2()} -- t = {"a", "b"}
assert(#t == 2 and t[1] == "a" and t[2] == "b")
t = {foo0(), foo2(), 4} -- t[1] = nil, t[2] = "a", t[3] = 4
assert(#t == 3 and t[1] == nil and t[2] == "a" and t[3] == 4)

function foo (i)
  if i == 0 then return foo0()
  elseif i == 1 then return foo1()
  elseif i == 2 then return foo2()
  end
end
-- print(foo(1)) --> a
-- print(foo(2)) --> a b
-- print(foo(0)) -- (no results)
-- print(foo(3)) -- (no results)

assert((foo0()) == nil) --> nil
assert((foo1()) == "a") --> a
assert((foo2()) == "a") --> a

-- table.unpack
print("ok")
