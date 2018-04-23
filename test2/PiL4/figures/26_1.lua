-- Figure 26.1. Function to download a Web page

function download (host, file)
  local c = assert(socket.connect(host, 80))
  local count = 0 -- counts number of bytes read
  local request = string.format(
      "GET %s HTTP/1.0\r\nhost: %s\r\n\r\n", file, host)
  c:send(request)
  while true do
    local s, status = receive(c)
    count = count + #s
    if status == "closed" then break end
  end
  c:close()
  print(file, count)
end
