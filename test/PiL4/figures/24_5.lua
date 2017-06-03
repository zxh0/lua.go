-- Figure 24.5. Running synchronous code on top of the asynchronous library

local lib = require "async-lib"

function run (code)
  local co = coroutine.wrap(function ()
    code()
    lib.stop()  -- finish event loop when done
  end)
  co()          -- start coroutine
  lib.runloop() -- start event loop
end

function putline (stream, line)
  local co = coroutine.running() -- calling coroutine
  local callback = (function () coroutine.resume(co) end)
  lib.writeline(stream, line, callback)
  coroutine.yield()
end

function getline (stream, line)
  local co = coroutine.running() -- calling coroutine
  local callback = (function (l) coroutine.resume(co, l) end)
  lib.readline(stream, callback)
  local line = coroutine.yield()
  return line
end
