-- Figure 21.3. Accounts using a dual representation

local balance = {}

Account = {}

function Account:withdraw (v)
  balance[self] = balance[self] - v
end

function Account:deposit (v)
  balance[self] = balance[self] + v
end

function Account:balance ()
  return balance[self]
end

function Account:new (o)
  o = o or {} -- create table if user does not provide one
  setmetatable(o, self)
  self.__index = self
  balance[o] = 0 -- initial balance
  return o
end
