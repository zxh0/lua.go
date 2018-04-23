assert(string.find(package.path, './?.lua', 1, true) ~= nil)
assert(package.config == "/\n;\n?\n!\n-\n")
assert(#package.searchers == 4)
assert(type(package.loaded) == "table")
assert(type(package.preload) == "table")
assert(package.loaded.math == math)
name, errMsg = package.searchpath("a.b", "x/?.lua;y/?.lua")
assert(name == nil and
	errMsg == "\n\tno file 'x/a/b.lua'\n\tno file 'y/a/b.lua'")
assert(require("math") == math)

-- require("a.b")
--[[
mymod = require('mymod')
mymod.f("hw")

local m = {}
m.f = function(a) print(a) end
return m
]]