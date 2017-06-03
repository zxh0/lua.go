-- Figure 17.2. A simple module for complex numbers

local M = {}         -- the module

-- creates a new complex number
local function new (r, i)
  return {r=r, i=i}
end

M.new = new        -- add 'new' to the module

-- constant 'i'
M.i = new(0, 1)

function M.add (c1, c2)
  return new(c1.r + c2.r, c1.i + c2.i)
end

function M.sub (c1, c2)
  return new(c1.r - c2.r, c1.i - c2.i)
end

function M.mul (c1, c2)
  return new(c1.r*c2.r - c1.i*c2.i, c1.r*c2.i + c1.i*c2.r)
end

local function inv (c)
  local n = c.r^2 + c.i^2
  return new(c.r/n, -c.i/n)
end

function M.div (c1, c2)
  return M.mul(c1, inv(c2))
end

function M.tostring (c)
  return string.format("(%g,%g)", c.r, c.i)
end

return M
