days = {"Sunday", "Monday", "Tuesday", "Wednesday",
        "Thursday", "Friday", "Saturday"}
assert(days[4] == "Wednesday")

a = {x = 10, y = 20}
a = {}; a.x = 10; a.y = 20

w = {x=0, y=0, label="console"}
x = {
  math.sin(0), math.sin(1), math.sin(2)
}
w[1] = "another field"
x.f = w
assert(w["x"] == 0)
assert(w[1]   == "another field")
assert(x.f[1] == "another field")
w.x = nil

polyline = {
  color="blue",
  thickness=2,
  npoints=4,
  {x=0,   y=0}, -- polyline[1]
  {x=-10, y=0}, -- polyline[2]
  {x=-10, y=1}, -- polyline[3]
  {x=0,   y=1}  -- polyline[4]
}
assert(polyline[2].x == -10)
assert(polyline[4].y == 1)

opnames = {["+"] = "add", ["-"] = "sub",
           ["*"] = "mul", ["/"] = "div"}
i = 20; s = "-"
a = {[i+0] = s, [i+1] = s..s, [i+2] = s..s..s}
assert(opnames[s] == "sub")
assert(a[22]      == "---")

a = {[1]="red", [2]="green", [3]="blue",}
-- a = {x=10, y=45; "one", "two", "three"}

print("ok")
