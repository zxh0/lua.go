-- Figure 25.2. Hook for counting number of calls

local function hook ()
  local f = debug.getinfo(2, "f").func
  local count = Counters[f]
  if count == nil then -- first time 'f' is called?
    Counters[f] = 1
    Names[f] = debug.getinfo(2, "Sn")
  else -- only increment the counter
    Counters[f] = count + 1
  end
end
