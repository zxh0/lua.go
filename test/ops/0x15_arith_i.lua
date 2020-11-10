local a = 1
a = a + 10 -- OP_ADDI
a = a >> 1 -- OP_SHRI
a = 1 << a -- OP_SHLI
print(a)
