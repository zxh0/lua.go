-- Relational Operators
a = {}; a.x = 1; a.y = 0
b = {}; b.x = 1; b.y = 0
c = a
assert(a == c)
assert(a ~= b)

-- Logical Operators
assert((4 and 5)      == 5)
assert((nil and 13)   == nil)
assert((false and 13) == false)
assert((4 or 5)       == 4)
assert((false or 5)   == 5)
assert(not nil        == true)
assert(not false      == true)
assert(not 0          == false)
assert(not not 1      == true)
assert(not not nil    == false)

-- Concatenation
assert("Hello " .. "World" == "Hello World")
assert(0 .. 1 == "01")
assert(000 .. 01 == "01")
a = "Hello"
assert(a .. " World" == "Hello World")
assert(a == "Hello")

print("ok")
