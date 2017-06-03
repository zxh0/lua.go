-- Figure 24.4. Reversing a file in event-driven fashion

local lib = require "async-lib"

local t = {}
local inp = io.input()
local out = io.output()
local i

-- write-line handler
local function putline ()
  i = i - 1
    if i == 0 then -- no more lines?
    lib.stop()     -- finish the main loop
  else             -- write line and prepare next one
    lib.writeline(out, t[i] .. "\n", putline)
  end
end

-- read-line handler
local function getline (line)
  if line then                 -- not EOF?
    t[#t + 1] = line           -- save line
    lib.readline(inp, getline) -- read next one
  else                         -- end of file
    i = #t + 1                 -- prepare write loop
    putline()                  -- enter write loop
  end
end

lib.readline(inp, getline) -- ask to read first line
lib.runloop()              -- run the main loop
