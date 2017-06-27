assert(math.sin(math.pi / 2) == 1.0)
assert(math.max(10.4, 7, -3, 20) == 20)
-- assert(math.huge == inf)

assert(math.floor(3.3) ==  3)
assert(math.floor(-3.3) ==  -4)
-- assert(math.floor(2^70) ==  1.1805916207174e+21)
assert(math.ceil(3.3) ==  4)
assert(math.ceil(-3.3) ==  -3)
r,f = math.modf(3.3); assert(r == 3 and tostring(f) == "0.3")
r,f = math.modf(-3.3); assert(r == -3 and tostring(f) == "-0.3")

x = 2^52 + 1
assert(string.format("%d %d", x, math.floor(x + 0.5))
    == "4503599627370497 4503599627370498")

assert(math.floor(3.5 + 0.5) == 4)
assert(math.floor(2.5 + 0.5) == 3)

print("ok")
