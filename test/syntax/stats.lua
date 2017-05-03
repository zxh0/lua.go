local u1,u2,u3
--[==[
function stat_local_assign(...)
  local v1 = nil
  local v2 = true
  local v3 = false
  local v4 = 100
  local v5 = "foo"
  local v6 = {}
  local v7 = ...
  local v8 = v7
  local v9 = v9
  local va = u1
  local vb = x
  local v1,v2,v3 = 1,2,3
  local v1,v1,v1 = a,b,c
  local v1,v2,v1 = a,b,c
  local a,b,c
  local a,b,c = 1
  local a,b,c = 1,2
  local a,b,c = 1,2,3,4,5
  local a,b,c = nil, "bar"
  local a,b,c = ...
  local a,b,c = ...,...
  local a,b,c = ...,...,...
  local a,b,c = ...,...,...,...
  local a,b,c = 1,...
  local a,b,c = 1,2,...
  local a,b,c = 1,2,3,...
  local a,b,c = 1,2,3,4,...
  local a,b,c = 1,2,3,4,...,5
  local a,b,c = 1,...,2,3,4,5
  local a,b,c = (...)
  local a,b,c = c,b,a
  local a,b,c = c,c,c,4,5
  local a,b,c = a,a,a,a,a,a
  local a = f()
  local a,b,c = f()
  local a,b,c = f(),f()
  local a,b,c = f(),f(),f()
  local a,b,c = f(),f(),f(),f()
  local a,b,c = 1,f()
  local a,b,c = 1,2,f()
  local a,b,c = 1,2,3,f()
  local a,b,c = 1,2,3,4,f()
  local a,b,c = 1,2,3,4,f(),5
  local a,b,c = 1,f(),2,3,4,5
  local a,b,c = a(b,c)
  local f = function() end
  local a = b[1]
  local a = b["foo"]
  local a = b.foo
  local a = b.c
  local a = b[u1]
  local a = b[x]
  local a = u1[1]
  local a = u1.foo
  local a = u1[u2]
  local a = u1[x]
  local a = x[1]
  local a = x.foo
  local a = x.foo.bar
  local a = x[u1]
  local a = x[u1][u2]
  local a = x[y]
  local a = x[y][z]
  local a = b + c
  local a = v1 - v2 - v3
  local a = v1 * v2 * v3 * v4
  local a = b + 1
  local a = u1 + u2
  local a = u1 + u2 + u3
  local a = x + y
  local a = x + y + z
  local a = b / u1 / x / 1
  local a = b ^ u1 ^ x ^ 1
  local a = b ^ b ^ b ^ b
  --local a = 1 ^ x ^ u1 ^ b ^ c
  local a = b .. c .. u1 .. u2 .. x .. 1
  local a = 1 < 2
  local a = b == c
  local a = b ~= c
  local a = b > c
  local a = b < c
  local a = b >= c
  local a = b <= c
  local a = v1 < v2 < v3
  local a = u1 ~= u2
  local a = u1 > u2 > u3
  local a = x == y
  local a = x > y > z
  local a = b == c ~= u1 > x >= y < 1 <= false
  local a = 1 or 2
  local a = b or c
  local a = u1 or u2
  local a = x or y
  local a = x or y or z or y or x
  local a = b or u1 or x or true
  local a = b and c
  local a = v1 and v2 and v3
  local a = u1 and u2 and u3
  local a = x and y and z
  local a = b and u1 and x and true
  -- local a = x and y or x and y
  -- local a = 1
  a=1
end --]==]
--[==[
function stat_assign_1(...)
  local a,b,c,d=1,2,3,4
  a = nil
  b = true
  c = false
  a = 100
  b = "foo"
  c = {}
  a = ...
  a = b
  a = u1
  a = x
  a = f()
  a = a%360
  u1 = nil
  u2 = true
  u3 = 100
  u1 = "foo"
  u2 = {}
  u3 = ...
  u1 = a
  u1 = u2
  u1 = _ENV.x
  u1 = x
  x = nil
  x = false
  x = 100
  x = "foo"
  x = {}
  x = ...
  x = a
  x = u1
  x = y
  a[nil] = nil
  a[true] = false
  a[1] = 2
  a[1] = b[2]
  a[b] = c
  a[b] = b[a]
  a[u1] = u2
  a[u1] = u2[b]
  b[x] = y[a]
  u1[nil] = nil
  u1[false] = true
  u1[1] = u2[2]
  u1[a] = b[u2]
  u1[u2] = u2[u1]
  u2[x] = x[u2]
  x[nil] = nil
  x[true] = false
  x[1] = y[2]
  x[b] = c[y]
  y[u2] = u1[x]
  y[x] = x[y]
  c[b][a] = u3[u2][u1]
  c[b[a]] = u3[u2[u1]]
  b[1][2] = y[true][false]
  a,b,c = nil
  a,b,c = 1,2,3
end --]==]
--[==[
function stat_assign_n(...)
  local a,b,c = 1,2,3
  a,b,c = nil
  a,b,c = 1
  a,b,c = 1,2
  a,b,c = 1,2,3
  a,b,c = 1,2,3,4
  a,u1,x = y,b,u2
  -- u1,u2,u3 = a,b,c
  -- a,a,a = 1,2,3
  -- u1,u2,x = 1,2,3
  -- x,y,z[1] = 1,2,3
end --]==]
--[==[
function stat_fc(...)
  f()
  f(a)
  f(a, b)
  f(a, b, c)
  f.g.h()
  f.g:h()
  f[g][h]()
  f[g[h]]()
  f.g.h(a.b.c)
  f(g())
  f(g(), 1)
  f(1, g())
  return f(g())
end --]==]
--[==[
function stat_return(...)
  local a,b,c = 1,2,3
  return a,b,c
end --]==]
--[==[
function stat_if(...)
  local a = 1,(x == 0)
  if true   then print(1) end
  if 1234   then print(2) end
  if "ab"   then print(3) end
  if false  then print(4) end
  if nil    then print(4) end
  if x      then print(5) end
  if x == 0 then print(6) end
  if x ~= 0 then print(6) end
  if x >  0 then print(6) end
  if x >= 0 then print(6) end
  if x <  0 then print(6) end
  if x <= 0 then print(6) end
  if 0 == x then print(7) end
  if x == y then print(8) end
  if 1 == 2 then print(9) end
  if n == 0 then return 1 else return fact(n-1) end
  if x then print(a) elseif y then print(b) end
  if x then print(a) elseif y then print(b) elseif z then print(c) end
end --]==]
--[==[
function stat_for_num(...)
  for i = 0, 100, 2  do print(i) end
  for i = 0, 100     do print(i) end
  for i = 100, 0, -1 do print(i) end
end --]==]
--[==[
function stat_while(...)
  while true  do print(1) end
  while 1234  do print(2) end
  while x     do print(3) end
  while false do print(4) end
  while nil   do print(5) end
end --]==]
--[==[
function stat_repeat(...)
  repeat print(1) until true
  repeat print(2) until 1234
  repeat print(3) until x
  repeat print(3) until false
  repeat print(3) until nil
end --]==]
--[==[
function tc(...)
  -- body
end --]==]
