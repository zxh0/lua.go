Lib = {}
Lib.foo = function(x,y) return x + y end
Lib.goo = function(x,y) return x - y end
assert(Lib.foo(2, 3) == 5 and Lib.goo(2, 3) == -1)

Lib = {
  foo = function (x,y) return x + y end,
  goo = function (x,y) return x - y end
}
assert(Lib.foo(2, 3) == 5 and Lib.goo(2, 3) == -1)


Lib = {}
function Lib.foo (x,y) return x + y end
function Lib.goo (x,y) return x - y end
assert(Lib.foo(2, 3) == 5 and Lib.goo(2, 3) == -1)

print("ok")
