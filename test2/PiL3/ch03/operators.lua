-- Relational Operators
a = {}; a.x = 1; a.y = 0
b = {}; b.x = 1; b.y = 0
c = a
assert(a == c)
assert(a ~= b)

print("ok")
