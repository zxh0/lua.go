-- Figure 25.5. Controlling memory use

-- maximum memory (in KB) that can be used
local memlimit = 1000

-- maximum "steps" that can be performed
local steplimit = 1000

local function checkmem ()
  if collectgarbage("count") > memlimit then
    error("script uses too much memory")
  end
end

local count = 0
local function step ()
  checkmem()
  count = count + 1
  if count > steplimit then
    error("script uses too much CPU")
  end
end
