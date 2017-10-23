package lexer

import "bytes"
import "fmt"
import "regexp"
import "strconv"
import "strings"

//var reSpaces = regexp.MustCompile(`^\s+`)
var reIdentifier = regexp.MustCompile(`^[_\d\w]+`)
var reNumber = regexp.MustCompile(`^-?0[xX][0-9a-fA-F]+(\.[0-9a-fA-F]+)?([pP][+\-]?[0-9]+)?|^-?[0-9]+(\.[0-9]+)?([eE][+\-]?[0-9]+)?`)
var reShortStr = regexp.MustCompile(`(?s)(^'(\\'|[^'])*')|(^"(\\"|[^"])*")`)
var reLongStrStart = regexp.MustCompile(`^\[=*\[`)

var reDecEscapeSeq = regexp.MustCompile(`^\\[0-9]{1,3}`)
var reHexEscapeSeq = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)
var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

type Lexer struct {
	source string
	chunk  string
	line   int // current line number
}

func NewLexer(source string, chunk string) *Lexer {
	return &Lexer{source, chunk, 1}
}

func (self *Lexer) Line() int {
	return self.line
}

func (self *Lexer) Backup() Lexer {
	return Lexer{"", self.chunk, self.line}
}
func (self *Lexer) Restore(backup Lexer) {
	self.chunk = backup.chunk
	self.line = backup.line
}

// todo
func (self *Lexer) LookAhead(n int) int {
	backup := self.Backup()
	for i := 1; i < n; i++ {
		self.NextToken()
	}
	_, kind, _ := self.NextToken()
	self.Restore(backup)
	return kind
}

func (self *Lexer) NextIdentifier() (line int, token string) {
	return self.NextTokenOfKind(TOKEN_IDENTIFIER)
}

func (self *Lexer) NextTokenOfKind(kind int) (line int, token string) {
	line, _kind, token := self.NextToken()
	if kind != _kind {
		self.error("syntax error near '%s'", token)
	}
	return line, token
}

func (self *Lexer) NextToken() (line, kind int, token string) {
	self.skipWhiteSpaces()
	line = self.line

	if len(self.chunk) == 0 {
		return line, TOKEN_EOF, "EOF"
	}

	switch self.chunk[0] {
	case ';':
		self.next(1)
		return line, TOKEN_SEP_SEMI, ""
	case ',':
		self.next(1)
		return line, TOKEN_SEP_COMMA, ""
	case '(':
		self.next(1)
		return line, TOKEN_SEP_LPAREN, ""
	case ')':
		self.next(1)
		return line, TOKEN_SEP_RPAREN, ""
	case ']':
		self.next(1)
		return line, TOKEN_SEP_RBRACK, ""
	case '{':
		self.next(1)
		return line, TOKEN_SEP_LCURLY, ""
	case '}':
		self.next(1)
		return line, TOKEN_SEP_RCURLY, ""
	case '+':
		self.next(1)
		return line, TOKEN_OP_ADD, ""
	case '-':
		self.next(1)
		return line, TOKEN_MINUS, ""
	case '*':
		self.next(1)
		return line, TOKEN_OP_MUL, ""
	case '^':
		self.next(1)
		return line, TOKEN_OP_POW, ""
	case '%':
		self.next(1)
		return line, TOKEN_OP_MOD, ""
	case '&':
		self.next(1)
		return line, TOKEN_OP_BAND, ""
	case '|':
		self.next(1)
		return line, TOKEN_OP_BOR, ""
	case '#':
		self.next(1)
		return line, TOKEN_OP_LEN, ""
	case '.':
		if self.test("...") {
			self.next(3)
			return line, TOKEN_VARARG, ""
		} else if self.test("..") {
			self.next(2)
			return line, TOKEN_OP_CONCAT, ""
		} else {
			self.next(1)
			return line, TOKEN_SEP_DOT, ""
		}
	case ':':
		if self.test("::") {
			self.next(2)
			return line, TOKEN_SEP_LABEL, ""
		} else {
			self.next(1)
			return line, TOKEN_SEP_COLON, ""
		}
	case '/':
		if self.test("//") {
			self.next(2)
			return line, TOKEN_OP_IDIV, ""
		} else {
			self.next(1)
			return line, TOKEN_OP_DIV, ""
		}
	case '~':
		if self.test("~=") {
			self.next(2)
			return line, TOKEN_OP_NE, ""
		} else {
			self.next(1)
			return line, TOKEN_WAVE, ""
		}
	case '=':
		if self.test("==") {
			self.next(2)
			return line, TOKEN_OP_EQ, ""
		} else {
			self.next(1)
			return line, TOKEN_ASSIGN, ""
		}
	case '<':
		if self.test("<<") {
			self.next(2)
			return line, TOKEN_OP_SHL, ""
		} else if self.test("<=") {
			self.next(2)
			return line, TOKEN_OP_LE, ""
		} else {
			self.next(1)
			return line, TOKEN_OP_LT, ""
		}
	case '>':
		if self.test(">>") {
			self.next(2)
			return line, TOKEN_OP_SHR, ""
		} else if self.test(">=") {
			self.next(2)
			return line, TOKEN_OP_GE, ""
		} else {
			self.next(1)
			return line, TOKEN_OP_GT, ""
		}
	case '[':
		if self.test("[[") || self.test("[=") {
			return line, TOKEN_STRING, self.scanLongString()
		} else {
			self.next(1)
			return line, TOKEN_SEP_LBRACK, ""
		}
	case '\'', '"':
		return line, TOKEN_STRING, self.scanShortString()
	}

	c := self.chunk[0]
	if c == '_' || isLatter(c) {
		token := self.scanIdentifier()
		if kind, found := keywords[token]; found {
			return line, kind, "" // keyword
		} else {
			return line, TOKEN_IDENTIFIER, token
		}
	} else if isDigit(c) {
		token := self.scanNumber()
		return line, TOKEN_NUMBER, token
	}

	self.error("unexpected symbol near %q", c)
	return
}

