function newCounter ()
  local i = 0
  return function () -- anonymous function
    i = i + 1
    return i
  end
end

c1 = newCounter()
assert(c1() == 1) --> 1
assert(c1() == 2) --> 2

c2 = newCounter()
assert(c2() == 1) --> 1
assert(c1() == 3) --> 3
assert(c2() == 2) --> 2

print("ok")
