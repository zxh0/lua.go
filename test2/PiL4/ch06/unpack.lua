print(table.unpack{10,20,30}) --> 10   20   30
a,b = table.unpack{10,20,30}  -- a=10, b=20, 30 is discarded
assert(a == 10 and b == 20)

print(string.find("hello", "ll"))

f = string.find
a = {"hello", "ll"}
print(f(table.unpack(a)))

print(table.unpack({"Sun", "Mon", "Tue", "Wed"}, 2, 3)) --> Mon    Tue

function unpack(t, i, n)
  i = i or 1
  n = n or #t
  if i <= n then
    return t[i], unpack(t, i + 1, n)
  end
end

print("ok")
