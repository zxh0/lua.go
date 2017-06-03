-- Figure 13.2. Dumping the dump program

local f = assert(io.open(arg[1], "rb"))
local blocksize = 16
  for bytes in f:lines(blocksize) do
    for i = 1, #bytes do
      local b = string.unpack("B", bytes, i)
      io.write(string.format("%02X ", b))
    end
  io.write(string.rep("   ", blocksize - #bytes))
  bytes = string.gsub(bytes, "%c", ".")
  io.write(" ", bytes, "\n")
end
