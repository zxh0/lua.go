-- setmetatable(1)
mt = {}
debug.setmetatable(1, mt)
assert(debug.getmetatable(2) == mt)


registry = debug.getregistry()
assert(registry[2] == _G)


local uv = 100
f = function() uv = 2 end
name, val = debug.getupvalue(f, 1)
assert(name == "uv" and val == 100)
debug.setupvalue(f, 1, 200)
assert(uv == 200)


info = debug.getinfo(1)
print(info.source)
print(info.short_src)
print(info.linedefined)
print(info.lastlinedefined)
print(info.what)
print(info.currentline)
print(info.nups)
print(info.nparams)
print(info.isvararg)
print(info.name)
print(info.namewhat)
print(info.istailcall)
