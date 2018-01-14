mt = {}
mt.__len = function(a, b) return 100 end

s = "foo"
-- debug.setmetatable(s, mt)
-- assert(#s == 3)

t = {}
setmetatable(t, mt)
assert(#t == 100)

mt.__eq = function(a, b) return 100 end
assert(t == t)
assert((t == {}) == true)
assert({} ~= print)

setmetatable(t, nil)