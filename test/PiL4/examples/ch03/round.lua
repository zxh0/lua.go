function round (x)
  local f = math.floor(x)
  if (x == f) or (x % 2.0 == 0.5) then
    return f
  else
    return math.floor(x + 0.5)
  end
end

assert(round(2.5)  == 2)
assert(round(3.5)  == 4)
assert(round(-2.5) == -2)
assert(round(-1.5) == -2)

print("ok")
