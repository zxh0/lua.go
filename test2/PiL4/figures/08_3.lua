-- Figure 8.3. A strange (and invalid) use of a goto

function getlabel ()
  return function () goto L1 end
  ::L1::
  return 0
end

function f (n)
  if n == 0 then return getlabel()
  else
    local res = f(n - 1)
    print(n)
    return res
  end
end

x = f(10)
x()
