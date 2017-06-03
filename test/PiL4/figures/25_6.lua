-- Figure 25.6. Using hooks to bar calls to unauthorized functions

local debug = require "debug"

-- maximum "steps" that can be performed
local steplimit = 1000

local count = 0 -- counter for steps

-- set of authorized functions
local validfunc = {
  [string.upper] = true,
  [string.lower] = true,
  -- ... -- other authorized functions
}

local function hook (event)
  if event == "call" then
    local info = debug.getinfo(2, "fn")
    if not validfunc[info.func] then
      error("calling bad function: " .. (info.name or "?"))
    end
  end
  count = count + 1
  if count > steplimit then
    error("script uses too much CPU")
  end
end

-- load chunk
local f = assert(loadfile(arg[1], "t", {}))

debug.sethook(hook, "", 100) -- set hook

f() -- run chunk
