package stdlib

import "bytes"

//var magicCharacters = "().%+-*?[]^$"
var characterClasses = map[byte]string{
	'a': "[[:alpha:]]",   // letters
	'A': "[[:^alpha:]]",  //
	'c': "[[:cntrl:]]",   // control characters
	'C': "[[:^cntrl:]]",  //
	'd': "[[:digit:]]",   // digits
	'D': "[[:^digit:]]",  //
	'g': "[[:print:]]",   // printable characters except spaces
	'G': "[[:^print:]]",  //
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

// todo: optimize
func toRegexp(pattern string) (re, err string) {
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
