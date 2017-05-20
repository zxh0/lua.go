a, b, c = 0, 1
assert(a == 0 and b == 1 and c == nil)
a, b = a+1, b+1, b+2
assert(a == 1 and b == 2)
a, b, c = 0
assert(a == 0 and b == nil and c == nil)
-- TODO: fix compiler
-- a, b, c = 0, 0, 0
-- assert(a == 0 and b == 0 and c == 0)

print("ok")
