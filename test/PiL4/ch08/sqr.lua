function sqr(x)
  -- computes the square root of 'x' using Newton-Raphson method
  local sqr = x / 2
  repeat
    sqr = (sqr + x/sqr) / 2
    local error = math.abs(sqr^2 - x)
  until error < x/10000 -- local 'error' still visible here
  return sqr
end

assert(math.abs(sqr(16) - 4) < 0.0001)
print("ok")
