local a = true
local b = false
local c = a and b -- OP_TESTSET
local d = a or b  -- OP_TESTSET
local e = not a   -- OP_NOT
a = a or b        -- OP_TEST
print(c, d, e, a)