func (self *Lexer) next(n int) {
	self.chunk = self.chunk[n:]
}

func (self *Lexer) test(s string) bool {
	return strings.HasPrefix(self.chunk, s)
}

func (self *Lexer) error(f string, a ...interface{}) {
	err := fmt.Sprintf(f, a...)
	err = fmt.Sprintf("%s:%d: %s", self.source, self.line, err)
	panic(err)
}

func (self *Lexer) skipWhiteSpaces() {
	for len(self.chunk) > 0 {
		if self.test("--") {
			self.skipComment()
		} else if c := self.chunk[0]; isSpace(c) {
			if self.test("\r\n") {
				self.next(2)
				self.line += 1
			} else {
				self.next(1)
				if c == '\r' || c == '\n' {
					self.line += 1
				}
			}
		} else {
			break
		}
	}
}

func (self *Lexer) skipComment() {
	self.next(2) // skip --
	if len(self.chunk) == 0 {
		return
	}

	// long comment
	if self.chunk[0] == '[' {
		if reLongStrStart.FindString(self.chunk) != "" {
			self.scanLongString() // todo
			return
		}
	}

	// short comment
	for len(self.chunk) > 0 &&
		self.chunk[0] != '\n' && self.chunk[0] != '\r' {

		self.next(1)
	}
}

func (self *Lexer) scanIdentifier() string {
	return self.scan(reIdentifier)
}

func (self *Lexer) scanNumber() string {
	return self.scan(reNumber)
}

func (self *Lexer) scan(re *regexp.Regexp) string {
	if token := re.FindString(self.chunk); token != "" {
		self.next(len(token))
		return token
	}
	panic("unreachable!")
}

func (self *Lexer) scanLongString() string {
	startStr := reLongStrStart.FindString(self.chunk)
	if startStr == "" {
		self.error("invalid long string delimiter near '%s'",
			self.chunk[0:2])
	}

	endStr := strings.Replace(startStr, "[", "]", -1)
	endIdx := strings.Index(self.chunk, endStr)
	if endIdx < 0 {
		self.error("unfinished long string or comment")
	}

	str := self.chunk[len(startStr):endIdx]
	self.next(endIdx + len(endStr))
	self.line += strings.Count(str, "\n")
	return str
}

func (self *Lexer) scanShortString() string {
	if str := reShortStr.FindString(self.chunk); str != "" {
		self.next(len(str))
		str = str[1 : len(str)-1]
		if strings.Index(str, `\`) >= 0 {
			str = self.escape(str)
		}
		return str
	}
	self.error("unfinished string")
	return ""
}

func (self *Lexer) escape(str string) string {
	var buf bytes.Buffer

	for len(str) > 0 {
		if str[0] != '\\' {
			buf.WriteByte(str[0])
			str = str[1:]
		} else if len(str) > 1 {
			switch str[1] {
			case 'a':
				buf.WriteByte('\a')
				str = str[2:]
				continue
			case 'b':
				buf.WriteByte('\b')
				str = str[2:]
				continue
			case 'f':
				buf.WriteByte('\f')
				str = str[2:]
				continue
			case 'n':
				buf.WriteByte('\n')
				str = str[2:]
				continue
			case 'r':
				buf.WriteByte('\r')
				str = str[2:]
				continue
			case 't':
				buf.WriteByte('\t')
				str = str[2:]
				continue
			case 'v':
				buf.WriteByte('\v')
				str = str[2:]
				continue
			case '"':
				buf.WriteByte('"')
				str = str[2:]
				continue
			case '\'':
				buf.WriteByte('\'')
				str = str[2:]
				continue
			case '\\':
				buf.WriteByte('\\')
				str = str[2:]
				continue
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // \ddd
				if found := reDecEscapeSeq.FindString(str); found != "" {
					d, _ := strconv.ParseInt(found[1:], 10, 32)
					if d <= 0xFF {
						buf.WriteByte(byte(d))
						str = str[len(found):]
						continue
					}
					self.error("decimal escape too large near '%s'", found)
				}
			case 'x': // \xXX
				if found := reHexEscapeSeq.FindString(str); found != "" {
					d, _ := strconv.ParseInt(found[2:], 16, 32)
					buf.WriteByte(byte(d))
					str = str[len(found):]
					continue
				}
			case 'u': // \u{XXX}
				if found := reUnicodeEscapeSeq.FindString(str); found != "" {
					if len(found) <= 10 {
						d, _ := strconv.ParseInt(found[3:len(found)-1], 16, 32)
						if d <= 0x10FFFF {
							buf.WriteRune(rune(d))
							str = str[len(found):]
							continue
						}
					}
					self.error("UTF-8 value too large near '%s'", found)
				}
			case 'z':
				str = str[2:]
				for len(str) > 0 && isSpace(str[0]) {
					if str[0] == '\n' {
						self.line += 1
					}
					str = str[1:]
				}
				continue
			}
			self.error("invalid escape sequence near '\\%c'", str[1])
		}
	}

	return buf.String()
}

func isSpace(x byte) bool {
	switch x {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	default:
		return false
	}
}

func isDigit(x byte) bool {
	return x >= '0' && x <= '9'
}

func isLatter(x byte) bool {
	return x >= 'a' && x <= 'z' || x >= 'A' && x <= 'Z'
}
