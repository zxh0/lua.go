-- Figure 24.3. An ugly implementation of the asynchronous I/O library

local cmdQueue = {} -- queue of pending operations

local lib = {}

function lib.readline (stream, callback)
  local nextCmd = function ()
    callback(stream:read())
  end
  table.insert(cmdQueue, nextCmd)
end

function lib.writeline (stream, line, callback)
  local nextCmd = function ()
    callback(stream:write(line))
  end
  table.insert(cmdQueue, nextCmd)
end

function lib.stop ()
  table.insert(cmdQueue, "stop")
end

function lib.runloop ()
  while true do
    local nextCmd = table.remove(cmdQueue, 1)
    if nextCmd == "stop" then
      break
    else
      nextCmd() -- perform next operation
    end
  end
end

return lib
