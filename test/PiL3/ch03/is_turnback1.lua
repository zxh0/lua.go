local tolerance = 10
function isturnback(angle)
  angle = angle % 360
  return (math.abs(angle - 180) < tolerance)
end

assert(isturnback(-180) == true)

-- function mod(a, b) return a % b end
-- assert(mod( 2,  5) ==  2)
-- assert(mod( 5,  2) ==  1)
-- assert(mod( 5, -2) == -1)
-- assert(mod(-5,  2) ==  1)
-- assert(mod(-5, -1) ==  0)

print("ok")
