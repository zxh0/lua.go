print "Hello World"
--dofile 'a.lua'
print [[a multi-line
 message]]
--f{x=10, y=20}
type{}

function f(a, b) print(a, b) end
f()        --> nil    nil
f(3)       --> 3      nil
f(3, 4)    --> 3      4
f(3, 4, 5) --> 3      4      (5 is discarded)

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
