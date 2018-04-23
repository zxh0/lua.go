assert(10 .. 20 == "1020")
assert("10" + 1 == 11.0)

assert(tonumber("  -3 ")    == -3)
assert(tonumber(" 10e4 ")   == 100000.0)
assert(tonumber("10e")      == nil) -- not a valid number
assert(tonumber("0x1.3p-4") == 0.07421875)

assert(tonumber("100101", 2) == 37)
assert(tonumber("fff", 16)   == 4095)
assert(tonumber("-ZZ", 36)   == -1295)
assert(tonumber("987", 8)    == nil)

assert(tostring(10) == "10")

print("ok")
