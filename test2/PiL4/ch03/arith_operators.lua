a = 13; b = 15; c = a + b -- 13 + 15 --> 28
assert(c == 28 and math.type(c) == "integer")
a = 13.0; b = 15.0; c = a + b -- 13.0 + 15.0 --> 28.0
assert(c == 28.0 and math.type(c) == "float")

a = 13.0; b = 25; c = a + b -- 13.0 + 25 --> 38.0
assert(c == 38.0 and math.type(c) == "float")
a = 3; b = 6.0; c = -(a * b) -- -(3 * 6.0) --> -18.0
assert(c == -18.0 and math.type(c) == "float")

a = 3.0; b = 2.0; c = a / b -- 3.0 / 2.0 --> 1.5
assert(c == 1.5 and math.type(c) == "float")
a = 3; b = 2; c = a / b -- 3 / 2 --> 1.5
assert(c == 1.5 and math.type(c) == "float")

a = 3; b = 2; c = a // b -- 3 // 2 --> 1
assert(c == 1 and math.type(c) == "integer")
a = 3.0; b = 2; c = a // b -- 3.0 // 2 --> 1.0
assert(c == 1.0 and math.type(c) == "float")
a = 6; b = 2; c = a // b -- 6 // 2 --> 3
assert(c == 3 and math.type(c) == "integer")
a = 6.0; b = 2.0; c = a // b -- 6.0 // 2.0 --> 3.0
assert(c == 3.0 and math.type(c) == "float")
a = -9; b = 2; c = a // b -- -9 // 2 --> -5
assert(c == -5 and math.type(c) == "integer")
a = 1.5; b = 0.5; c = a // b -- 1.5 // 0.5 --> 3.0
assert(c == 3.0 and math.type(c) == "float")

x = math.pi
print(x - x%0.01)  --> 3.14
print(x - x%0.001) --> 3.141

function idiv(a, b) return a // b end
assert(idiv( 7,    3) ==  2  )
assert(idiv( 7,   -3) == -3  )
assert(idiv(-7,    3) == -3  )
assert(idiv(-7,   -3) ==  2  )
assert(idiv( 7.0,  3) ==  2.0)
assert(idiv( 7.0, -3) == -3.0)
assert(idiv(-7.0,  3) == -3.0)
assert(idiv(-7.0, -3) ==  2.0)

function mod(a, b) return a % b end
assert(mod( 2,  5) ==  2)
assert(mod( 5,  2) ==  1)
assert(mod( 5, -2) == -1)
assert(mod(-5,  2) ==  1)
assert(mod(-5, -1) ==  0)
assert(mod(-1.5, 2.0) == 0.5)

print("ok")
