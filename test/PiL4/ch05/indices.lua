a = {}
for i = 1, 1000 do a[i] = i*2 end
assert(a[9] == 18)
a["x"] = 10
assert(a["x"] == 10)
assert(a["y"] == nil)

a = {}
a.x = 10
assert(a.x == 10)
assert(a.y == nil)

a = {}
x = "y"
a[x] = 10
assert(a[x] == 10)
assert(a.x  == nil)
assert(a.y  == 10)

i = 10; j = "10"; k = "+10"
a = {}
a[i] = "one value"
a[j] = "another value"
a[k] = "yet another value"
assert(a[i] == "one value")
assert(a[j] == "another value")
assert(a[k] == "yet another value")
assert(a[tonumber(j)] == "one value")
assert(a[tonumber(k)] == "one value")

a = {}
a[2.0] = 10
a[2.1] = 20
assert(a[2]   == 10)
assert(a[2.1] == 20)

print("ok")
