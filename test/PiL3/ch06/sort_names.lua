names = {"Peter", "Paul", "Mary"}
grades = {Mary = 10, Paul = 7, Peter = 8}
table.sort(names, function(n1, n2)
  return grades[n1] > grades[n2] -- compare the grades
end)
assert(names[1] == "Mary")
assert(names[2] == "Peter")
assert(names[3] == "Paul")

names = {"Peter", "Paul", "Mary"}
grades = {Mary = 10, Paul = 7, Peter = 8}
function sortbygrade(names, grades)
  table.sort(names, function (n1, n2)
    return grades[n1] > grades[n2] -- compare the grades
  end)
end
sortbygrade(names, grades)
assert(names[1] == "Mary")
assert(names[2] == "Peter")
assert(names[3] == "Paul")

print("ok")
