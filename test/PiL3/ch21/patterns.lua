s = "Deadline is 30/05/1999, firm"
date = "%d%d/%d%d/%d%d%d%d"
assert(string.sub(s, string.find(s, date)) == "30/05/1999")

s, n = string.gsub("hello, up-down!", "%A", ".")
assert(s == "hello..up.down." and n == 4)

s, n = string.gsub("one, and two; and three", "%a+", "word")
assert(s == "word, word word; word word", n == 5)

assert(string.match("the number 1298 is even", "%d+") == "1298")

test = "int x; /* x */ int y; /* y */"
assert(string.match(test, "/%*.*%*/") == "/* x */ int y; /* y */")

test = "int x; /* x */ int y; /* y */"
assert(string.gsub(test, "/%*.-%*/", "") == "int x;  int y; ")

-- s = "a (enclosed (in) parentheses) line"
-- print(string.gsub(s, "%b()", "")) --> a line

-- s = "the anthem is the theme"
-- print(s:gsub("%f[%w]the%f[%W]", "one"))
-- --> one anthem is one theme


pair = "name = Anna"
key, value = string.match(pair, "(%a+)%s*=%s*(%a+)")
assert(key == "name", value == "Anna")

date = "Today is 17/7/1990"
d, m, y = string.match(date, "(%d+)/(%d+)/(%d+)")
assert(d == "17" and m == "7" and y == "1990")

-- s = [[then he said: "it's all right"!]]
-- q, quotedPart = string.match(s, "([\"'])(.-)%1")
-- print(quotedPart) --> it's all right
-- print(q) --> "

-- p = "%[(=*)%[(.-)%]%1%]"
-- s = "a = [=[[[ something ]] ]==] ]=]; print(a)"
-- print(string.match(s, p)) --> = [[ something ]] ]==]

s, n = string.gsub("hello Lua!", "%a", "%0-%0")
assert(s == "h-he-el-ll-lo-o L-Lu-ua-a!" and n == 8)

s, n = string.gsub("hello Lua", "(.)(.)", "%2%1")
assert(s == "ehll ouLa" and n == 4)

s = [[the \quote{task} is to \em{change} that.]]
s, n = string.gsub(s, "\\(%a+){(.-)}", "<%1>%2</%1>")
assert(s == "the <quote>task</quote> is to <em>change</em> that." and n == 2)

print("ok")
