co = coroutine.create(function() print("hi") end)
assert(type(co) == "thread")
assert(coroutine.status(co) == "suspended")
coroutine.resume(co) --> hi
assert(coroutine.status(co) == "dead")

co = coroutine.create(function()
  for i = 1, 10 do
    print("co", i)
    coroutine.yield()
  end
end)
coroutine.resume(co) --> co 1
assert(coroutine.status(co) == "suspended")
coroutine.resume(co) --> co 2
coroutine.resume(co) --> co 3
coroutine.resume(co) --> co 4
coroutine.resume(co) --> co 5
coroutine.resume(co) --> co 6
coroutine.resume(co) --> co 7
coroutine.resume(co) --> co 8
coroutine.resume(co) --> co 9
coroutine.resume(co) --> co 10
coroutine.resume(co) -- prints nothing
print(coroutine.resume(co)) --> false	cannot resume dead coroutine

co = coroutine.create(function(a, b, c)
  print("co", a, b, c + 2)
end)
coroutine.resume(co, 1, 2, 3) --> co 1 2 5

co = coroutine.create(function(a,b)
  coroutine.yield(a + b, a - b)
end)
print(coroutine.resume(co, 20, 10)) --> true 30 10

co = coroutine.create(function(x)
  print("co1", x)
  print("co2", coroutine.yield())
end)
coroutine.resume(co, "hi") --> co1 hi
coroutine.resume(co, 4, 5) --> co2 4 5

co = coroutine.create(function()
  return 6, 7
end)
print(coroutine.resume(co)) --> true 6 7

print("ok")
