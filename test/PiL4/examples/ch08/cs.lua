-- Control Structures

function if_then_else()
  if a < 0 then a = 0 end
  if a < b then return a else return b end
  if line > MAXLINES then
    showpage()
    line = 0
  end
  
  if op == "+" then
    r=a+ b
  elseif op == "-" then
    r=a- b
  elseif op == "*" then
    r = a*b
  elseif op == "/" then
    r = a/b
  else
    error("invalid operation")
  end
end

function _while()
  local i = 1
  while a[i] do
    print(a[i])
    i=i+ 1
  end
end

function _repeat()
  -- print the first non-empty input line
  local line
  repeat
    line = io.read()
  until line ~= ""
  print(line)

  -- computes the square root of 'x' using Newton-Raphson method
  local sqr = x / 2
  repeat
    sqr = (sqr + x/sqr) / 2
    local error = math.abs(sqr^2 - x)
  until error < x/10000 -- local 'error' still visible here
end

function for_num()
  for i = 1, math.huge do
    if (0.3*i^3 - 20*i^2 - 500 >= 0) then
      print(i)
      break
    end
  end

  for i = 1, 10 do print(i) end
  max = i      -- probably wrong!

  -- find a value in a list
  local found = nil
  for i = 1, #a do
    if a[i] < 0 then
      found =i -- save value of 'i' break
    end
  end
  print(found)
end

function while2()
  local i = 1
  while a[i] do
    if a[i] == v then return i end
    i=i+ 1
  end
end

print("ok")
