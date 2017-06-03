assert(string.rep("abc", 3)           == "abcabcabc")
assert(string.reverse("A Long Line!") == "!eniL gnoL A")
assert(string.lower("A Long Line!")   == "a long line!")
assert(string.upper("A Long Line!")   == "A LONG LINE!")

s = "[in brackets]"
assert(string.sub(s, 2, -2)  == "in brackets")
assert(string.sub(s, 1, 1)   == "[")
assert(string.sub(s, -1, -1) == "]")

assert(string.char(97)                  == "a")
i = 99; assert(string.char(i, i+1, i+2) == "cde")
assert(string.byte("abc")               == 97)
assert(string.byte("abc", 2)            == 98)
assert(string.byte("abc", -1)           == 99)

a,b = string.byte("abc", 1, 2)
assert(a == 97 and b == 98)

assert(string.format("x = %d  y = %d", 10, 20) == "x = 10  y = 20")
assert(string.format("x = %x", 200)            == "x = c8")
assert(string.format("x = 0x%X", 200)          == "x = 0xC8")
assert(string.format("x = %f", 200)            == "x = 200.000000")
tag, title = "h1", "a title"
assert(string.format("<%s>%s</%s>", tag, title, tag)
    == "<h1>a title</h1>")	

assert(string.format("pi = %.4f", math.pi) == "pi = 3.1416")
d = 5; m = 11; y = 1990
assert(string.format("%02d/%02d/%04d", d, m, y) == "05/11/1990") 

i,j = string.find("hello world", "wor")
assert(i == 7 and j == 9)
assert(string.find("hello world", "war") == nil)

s,n = string.gsub("hello world", "l", ".")
assert(s == "he..o wor.d" and n == 3)
s,n = string.gsub("hello world", "ll", "..")
assert(s == "he..o world" and n == 1)
s,n = string.gsub("hello world", "a", ".")
assert(s == "hello world" and n == 0)

print("ok")
