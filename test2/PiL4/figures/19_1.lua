-- Figure 19.1. Auxiliary definitions for the Markov program

function allwords()
  local line = io.read() -- current line
  local pos = 1          -- current position in the line
  return function()      -- iterator function
    while line do        -- repeat while there are lines
      local w, e = string.match(line, "(%w+[,;.:]?)()", pos)
      if w then          -- found a word?
        pos = e          -- update next position
        return w         -- return the word
      else
        line = io.read() -- word not found; try next line
        pos = 1 
      end
    end
    return nil           -- no more lines: end of traversal
  end
end

function prefix(w1, w2)
  return w1 .. " " .. w2
end

local statetab = {}

function insert(prefix, value)
  local list = statetab[prefix]
  if list == nil then
    statetab[prefix] = {value}
  else
    list[#list + 1] = value
  end
end
