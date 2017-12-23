local t = {1,2,3}
assert(#t == 3)

t = {1,2,3,nil}
assert(#t == 3)

t = {1,2,3,nil,4,nil}
print(#t)

t = {1,2,3,4,5}
t[4] = nil
t[5] = nil 
print(#t)