t = {10, print, x = 12, k = "hi"}
for k, v in pairs(t) do
  -- print(k, v)
end

t = {10, print, 12, "hi"}
for k, v in ipairs(t) do
  -- print(k, v)
end

t = {10, print, 12, "hi"}
for k = 1, #t do
  -- print(k, t[k])
end

print("ok")
