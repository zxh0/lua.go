-- Figure 26.3. Dispatcher using select

function dispatch ()
  local i = 1
  local timedout = {}
  while true do
    if tasks[i] == nil then -- no other tasks?
      if tasks[1] == nil then -- list is empty?
        break -- break the loop
      end
      i = 1 -- else restart the loop
      timedout = {}
    end
    local res = tasks[i]() -- run a task
    if not res then -- task finished?
      table.remove(tasks, i)
    else -- time out
      i = i + 1
      timedout[#timedout + 1] = res
      if #timedout == #tasks then -- all tasks blocked?
        socket.select(timedout) -- wait
      end
    end
  end
end
