local u,v,w
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
  local va = u
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
  local a = b[u]
  local a = b[x]
  local a = u[1]
  local a = u.foo
  local a = u[v]
  local a = u[x]
  local a = x[1]
  local a = x.foo
  local a = x.foo.bar
  local a = x[u]
  local a = x[u][v]
  local a = x[y]
  local a = x[y][z]
  local a = b + c
  local a = v1 - v2 - v3
  local a = v1 * v2 * v3 * v4
  local a = b + 1
  local a = u + v
  local a = u + v + w
  local a = x + y
  local a = x + y + z
  local a = b / u / x / 1
  local a = b ^ u ^ x ^ 1
  local a = b ^ b ^ b ^ b
  --local a = 1 ^ x ^ u ^ b ^ c
  local a = b .. c .. u .. v .. x .. 1
  local a = 1 < 2
  local a = b == c
  local a = b ~= c
  local a = b > c
  local a = b < c
  local a = b >= c
  local a = b <= c
  local a = v1 < v2 < v3
  local a = u ~= v
  local a = u > v > w
  local a = x == y
  local a = x > y > z
  local a = b == c ~= u > x >= y < 1 <= false
  local a = 1 or 2
  local a = b or c
  local a = u or v
  local a = x or y
  local a = x or y or z or y or x
  local a = b or u or x or true
  local a = b and c
  local a = v1 and v2 and v3
  local a = u and v and w
  local a = x and y and z
  local a = b and u and x and true
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
  a = u
  a = x
  a = f()
  a = a%360
  u = nil
  v = true
  w = 100
  u = "foo"
  v = {}
  w = ...
  u = a
  u = v
  u = _ENV.x
  u = x
  x = nil
  x = false
  x = 100
  x = "foo"
  x = {}
  x = ...
  x = a
  x = u
  x = y
  a[nil] = nil
  a[true] = false
  a[1] = 2
  a[1] = b[2]
  a[b] = c
  a[b] = b[a]
  a[u] = v
  a[u] = v[b]
  b[x] = y[a]
  u[nil] = nil
  u[false] = true
  u[1] = v[2]
  u[a] = b[v]
  u[v] = v[u]
  v[x] = x[v]
  x[nil] = nil
  x[true] = false
  x[1] = y[2]
  x[b] = c[y]
  y[v] = u[x]
  y[x] = x[y]
  c[b][a] = w[v][u]
  c[b[a]] = w[v[u]]
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
  a,u,x = y,b,v
  -- u,v,w = a,b,c
  -- a,a,a = 1,2,3
  -- u,v,x = 1,2,3
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
  local a,b,c
  t = {[a]=1, [u]=2, [x]=3, [5]=4,}
  t = {x=a, y=u, z=x,}
  days = {"Sunday", "Monday", "Tuesday", "Wednesday",
        "Thursday", "Friday", "Saturday"}
  w = {x=0, y=0, label="console"}
  opnames = {["+"] = "add", ["-"] = "sub",
           ["*"] = "mul", ["/"] = "div"}
  a = {[i+0] = s, [i+1] = s..s, [i+2] = s..s..s}
  a = {[1]="red", [2]="green", [3]="blue",}
  a = {x=10, y=45; "one", "two", "three"}
  -- polyline = {
  --   color="blue",
  --   thickness=2,
  --   npoints=4,
  --   {x=0, y=0},   -- polyline[1]
  --   {x=-10, y=0}, -- polyline[2]
  --   {x=-10, y=1}, -- polyline[3]
  --   {x=0, y=1}    -- polyline[4]
  -- }
end --]==]
