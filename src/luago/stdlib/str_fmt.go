package stdlib

import "regexp"
import "strings"

// tag = %[flags][width][.precision]specifier
var tagPattern = regexp.MustCompile(`%[ #+-0]?[0-9]*(\.[0-9]+)?[cdeEfgGioqsuxX%]`)

func parseFmtStr(fmt string) []string {
	if fmt == "" || strings.IndexByte(fmt, '%') < 0 {
		return []string{fmt}
	}

	parsed := make([]string, 0, len(fmt)/2)
	for {
		if fmt == "" {
			break
		}

		loc := tagPattern.FindStringIndex(fmt)
		if loc == nil {
			parsed = append(parsed, fmt)
			break
		}

		head := fmt[:loc[0]]
		tag := fmt[loc[0]:loc[1]]
		tail := fmt[loc[1]:]

		if head != "" {
			parsed = append(parsed, head)
		}
		parsed = append(parsed, tag)
		fmt = tail
	}
	return parsed
}
