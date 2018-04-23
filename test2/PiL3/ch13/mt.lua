t = {}
assert(getmetatable(t) == nil)

t1 = {}
setmetatable(t, t1)
assert(getmetatable(t) == t1)

assert(getmetatable("hi") == getmetatable("xuxu"))
assert(getmetatable(10) == nil)
assert(getmetatable(print) == nil)

print("ok")
