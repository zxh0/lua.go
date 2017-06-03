function disk1(x, y)
  return (x - 1.0)^2 + (y - 3.0)^2 <= 4.5^2
end

function disk(cx, cy, r)
  return function(x, y)
    return (x - cx)^2 + (y - cy)^2 <= r^2
  end
end

function rect(left, right, bottom, up)
  return function(x, y)
    return left <= x and x <= right and
            bottom <= x and x <= up
  end
end

function complement(r)
  return function(x, y)
    return not r(x, y)
  end
end

function union(r1, r2)
  return function(x, y)
    return r1(x, y) or r2(x, y)
  end
end

function intersection(r1, r2)
  return function(x, y)
    return r1(x, y) and r2(x, y)
  end
end

function difference (r1, r2)
  return function(x, y)
    return r1(x, y) and not r2(x, y)
  end
end

function translate (r, dx, dy)
  return function(x, y)
    return r(x - dx, y - dy)
  end
end

print("ok")
