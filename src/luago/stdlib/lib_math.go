package stdlib

import "math"
import "math/rand"
import . "luago/lua"

var mathLib = map[string]LuaGoFunction{
	"random":     mathRandom,
	"randomseed": mathRandomSeed,
	"max":        mathMax,
	"min":        mathMin,
	"exp":        mathExp,
	"log":        mathLog,
	"deg":        mathDeg,
	"rad":        mathRad,
	"sin":        mathSin,
	"cos":        mathCos,
	"tan":        mathTan,
	"asin":       mathAsin,
	"acos":       mathAcos,
	"atan":       mathAtan,
	"ceil":       mathCeil,
	"floor":      mathFloor,
	"fmod":       mathFmod,
	"modf":       mathModf,
	"abs":        mathAbs,
	"sqrt":       mathSqrt,
	"ult":        mathUlt,
	"tointeger":  mathToInt,
	"type":       mathType,
	/* placeholders */
	"pi":         nil,
	"huge":       nil,
	"maxinteger": nil,
	"mininteger": nil,
}

func OpenMathLib(ls LuaState) int {
	ls.NewLib(mathLib)
	ls.PushNumber(math.Pi)
	ls.SetField(-2, "pi")
	ls.PushNumber(math.Inf(1))
	ls.SetField(-2, "huge")
	ls.PushInteger(math.MaxInt64)
	ls.SetField(-2, "maxinteger")
	ls.PushInteger(math.MinInt64)
	ls.SetField(-2, "mininteger")
	return 1
}

/* pseudo-random numbers */

// math.random ([m [, n]])
// http://www.lua.org/manual/5.3/manual.html#pdf-math.random
func mathRandom(ls LuaState) int {
	switch ls.GetTop() {
	case 0:
		r := rand.Float64()
		ls.PushNumber(r)
		return 1
	case 1: // todo
		n := ls.ToInteger(1)
		r := 1 + rand.Int63n(n)
		ls.PushInteger(r)
		return 1
	case 2: // todo
		m := ls.ToInteger(1)
		n := ls.ToInteger(1)
		r := m + rand.Int63n(n-m)
		ls.PushInteger(r)
		return 1
	}

	panic("todo!")
}

// math.randomseed (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.randomseed
func mathRandomSeed(ls LuaState) int {
	if seed, ok := ls.ToIntegerX(1); ok {
		rand.Seed(seed)
		return 0
	} else {
		panic("todo!")
	}
}

/* max & min */

// math.max (x, ···)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.max
func mathMax(ls LuaState) int {
	return _maxOrMin(ls, true)
}

// math.min (x, ···)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.min
func mathMin(ls LuaState) int {
	return _maxOrMin(ls, false)
}

func _maxOrMin(ls LuaState, max bool) int {
	top := ls.GetTop()
	switch top {
	case 0:
		panic("todo!")
	case 1:
		return 1
	default: // todo
		idx := 1
		for i := 1; i < top; i++ {
			if max && ls.Compare(i, i+1, LUA_OPLT) || // arg[i] < arg[i+1] ?
				!max && ls.Compare(i+1, i, LUA_OPLT) { // arg[i+1] < arg[i] ?

				idx = i + 1
			}
		}
		ls.PushValue(idx)
		return 1
	}
}

/* exponentiation and logarithms */

// math.exp (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.exp
func mathExp(ls LuaState) int {
	x := ls.ToNumber(1)
	ls.PushNumber(math.Exp(x))
	return 1
}

// math.log (x [, base])
// http://www.lua.org/manual/5.3/manual.html#pdf-math.log
func mathLog(ls LuaState) int {
	x := ls.ToNumber(1)
	b := ls.ToNumber(2)
	l := math.Log(x) / math.Log(b)
	ls.PushNumber(l)
	return 1
}

/* trigonometric functions */

// math.deg (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.deg
func mathDeg(ls LuaState) int {
	x := ls.ToNumber(1)
	d := x * 180 / math.Pi
	ls.PushNumber(d)
	return 1
}

