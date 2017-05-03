function derivative(f, delta)
  delta = delta or 1e-4
  return function(x)
    return (f(x + delta) - f(x))/delta
  end
end
c = derivative(math.sin)
print(math.cos(5.2), c(5.2)) --> 0.46851667130038 0.46856084325086
print(math.cos(10), c(10)) --> -0.83907152907645 -0.83904432662041
