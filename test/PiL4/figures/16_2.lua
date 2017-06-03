-- Figure 16.2. String repetition

function stringrep (s, n)
  local r = ""
  if n > 0 then
    while n > 1 do
      if n % 2 ~= 0 then
      s = s .. s
      n = math.floor(n / 2)
    end
    r = r .. s 
  end
  return r 
end
