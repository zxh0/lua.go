lines = {
  luaH_set = 10,
  luaH_get = 24,
  luaH_present = 48,
}

a = {}
for n in pairs(lines) do
  a[#a + 1] = n
end
table.sort(a)
for _, n in ipairs(a) do
  print(n)
end


function pairsByKeys(t, f)
  local a = {}
  for n in pairs(t) do
    a[#a + 1] = n
  end
  table.sort(a, f)
  local i = 0
  return function()
    i = i + 1
    return a[i], t[a[i]]
  end
end

for name, line in pairsByKeys(lines) do
  print(name, line)
end

print("ok")
