-- Figure 9.2. Drawing a region in a PBM file

function plot (r, M, N)
  io.write("P1\n", M, " ", N, "\n") -- header
  for i = 1, N do -- for each line
    local y = (N - i*2)/N
    for j = 1, M do -- for each column
      local x = (j*2 - M)/M
      io.write(r(x, y) and "1" or "0")
    end
    io.write("\n")
  end
end
