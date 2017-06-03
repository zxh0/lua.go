-- add the elements of sequence 'a'
function add1(a)
  local sum = 0
  for i = 1, #a do
    sum = sum + a[i]
  end
  return sum
end

function add2(...)
  local s = 0
  for i, v in ipairs{...} do
    s = s + v
  end
  return s
end

function add3(...)
  local s = 0
  for i = 1, select("#", ...) do
    s = s + select(i, ...)
  end
  return s 
end

assert(add1({3, 4, 10, 25, 12}) == 54)
assert(add2(3, 4, 10, 25, 12) == 54)
assert(add3(3, 4, 10, 25, 12) == 54)

-- print(select(1, "a", "b", "c"))   --> a    b    c
-- print(select(2, "a", "b", "c"))   --> b    c
-- print(select(3, "a", "b", "c"))   --> c
assert(select("#", "a", "b", "c") == 3)

print("ok")
