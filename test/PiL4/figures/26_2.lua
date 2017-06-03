-- Figure 26.2. The dispatcher

tasks = {} -- list of all live tasks

function get (host, file)
  -- create coroutine for a task
  local co = coroutine.wrap(function ()
    download(host, file)
  end)
  -- insert it in the list
  table.insert(tasks, co)
end

function dispatch ()
  local i = 1
  while true do
    if tasks[i] == nil then -- no other tasks?
    if tasks[1] == nil then -- list is empty?
      break -- break the loop
    end
    i = 1 -- else restart the loop
  end
  local res = tasks[i]() -- run a task
  if not res then -- task finished?
    table.remove(tasks, i)
  else
    i = i + 1 -- go to next task
  end
  end
end
