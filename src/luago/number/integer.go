package number

import "regexp"
import "strings"
import "strconv"

var reInteger = regexp.MustCompile(`^[+-]?[0-9]+$|^-?0x[0-9a-f]+$`)

func ParseInteger(str string, base int) (int64, bool) {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)

	if base != 10 {
		i, err := strconv.ParseInt(str, base, 64)
		return i, err == nil
	}

	if !reInteger.MatchString(str) { // float?
		return 0, false
	}
	if str[0] == '+' {
		str = str[1:]
	}
	if strings.Index(str, "0x") < 0 { // decimal
		i, err := strconv.ParseInt(str, base, 64)
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
