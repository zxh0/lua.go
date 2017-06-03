local fact = function(n)
  if n == 0 then return 1
  else return n*fact(n-1) -- buggy
  end
end

local fact
fact = function(n)
  if n == 0 then return 1
    else return n*fact(n-1)
  end
end
assert(fact(10) == 3628800)

print("ok")
