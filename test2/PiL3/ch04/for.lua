days = {"Sunday", "Monday", "Tuesday", "Wednesday",
        "Thursday", "Friday", "Saturday"}
revDays = {}
for k, v in pairs(days) do
  revDays[v] = k
end
assert(revDays["Tuesday"] == 3)

--TODO
print("ok")
