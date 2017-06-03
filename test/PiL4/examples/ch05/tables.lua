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

print("ok")
