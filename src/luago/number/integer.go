package number

import "regexp"
import "strings"
import "strconv"

var reInteger = regexp.MustCompile(`^-?[0-9]+$|^-?0x[0-9a-f]+$`)

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

func ParseInteger(str string) (int64, bool) {
	str = strings.ToLower(str)
	if !reInteger.MatchString(str) { // float?
		return 0, false
	}
	if strings.Index(str, "0x") < 0 { // decimal
		i, err := strconv.ParseInt(str, 10, 64)
		return i, err == nil
	}

	// hex
	var sign int64 = 1
	if str[0] == '-' {
		sign = -1
		str = str[3:]
	} else {
		str = str[2:]
	}

	if len(str) > 16 {
		str = str[len(str)-16:] // cut long hex string
	}

	i, err := strconv.ParseUint(str, 16, 64)
	return sign * int64(i), err == nil
}
