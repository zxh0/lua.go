print(4)
print(0.4)
assert(4.57e-3 == 0.00457)
assert(0.3e12  == 300000000000.0)
assert(5E+20   == 5e+20)

assert(type(3)   == "number")
assert(type(3.5) == "number")
assert(type(3.0) == "number")

print(3)    -- 3
print(1000) -- 100
-- TODO
-- print(3.0)  -- 3.0
-- print(1e3)  -- 1000.0

assert(1     == 1.0)
assert(-3    == -3.0)
assert(0.2e3 == 200)

assert(math.type(3)   == "integer")
assert(math.type(3.0) == "float")

assert(0xff    == 255)
assert(0x1A3   == 419)
assert(0x0.2   == 0.125)
assert(0x1p-1  == 0.5)
assert(0xa.bp2 == 42.75)

-- TODO
-- assert(string.format("%a", 419) == "0x1.a3p+8")
-- assert(string.format("%a", 0.1) == "0x1.999999999999ap-4")

print("ok")
