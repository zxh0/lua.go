local function fact(n)
  if n == 0 then return 1
    else return n*fact(n-1) -- buggy
  end
end
assert(fact(10) == 3628800)
print("ok")
