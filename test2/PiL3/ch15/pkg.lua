local m = require "math"
assert(m == math)

assert(#package.searchers == 4)
assert(type(package.preload) == "table")

print("ok")
