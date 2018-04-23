function norm(x, y)
  return math.sqrt(x^2 + y^2)
end

function twice(x)
  return 2.0 * x
end

n = norm(3.4, 1.0)
n2 = twice(n)
--print(n2)
assert(math.abs(n2 - 7.0880180586677) < 0.01)

print("ok")
