assert(type("Hello world") == "string")
assert(type(10.4*3) == "number")
assert(type(print) == "function")
assert(type(type) == "function")
assert(type(nil) == "nil")
assert(type(type(X)) == "string")

assert(type(a) == "nil")
a = 10
assert(type(a) == "number")
a = "a string!!"
assert(type(a) == "string")
a = print
assert(type(a) == "function")

print("ok")
