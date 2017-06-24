package stdlib

import "regexp"
import "strings"

// [max(1, i), min(len(s), j)]
func subStr(s string, i, j int) string {
	if i < 0 {
		i = len(s) + i + 1
	}
	if j < 0 {
		j = len(s) + j + 1
	}

	i = max(i, 1)
	j = min(j, len(s))

	if i > j {
		return ""
	}
	return s[i-1 : j]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parseFmtStr(fmt string) []string {
	if fmt == "" || strings.IndexByte(fmt, '%') < 0 {
		return []string{fmt}
	}

	// tag = %[flags][width][.precision]specifier
	pattern := regexp.MustCompile(`%[ #+-0]?[0-9]*(\.[0-9]+)?[cdeEfgGioqsuxX%]`)

	parsed := make([]string, 0, len(fmt)/2)
	for {
		if fmt == "" {
			break
		}

		loc := pattern.FindStringIndex(fmt)
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

func quote(s string) string {
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, `"`, `\"`, -1)
	s = strings.Replace(s, "\n", "\\\n", -1)
	s = strings.Replace(s, "\x00", "\\0", -1)
	return "\"" + s + "\""
}

func find(s, pattern string, init int, plain bool) (start, end int) {
	tail := s
	if init != 1 && init != 0 {
		tail = subStr(s, init, -1)
	}

	if plain {
		start = strings.Index(tail, pattern)
		end = start + len(pattern) - 1
	} else {
		re, err := compile(pattern)
		if err != "" {
			panic(err) // todo
		} else {
			loc := re.FindStringIndex(tail)
			if loc == nil {
				start, end = -1, -1
			} else {
				start, end = loc[0], loc[1]-1
			}
		}
	}
	if start >= 0 {
		start += len(s) - len(tail) + 1
		end += len(s) - len(tail) + 1
	}

	return
}

func match(s, pattern string, init int) []int {
	tail := s
	if init != 1 && init != 0 {
		tail = subStr(s, init, -1)
	}

	re, err := compile(pattern)
	if err != "" {
		panic(err) // todo
	} else {
		found := re.FindStringSubmatchIndex(tail)
		if len(found) > 2 {
			return found[2:]
		} else {
			return found
		}
	}
}

// todo
func gsub(s, pattern, repl string, n int) (string, int) {
	re, err := compile(pattern)
	if err != "" {
		panic(err) // todo
	} else {
		indexes := re.FindAllStringIndex(s, n)
		if indexes == nil {
			return s, 0
		}

		nMatches := len(indexes)
		lastEnd := indexes[nMatches-1][1]
		head, tail := s[:lastEnd], s[lastEnd:]

		repl = strings.Replace(repl, "%", "$", -1)

		newHead := re.ReplaceAllString(head, repl)
		return newHead + tail, nMatches
	}
}

func compile(pattern string) (*regexp.Regexp, string) {
	expr, errStr := toRegexp(pattern)
	if errStr != "" {
		return nil, errStr
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err.Error() // todo
	} else {
		return re, ""
	}
}
