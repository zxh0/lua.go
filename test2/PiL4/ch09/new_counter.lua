function newCounter ()
  local count = 0
  return function () -- anonymous function
    count = count + 1
    return count
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
