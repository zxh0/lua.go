function f(a, b)
  return a .. tostring(b)
end
assert(f(3) == "3nil")
assert(f(3, 4) == "34")
assert(f(3, 4, 5) == "34")

print("ok")
