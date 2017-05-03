function newCounter ()
	do
	local i = 0
	end
	return function () -- anonymous function
		i = i + 1
		return i
	end
end
c1 = newCounter()
print(c1()) --> 1
print(c1()) --> 2