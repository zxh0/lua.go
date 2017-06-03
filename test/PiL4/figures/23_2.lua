-- Figure 23.2. Running a function at every GC cycle

do
  local mt = {__gc = function (o)
    -- whatever you want to do
    print("new cycle")
    -- creates new object for next cycle
    setmetatable({}, getmetatable(o))
  end}
  -- creates first object
  setmetatable({}, mt)
end

collectgarbage() --> new cycle
collectgarbage() --> new cycle
collectgarbage() --> new cycle
