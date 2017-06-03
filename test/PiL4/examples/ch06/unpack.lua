function unpack(t, i, n)
  i = i or 1
  n = n or #t
  if i <= n then
    return t[i], unpack(t, i + 1, n)
  end
end

a,b,c = table.unpack{10,20,30} --> 10   20   30
assert(a == 10 and b == 20 and c == 30)
a,b = table.unpack{10,20,30}   -- a=10, b=20, 30 is discarded
assert(a == 10 and b == 20)

f = string.find
a = {"hello", "ll"}
i,j = f(table.unpack(a))
assert(i == 3 and j == 4)

a,b = table.unpack({"Sun", "Mon", "Tue", "Wed"}, 2, 3)
assert(a == "Mon" and b == "Tue")

print("ok")
