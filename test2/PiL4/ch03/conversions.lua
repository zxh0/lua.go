a = -3; b = a + 0.0 -- -3 + 0.0 --> -3.0
assert(b == -3.0 and math.type(b) == "float")
a = 0x7fffffffffffffff; b = a + 0.0 -- 0x7fffffffffffffff + 0.0 --> 9.2233720368548e+18
-- assert(a + 0.0 == 9.2233720368548e+18) -- TODO: FIX ME

assert(9007199254740991 + 0.0 == 9007199254740991)
assert(9007199254740992 + 0.0 == 9007199254740992)
-- assert(9007199254740993 + 0.0 ~= 9007199254740993) -- TODO

-- assert(tostring(2^53) == "9.007199254741e+15") -- (float) -- TODO
assert(2^53 | 0 == 9007199254740992)           -- (integer)

-- TODO: FIX ME
-- ok,err = pcall(function() return 3.2 | 0 end) -- fractional part
-- assert(not ok and err:find("number has no integer representation") > 0)
-- ok,err = pcall(function() return 2^64 | 0 end) -- out of range
-- assert(not ok and err:find("number has no integer representation") > 0)
-- ok,err = pcall(math.random, 1, 3.5)
-- assert(not ok and 
--   err == "bad argument #2 to 'math.random' (number has no integer representation)")

assert(math.tointeger(-258.0) == -258)
assert(math.tointeger(2^30) == 1073741824)
assert(math.tointeger(5.01) == nil) -- (not an integral value)
assert(math.tointeger(2^64) == nil) -- (out of range)

function cond2int (x)
  return math.tointeger(x) or x
end

print("ok")
