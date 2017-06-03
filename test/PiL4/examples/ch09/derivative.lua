function derivative(f, delta)
  delta = delta or 1e-4
  return function(x)
    return (f(x + delta) - f(x))/delta
  end
end
c = derivative(math.sin)
assert(tostring(math.cos(5.2)) == "0.46851667130038")
assert(tostring(c(5.2))        == "0.46856084325086")
assert(tostring(math.cos(10))  == "-0.83907152907645")
assert(tostring(c(10))         == "-0.83904432662041")

print("ok")
