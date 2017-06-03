local mt = {}    -- create the matrix
for i = 1, N do
  local row = {} -- create a new row
  mt[i] = row
  for j = 1, M do
row[j] = 0 end
end

local mt = {} -- create the matrix
for i = 1, N do
  local aux = (i - 1) * M
  for j = 1, M do
    mt[aux + j] = 0
  end
end

for i = 1, M do
  for j = 1, N do
    c[i][j] = 0
    for k = 1, K do
      c[i][j] = c[i][j] + a[i][k] * b[k][j]
    end
  end
end

-- assumes 'c' has zeros in all elements
for i = 1, M do
  for k = 1, K do
    for j = 1, N do
      c[i][j] = c[i][j] + a[i][k] * b[k][j]
    end
  end
end
