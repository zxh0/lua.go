-- Figure 14.1. Multiplication of sparse matrices

function mult(a, b)
  local c = {} -- resulting matrix
  for i = 1, #a do
    local resultline = {}         -- will be 'c[i]'
    for k, va in pairs(a[i]) do   -- 'va' is a[i][k]
      for j, vb in pairs(b[k]) do -- 'vb' is b[k][j]
        local res = (resultline[j] or 0) + va * vb
        resultline[j] = (res ~= 0) and res or nil
      end
    end
    c[i] = resultline
  end
  return c
end
