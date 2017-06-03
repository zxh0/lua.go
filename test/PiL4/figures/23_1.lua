-- Figure 23.1. Constant-function factory with memorization

do
  local mem = {} -- memorization table
  setmetatable(mem, {__mode = "k"})
  function factory (o)
    local res = mem[o]
    if not res then
      res = (function () return o end)
      mem[o] = res
    end
    return res
  end
end
