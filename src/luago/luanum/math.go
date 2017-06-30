package luanum

import "math"

func CastToInteger(x float64) (int64, bool) {
	i, f := math.Modf(x)
	if f != 0 {
		return 0, false
	}
	if i >= 0 && i <= math.MaxInt64 ||
		i < 0 && i >= math.MinInt64 {
		return int64(i), true
	}
	return 0, false
}

func SHL(x, y int64) int64 {
	if y >= 0 {
		return x << uint64(y)
	} else {
		return int64(uint64(x) >> uint64(-y))
	}
}

func SHR(x, y int64) int64 {
	if y >= 0 {
		return int64(uint64(x) >> uint64(y))
	} else {
		return x << uint64(-y)
	}
}

func MOD(x, y int64) int64 {
	if x > 0 && y < 0 || x < 0 && y > 0 {
		return x%y + y
	} else {
		return x % y
	}
}

func IDIV(x, y int64) int64 {
	if x > 0 && y > 0 || x < 0 && y < 0 || x%y == 0 {
		return x / y
	} else {
		return x/y - 1
	}
}
