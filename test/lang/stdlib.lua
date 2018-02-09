f = load("function add(a, b) return a + b end")
f()
print(add(1, 2))

print(type(_G))
print(_VERSION)