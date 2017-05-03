a = {}
k = "x"
a[k] = 10
a[20] = "great"
assert(a["x"] == 10)
k = 20
assert(a[k] == "great")
a["x"] = a["x"] + 1
assert(a["x"] == 11)

a = {}
a["x"] = 10
b = a
assert(b["x"] == 10)
b["x"] = 20
assert(a["x"] == 20)
a = nil
b = nil

a = {}
for i = 1, 1000 do a[i] = i*2 end
assert(a[9] == 18)
a["x"] = 10
assert(a["x"] == 10)
assert(a["y"] == nil)

a.x = 10
assert(a.x == 10)
assert(a.y == nil)

a = {}
x = "y"
a[x] = 10
assert(a[x] == 10)
assert(a.x == nil)
assert(a.y == 10)

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

print("ok")

