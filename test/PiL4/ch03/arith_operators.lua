assert(13 + 15     == 28)
assert(13.0 + 15.0 == 28.0)

assert(13.0 + 25  == 38.0)
assert(-(3 * 6.0) == -18.0)

assert(3.0 / 2.0 ==  1.5)
assert(3 / 2     ==  1.5)

assert(3 // 2     ==  1)
assert(3.0 // 2   ==  1.0)
assert(6 // 2     ==  3)
assert(6.0 // 2.0 ==  3.0)
assert(-9 // 2    ==  -5)
assert(1.5 // 0.5 ==  3.0)

x = math.pi
assert(x - x%0.01  ==  3.14)
assert(x - x%0.001 ==  3.141)

print("ok")
