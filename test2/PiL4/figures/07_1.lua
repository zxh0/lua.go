-- Figure 7.1. A program to sort a file

local lines = {}

-- read the lines in table 'lines'
for line in io.lines() do
  lines[#lines + 1] = line
end

-- sort
table.sort(lines)

-- write all the lines
for _, l in ipairs(lines) do
  io.write(l, "\n")
end
