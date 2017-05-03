-- for n in pairs(_G) do
--   print(n)
-- end

assert(_G["print"] == print)

_G["hello"] = "world"
assert(hello == "world")

function getfield(f)
  local v = _G -- start with the table of globals
  for w in string.gmatch(f, "[%w_]+") do
    v = v[w]
  end
  return v
end
assert(getfield("table.sort") == table.sort)

function setfield (f, v)
  local t = _G          -- start with the table of globals
  for w, d in string.gmatch(f, "([%w_]+)(%.?)") do
    if d == "." then    -- not last name?
      t[w] = t[w] or {} -- create table if absent
      t = t[w]          -- get the table
    else                -- last name
      t[w] = v          -- do the assignment
    end
  end
end
setfield("t.x.y", 10)
assert(t.x.y == 10)
assert(getfield("t.x.y") == 10)

print("ok")
