-- Figure 21.2. An implementation of multiple inheritance

-- look up for 'k' in list of tables 'plist'
local function search (k, plist)
  for i = 1, #plist do
    local v = plist[i][k] -- try 'i'-th superclass
    if v then return v end
  end
end

function createClass (...)
  local c = {} -- new class
  local parents = {...} -- list of parents
  
  -- class searches for absent methods in its list of parents
  setmetatable(c, {__index = function (t, k)
    return search(k, parents)
  end})
  
  -- prepare 'c' to be the metatable of its instances
  c.__index = c
  
  -- define a new constructor for this new class
  function c:new (o)
    o = o or {}
    setmetatable(o, c)
    return o
  end
  
  return c -- return new class
end
