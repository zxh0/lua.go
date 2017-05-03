t = {10, 20, 30}
table.insert(t, 1, 15)
assert(table.concat(t, ",") == "15,10,20,30")
print("ok")
