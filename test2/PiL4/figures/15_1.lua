-- Figure 15.1. Quoting arbitrary literal strings

function quote(s)
  -- find maximum length of sequences of equals signs
  local n = -1
  for w in string.gmatch(s, "]=*%f[%]]") do
    n = math.max(n, #w - 1)   -- -1 to remove the ']'
  end

  -- produce a string with 'n' plus one equals signs
  local eq = string.rep("=", n + 1)

  -- build quoted string
  return string.format(" [%s[\n%s]%s] ", eq, s, eq)
end
