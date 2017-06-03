-- Figure 25.4. A naive sandbox with hooks

local debug = require "debug"

-- maximum "steps" that can be performed
local steplimit = 1000

local count = 0 -- counter for steps

local function step ()
  count = count + 1
  if count > steplimit then
    error("script uses too much CPU")
  end
end

-- load file
local f = assert(loadfile(arg[1], "t", {}))

debug.sethook(step, "", 100) -- set hook

f() -- run file
