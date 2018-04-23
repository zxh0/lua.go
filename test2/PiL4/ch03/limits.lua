assert(math.maxinteger + 1   == math.mininteger)
assert(math.mininteger - 1   == math.maxinteger)
assert(-math.mininteger      == math.mininteger)
assert(math.mininteger // -1 == math.mininteger)

assert(math.maxinteger    ==  9223372036854775807)
assert(0x7fffffffffffffff ==  9223372036854775807)
assert(math.mininteger    == -9223372036854775808)
assert(0x8000000000000000 == -9223372036854775808)

assert(math.maxinteger + 2 == -9223372036854775807)
-- assert(tostring(math.maxinteger + 2.0) == "9.2233720368548e+18") -- TODO: FIX ME

assert(math.maxinteger + 2.0 == math.maxinteger + 1.0)

print("ok")
