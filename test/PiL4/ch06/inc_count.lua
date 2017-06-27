function incCount(n)
  n = n or 1
  globalCounter = globalCounter + n
end

globalCounter = 1
incCount(1)
assert(globalCounter == 2)
incCount(2)
assert(globalCounter == 4)
incCount()
assert(globalCounter == 5)

print("ok")
