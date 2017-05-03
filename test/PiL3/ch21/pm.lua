s = "hello world"
i, j = string.find(s, "hello")
assert(i == 1 and j == 5)
assert(string.sub(s, i, j) == "hello")
i, j = string.find(s, "world")
assert(i == 7 and j == 11)
i, j = string.find(s, "l")
assert(i == 3, j == 3)
assert(string.find(s, "lll") == nil)

assert(string.match("hello world", "hello") == "hello")
date = "Today is 17/7/1990"
d = string.match(date, "%d+/%d+/%d+")
assert(d == "17/7/1990")

s = string.gsub("Lua is cute", "cute", "great")
assert(s == "Lua is great")
s = string.gsub("all lii", "l", "x")
assert(s == "axx xii")
s = string.gsub("Lua is great", "Sol", "Sun")
assert(s == "Lua is great")

s = string.gsub("all lii", "l", "x", 1)
assert(s, "axl lii")
s = string.gsub("all lii", "l", "x", 2)
assert(s, "axx lii")

-- count = select(2, string.gsub(str, " ", " "))

print("ok")
