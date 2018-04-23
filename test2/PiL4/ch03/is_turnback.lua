local tolerance = 10
function isturnback(angle)
  angle = angle % 360
  return (math.abs(angle - 180) < tolerance)
end

assert(isturnback(-180) == true)


local tolerance = 0.17
function isturnback(angle)
  --angle = angle % (2*math.pi)
  angle = angle % (math.pi*2)
  return (math.abs(angle - math.pi) < tolerance)
end


print("ok")
