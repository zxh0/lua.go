a = "one string"
b = string.gsub(a, "one", "another")  -- change string parts
assert(a == "one string")
assert(b == "another string")

a = "hello"
assert(#a == 5)
assert(#"good bye" == 8)

assert("Hello " .. "World" == "Hello World")
assert("result is " .. 3   == "result is 3")

a = "Hello"
assert(a .. " World" == "Hello World")
assert(a == "Hello")

a = "a line"
b = 'another line'
c = 'a backslash inside quotes: \'\\\''
d = "a simpler way: '\\'"
e = "\u{3b1} \u{3b2} \u{3b3}"
assert("ALO\n123\"" == '\x41LO\10\04923"')
assert("ALO\n123\"" == '\x41\x4c\x4f\x0a\x31\x32\x33\x22')

-- Long strings
page = [[
<html>
<head>
<title>An HTML Page</title>
</head>
<body>
<a href="http://www.lua.org">Lua</a>
</body>
</html>
]]

data = "\x00\x01\x02\x03\x04\x05\x06\x07\z
        \x08\x09\x0A\x0B\x0C\x0D\x0E\x0F"

print("ok")
