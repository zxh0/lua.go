f = load("function add(a, b) return a + b end")
f()
print(add(1, 2))

print(type(_G))
print(_VERSION)


s = "你好，世界！"
print(string.len(s))           --> 18
print(string.len(s))           --> 18
print(utf8.len(s))             --> 6
print(utf8.offset(s, 2))       --> 4
print(utf8.offset(s, 1, 4))    --> 4
print(utf8.codepoint(s, 1))    --> 20320
print(utf8.codepoint(s, 4, 8)) --> 22909 65292
print(utf8.char(20320, 22909)) --> 你好
for p, c in utf8.codes(s) do print(p, c) end


co = coroutine.create(function()end)
assert(not coroutine.isyieldable(co))
assert(not coroutine.isyieldable(coroutine.running()))