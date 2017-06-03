assert((4 and 5)       == 5)
assert((nil and 13)    == nil)
assert((false and 13)  == false)
assert((0 or 5)        == 0)
assert((false or "hi") == "hi")
assert((nil or false)  == false)

assert(not nil     == true)
assert(not false   == true)
assert(not 0       == false)
assert(not not 1   == true)
assert(not not nil == false)

print("ok")
