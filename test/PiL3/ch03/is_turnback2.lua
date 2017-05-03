local tolerance = 0.17
function isturnback(angle)
  --angle = angle % (2*math.pi)
  angle = angle % (math.pi*2)
  return (math.abs(angle - math.pi) < tolerance)
end
