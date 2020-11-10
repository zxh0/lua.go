local t = {}
t[1] = 123  -- OP_SETI
print(t[1]) -- OP_GETI
t.k = 456   -- OP_SETFIELD
print(t.k)  -- OP_GETFIELD