// math.rad (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.rad
func mathRad(ls LuaState) int {
	x := ls.ToNumber(1)
	r := x * math.Pi / 180
	ls.PushNumber(r)
	return 1
}

// math.sin (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.sin
func mathSin(ls LuaState) int {
	x := ls.ToNumber(1)
	ls.PushNumber(math.Sin(x))
	return 1
}

// math.cos (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.cos
func mathCos(ls LuaState) int {
	x := ls.ToNumber(1)
	ls.PushNumber(math.Cos(x))
	return 1
}

// math.tan (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.tan
func mathTan(ls LuaState) int {
	x := ls.ToNumber(1)
	ls.PushNumber(math.Tan(x))
	return 1
}

// math.asin (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.asin
func mathAsin(ls LuaState) int {
	x := ls.ToNumber(1)
	ls.PushNumber(math.Asin(x))
	return 1
}

// math.acos (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.acos
func mathAcos(ls LuaState) int {
	x := ls.ToNumber(1)
	ls.PushNumber(math.Acos(x))
	return 1
}

// math.atan (y [, x])
// http://www.lua.org/manual/5.3/manual.html#pdf-math.atan
func mathAtan(ls LuaState) int {
	y := ls.ToNumber(1)
	x := ls.OptNumber(2, 1.0)
	ls.PushNumber(math.Atan2(y, x))
	return 1
}

/* rounding functions */

// math.ceil (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.ceil
func mathCeil(ls LuaState) int {
	x := ls.ToNumber(1)
	ls.PushNumber(math.Ceil(x))
	return 1
}

// math.floor (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.floor
func mathFloor(ls LuaState) int {
	x := ls.ToNumber(1)
	ls.PushNumber(math.Floor(x))
	return 1
}

/* others */

// math.fmod (x, y)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.fmod
func mathFmod(ls LuaState) int {
	x := ls.ToNumber(1)
	y := ls.ToNumber(2)
	ls.PushNumber(math.Remainder(x, y))
	return 1
}

// math.modf (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.modf
func mathModf(ls LuaState) int {
	x := ls.ToNumber(1)
	i, f := math.Modf(x)
	ls.PushNumber(i)
	ls.PushNumber(f)
	return 2
}

// math.abs (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.abs
func mathAbs(ls LuaState) int {
	if ls.GetTop() != 1 {
		panic("todo!")
	}

	if ls.IsInteger(1) {
		x := ls.ToInteger(1)
		if x < 0 {
			ls.PushInteger(-x)
		}
		return 1
	}

	if x, ok := ls.ToNumberX(1); ok {
		ls.PushNumber(math.Abs(x))
		return 1
	}

	panic("todo!")
}

// math.sqrt (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.sqrt
func mathSqrt(ls LuaState) int {
	if ls.GetTop() != 1 {
		panic("todo!")
	}
	if x, ok := ls.ToNumberX(1); ok {
		ls.PushNumber(math.Sqrt(x))
		return 1
	} else {
		panic("todo!")
	}
}

// math.ult (m, n)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.ult
func mathUlt(ls LuaState) int {
	// todo
	m := uint64(ls.ToInteger(1))
	n := uint64(ls.ToInteger(2))
	ls.PushBoolean(m < n)
	return 1
}

// math.tointeger (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.tointeger
func mathToInt(ls LuaState) int {
	if i, ok := ls.ToIntegerX(1); ok {
		ls.PushInteger(i)
	} else {
		ls.PushNil()
	}
	return 1
}

// math.type (x)
// http://www.lua.org/manual/5.3/manual.html#pdf-math.type
func mathType(ls LuaState) int {
	if ls.IsInteger(1) {
		ls.PushString("integer")
	} else if ls.Type(1) == LUA_TNUMBER {
		ls.PushString("float")
	} else {
		ls.PushNil()
	}
	return 1
}
