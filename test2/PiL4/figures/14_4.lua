-- Figure 14.4. Finding a path between two nodes

function findpath (curr, to, path, visited)
  path = path or {}
  visited = visited or {}
  if visited[curr] then  -- node already visited?
    return nil           -- no path here
  end
  visited[curr] = true   -- mark node as visited
  path[#path + 1] = curr -- add it to path
  if curr == to then     -- final node?
    return path
  end
  -- try all adjacent nodes
  for node in pairs(curr.adj) do
    local p = findpath(node, to, path, visited)
    if p then return p end
  end
  table.remove(path)     -- remove node from path
end
