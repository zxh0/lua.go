network = {
  {name = "grauna",  IP = "210.26.30.34"},
  {name = "arraial", IP = "210.26.30.23"},
  {name = "lua",     IP = "210.26.23.12"},
  {name = "derain",  IP = "210.26.23.20"},
}

table.sort(network, function(a, b) return (a.name > b.name) end)
assert(network[1].name == "lua")
assert(network[2].name == "grauna")
assert(network[3].name == "derain")
assert(network[4].name == "arraial")

print("ok")
