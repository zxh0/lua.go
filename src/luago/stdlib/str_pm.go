package stdlib

import "bytes"
import "regexp"
import "strings"

//var magicCharacters = "().%+-*?[]^$"
var characterClasses = map[byte]string{
	'a': "[[:alpha:]]",   // letters
	'A': "[[:^alpha:]]",  //
	'c': "[[:cntrl:]]",   // control characters
	'C': "[[:^cntrl:]]",  //
	'd': "[[:digit:]]",   // digits
	'D': "[[:^digit:]]",  //
	'g': "[[:graph:]]",   // printable characters except spaces
	'G': "[[:^graph:]]",  //
	'l': "[[:lower:]]",   // lower-case letters
	'L': "[[:^lower:]]",  //
	'p': "[[:punct:]]",   // punctuation characters
	'P': "[[:^punct:]]",  //
	's': "[[:space:]]",   // space characters
	'S': "[[:^space:]]",  //
	'u': "[[:upper:]]",   // upper-case letters
	'U': "[[:^upper:]]",  //
	'w': "[[:word:]]",    // alphanumeric characters
	'W': "[[:^word:]]",   //
	'x': "[[:xdigit:]]",  // hexadecimal digits
	'X': "[[:^xdigit:]]", //
}

func find(s, pattern string, init int, plain bool) (start, end int) {
	tail := s
	if init > 1 {
		tail = s[init-1:]
	}

	if plain {
		start = strings.Index(tail, pattern)
		end = start + len(pattern) - 1
	} else {
		re, err := _compile(pattern)
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
	if init > 1 {
		tail = s[init-1:]
	}

	re, err := _compile(pattern)
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
	re, err := _compile(pattern)
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

func _compile(pattern string) (*regexp.Regexp, string) {
	expr, errStr := _toRegexp(pattern)
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

// todo: optimize
func _toRegexp(pattern string) (re, err string) {
	var buf bytes.Buffer

	inBrackets := false
	for len(pattern) > 0 {
		var b0, b1 byte
		b0, pattern = pattern[0], pattern[1:]

		switch b0 {
		case '%':
			if len(pattern) == 0 {
				return "", "malformed pattern (ends with '%')"
			}

			b1, pattern = pattern[0], pattern[1:]
			switch b1 {
			case 'a', 'c', 'd', 'g', 'l', 'p', 's', 'u', 'w', 'x',
				'A', 'C', 'D', 'G', 'L', 'P', 'S', 'U', 'W', 'X':
				cc := characterClasses[b1]
				if inBrackets {
					buf.WriteString(cc[1 : len(cc)-1])
				} else {
					buf.WriteString(cc)
				}
			case '(', ')', '[', ']', '.', '?', '*', '+', '-', '^', '$':
				buf.WriteByte('\\')
				buf.WriteByte(b1)
			case '%':
				buf.WriteByte('%')
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				return "", "backreference is not supportted!"
			case 'b': // balanced string
				return "", "'%b' is not supportted!"
			case 'f': // lookahead
				return "", "'%f' is not supportted!"
			default:
				buf.WriteByte(b1)
			}
		case '{', '}', '\\':
			buf.WriteByte('\\')
			buf.WriteByte(b0)
		case '[':
			inBrackets = true
			buf.WriteByte(b0)
		case ']':
			inBrackets = false
			buf.WriteByte(b0)
		case '-':
			if inBrackets {
				buf.WriteByte(b0)
			} else {
				buf.WriteString("*?")
			}
		default:
			buf.WriteByte(b0)
		}
	}

	return buf.String(), ""
}
