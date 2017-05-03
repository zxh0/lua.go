function sigma(n)
  local s = 0
  for i = 1, n do
    s = s + i
  end
  return s
end

local x = sigma(100)
print(x)
