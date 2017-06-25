function a(   )            end -- IsVararg==0
function d(...) print(...) end -- IsVararg==1
function b(...)            end -- IsVararg==2
function c(...) print(123) end -- IsVararg==2
function e(...) print(arg) end -- IsVararg==2
print(...)
