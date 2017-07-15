function traceback()
  for level = 1, math.huge do
    local info = debug.getinfo(level, "Sl")
    if not info then break end
    if info.what == "C" then -- is a C function?
      print(string.format("%d\tC function", level))
    else -- a Lua function
      print(string.format("%d\t[%s]:%d", level, info.short_src, info.currentline))
    end
  end
end

traceback()
-- print(debug.traceback())
