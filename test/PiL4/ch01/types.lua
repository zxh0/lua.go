assert(type(nil)           == "nil")
assert(type(true)          == "boolean")
assert(type(10.4 * 3)      == "number")
assert(type("Hello world") == "string")
assert(type(io.stdin)      == "userdata")
assert(type(print)         == "function")
assert(type(type)          == "function")
assert(type({})            == "table")
assert(type(type(X))       == "string")

assert(type(a) == "nil")
a = 10
assert(type(a) == "number")
a = "a string!!"
assert(type(a) == "string")
a = nil
assert(type(a) == "nil")

print("ok")
