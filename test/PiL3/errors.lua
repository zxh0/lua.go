isOk, errMsg = pcall(pairs, 8)
assert(isOk == false
	and errMsg == "bad argument #1 to 'pairs' (table expected, got number)")

mt = {}
t = {}
setmetatable(t, mt)
mt.__metatable = "not your business"
assert(getmetatable(t) == "not your business")
isOk, errMsg = pcall(setmetatable, t, mt)
assert(isOk == false
	and errMsg == "cannot change a protected metatable")

print("ok")
