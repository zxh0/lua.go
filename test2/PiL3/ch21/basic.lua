s = "[in brackets]"
-- assert(s:sub(2, -2) == "in brackets")
assert(string.sub(s, 2, -2) == "in brackets")

assert(string.char(97) == "a")
i = 99; assert(string.char(i, i+1, i+2) == "cde")
assert(string.byte("abc") == 97)
assert(string.byte("abc", 2) == 98)
assert(string.byte("abc", -1) == 99)

x, y = string.byte("abc", 1, 2)
assert(x == 97 and y == 98)

assert(string.format("pi = %.4f", math.pi) == "pi = 3.1416")
d = 5; m = 11; y = 1990
assert(string.format("%02d/%02d/%04d", d, m, y) == "05/11/1990")
tag, title = "h1", "a title"
assert(string.format("<%s>%s</%s>", tag, title, tag) == "<h1>a title</h1>")

print("ok")
