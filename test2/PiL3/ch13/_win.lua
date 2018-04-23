-- create the prototype with default values
prototype = {x = 0, y = 0, width = 100, height = 100}
mt = {} -- create a metatable
-- declare the constructor function
function new(o)
  setmetatable(o, mt)
  return o
end

mt.__index = function(_, key)
  return prototype[key]
end
w = new{x=10, y=20}
assert(w.width == 100)

mt.__index = prototype
w = new{x=10, y=20}
assert(w.width == 100)

print("ok")
