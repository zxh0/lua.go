x = 10
local i = 1
while i <= x do
  local x = i*2
  --print(x)
  i = i + 1
end
if i > 20 then
  local x
  x = 20
  assert(x + 2 == 22)
else
  assert(x == 10)
end
assert(x == 10)

local a, b = 1, 10
if a < b then
  assert(a == 1)
  local a
  assert(a == nil)
end
assert(a == 1 and b == 10)

print("ok")
