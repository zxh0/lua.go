local a = {} -- new array
for i = 1, 1000 do
  a[i] = 0
end
assert(#a == 1000)

-- creates an array with indices from -5 to 5
a = {}
for i = -5, 5 do
  a[i] = 0
end
assert(#a == 5)

squares = {1, 4, 9, 16, 25, 36, 49, 64, 81}

print("ok")
