function incCount(n)
  n = n or 1
  count = count + n
end

count = 1
incCount(1)
assert(count == 2)
incCount(2)
assert(count == 4)
incCount()
assert(count == 5)

print("ok")
