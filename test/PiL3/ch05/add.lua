function add(...)
  local s = 0
  for i, v in ipairs{...} do
    s = s + v
  end
  return s
end

assert(add(3, 4, 10, 25, 12) == 54)

print("ok")
