x = 10
local i = 1       -- local to the chunk

while i <= x do
  local x = i * 2 -- local to the while body
  print(x)        --> 2, 4, 6, 8, ...
  i = i + 1
end

if i > 20 then
  local x         -- local to the "then" body
  x = 20
  print(x + 2)    -- (would print 22 if test succeeded)
else
  print(x)        --> 10  (the global one)
end

print(x)          --> 10  (the global one)


--[[
local x1, x2
do
  local a2 = 2*a
  local d = (b^2 - 4*a*c)^(1/2)
  x1 = (-b + d)/a2
  x2 = (-b - d)/a2
end           -- scope of 'a2' and 'd' ends here
print(x1, x2) -- 'x1' and 'x2' still in scope
]]


local a, b = 1, 10
if a < b then
  print(a)  --> 1
  local a   -- '= nil' is implicit
  print(a)  --> nil
end         -- ends the block started at 'then'
print(a, b) --> 1   10


print("ok")
