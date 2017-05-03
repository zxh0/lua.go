package lexer

import "fmt"
import "regexp"
import "strconv"
import "strings"
import "unicode"
import "unicode/utf8"

//var reSpaces = regexp.MustCompile(`^\s+`)
var reIdentifier = regexp.MustCompile(`^[_\d\w]+`)
var reNumber = regexp.MustCompile(`^-?0[xX][0-9a-fA-F]+(\.[0-9a-fA-F]+)?([pP][+\-]?[0-9]+)?|^-?[0-9]+(\.[0-9]+)?([eE][+\-]?[0-9]+)?`)
var reShortStr = regexp.MustCompile(`(?s)(^'(\\'|[^'])*')|(^"(\\"|[^"])*")`)
var reLongStringStart = regexp.MustCompile(`^\[=*\[`)

var reEscapeDecimalDigits = regexp.MustCompile(`^\\[0-9]{1,3}`)
var reEscapeHexDigits = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)
var reEscapeUnicode = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

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
	if kind == _kind {
		return line, token
	}

	panic(fmt.Sprintf("%s:%d: syntax error near %q",
		self.source, line, token))
}

func (self *Lexer) NextToken() (line, kind int, token string) {
	self.skipWhiteSpaces()
	line = self.line

	remains := len(self.chunk)
	if remains == 0 {
		return line, TOKEN_EOF, "EOF"
	}

	b0 := self.chunk[0]

	switch b0 {
	case ';':
		self.chunk = self.chunk[1:]
		return line, TOKEN_SEP_SEMI, ""
	case ',':
		self.chunk = self.chunk[1:]
		return line, TOKEN_SEP_COMMA, ""
		self.chunk = self.chunk[1:]
	case '(':
		self.chunk = self.chunk[1:]
		return line, TOKEN_SEP_LPAREN, ""
	case ')':
		self.chunk = self.chunk[1:]
		return line, TOKEN_SEP_RPAREN, ""
	case ']':
		self.chunk = self.chunk[1:]
		return line, TOKEN_SEP_RBRACK, ""
	case '{':
		self.chunk = self.chunk[1:]
		return line, TOKEN_SEP_LCURLY, ""
	case '}':
		self.chunk = self.chunk[1:]
		return line, TOKEN_SEP_RCURLY, ""
	case '+':
		self.chunk = self.chunk[1:]
		return line, TOKEN_OP_ADD, ""
	case '-':
		self.chunk = self.chunk[1:]
		return line, TOKEN_MINUS, ""
	case '*':
		self.chunk = self.chunk[1:]
		return line, TOKEN_OP_MUL, ""
	case '^':
		self.chunk = self.chunk[1:]
		return line, TOKEN_OP_POW, ""
	case '%':
		self.chunk = self.chunk[1:]
		return line, TOKEN_OP_MOD, ""
	case '&':
		self.chunk = self.chunk[1:]
		return line, TOKEN_OP_BAND, ""
	case '|':
		self.chunk = self.chunk[1:]
		return line, TOKEN_OP_BOR, ""
	case '#':
		self.chunk = self.chunk[1:]
		return line, TOKEN_OP_LEN, ""
	case '.':
		if remains > 2 && self.chunk[1] == '.' && self.chunk[2] == '.' {
			self.chunk = self.chunk[3:]
			return line, TOKEN_VARARG, ""
		} else if remains > 1 && self.chunk[1] == '.' {
			self.chunk = self.chunk[2:]
			return line, TOKEN_OP_CONCAT, ""
		} else {
			self.chunk = self.chunk[1:]
			return line, TOKEN_SEP_DOT, ""
		}
	case ':':
		if remains > 1 && self.chunk[1] == ':' {
			self.chunk = self.chunk[2:]
			return line, TOKEN_SEP_LABEL, ""
		} else {
			self.chunk = self.chunk[1:]
			return line, TOKEN_SEP_COLON, ""
		}
	case '/':
		if remains > 1 && self.chunk[1] == '/' {
			self.chunk = self.chunk[2:]
			return line, TOKEN_OP_IDIV, ""
		} else {
			self.chunk = self.chunk[1:]
			return line, TOKEN_OP_DIV, ""
		}
	case '~':
		if remains > 1 && self.chunk[1] == '=' {
			self.chunk = self.chunk[2:]
			return line, TOKEN_OP_NE, ""
		} else {
			self.chunk = self.chunk[1:]
			return line, TOKEN_WAVE, ""
		}
	case '=':
		if remains > 1 && self.chunk[1] == '=' {
			self.chunk = self.chunk[2:]
			return line, TOKEN_OP_EQ, ""
		} else {
			self.chunk = self.chunk[1:]
			return line, TOKEN_ASSIGN, ""
		}
	case '<':
		if remains > 1 && self.chunk[1] == '<' {
			self.chunk = self.chunk[2:]
			return line, TOKEN_OP_SHL, ""
		} else if remains > 1 && self.chunk[1] == '=' {
			self.chunk = self.chunk[2:]
			return line, TOKEN_OP_LE, ""
		} else {
			self.chunk = self.chunk[1:]
			return line, TOKEN_OP_LT, ""
		}
	case '>':
		if remains > 1 && self.chunk[1] == '>' {
			self.chunk = self.chunk[2:]
			return line, TOKEN_OP_SHR, ""
		} else if remains > 1 && self.chunk[1] == '=' {
			self.chunk = self.chunk[2:]
			return line, TOKEN_OP_GE, ""
		} else {
			self.chunk = self.chunk[1:]
			return line, TOKEN_OP_GT, ""
		}
	case '[':
		if remains > 1 && self.chunk[1] == '[' || self.chunk[1] == '=' {
			return line, TOKEN_STRING, self.scanLongString()
		} else {
			self.chunk = self.chunk[1:]
			return line, TOKEN_SEP_LBRACK, ""
		}
	}

	if b0 == '_' || isLatter(b0) {
		token := self.scanIdentifier()
		if kind, found := keywords[token]; found {
			return line, kind, "" // keyword
		} else {
			return line, TOKEN_IDENTIFIER, token
		}
	} else if isDigit(b0) {
		token := self.scanNumber()
		return line, TOKEN_NUMBER, token
	} else if b0 == '\'' || b0 == '"' {
		if token, ok := self.scanShortString(); ok {
			return line, TOKEN_STRING, token
		}
	}

	var msg string
	_, size := utf8.DecodeRuneInString(self.chunk)
	if size > 0 {
		msg = self.chunk[0:size]
	} else {
		msg = self.chunk[0:1]
	}
	panic(fmt.Sprintf("%s:%d: unexpected symbol near %q!",
		self.source, self.line, msg))
}

func (self *Lexer) skipWhiteSpaces() {
	for len(self.chunk) > 0 {
		if len(self.chunk) > 1 && self.chunk[0] == '-' && self.chunk[1] == '-' {
			self.skipComment()
		} else {
			r, size := utf8.DecodeRuneInString(self.chunk)
			if unicode.IsSpace(r) {
				self.chunk = self.chunk[size:]
				if r == '\n' {
					self.line += 1
				}
			} else {
				break
			}
		}
	}
}

func (self *Lexer) skipComment() {
	self.chunk = self.chunk[2:]

	if len(self.chunk) == 0 {
		return
	}
	if self.chunk[0] == '[' {
		if reLongStringStart.FindString(self.chunk) != "" {
			self.scanLongString() // todo
			return
		}
	}

	if idxOfNL := strings.IndexByte(self.chunk, '\n'); idxOfNL >= 0 {
		self.chunk = self.chunk[idxOfNL+1:]
		self.line += 1
	} else {
		self.chunk = ""
	}
}

func (self *Lexer) scanIdentifier() string {
	return self.scan(reIdentifier)
}

func (self *Lexer) scanNumber() string {
	return self.scan(reNumber)
}

func (self *Lexer) scan(re *regexp.Regexp) string {
	token := re.FindString(self.chunk)
	if token != "" {
		self.chunk = self.chunk[len(token):]
		return token
	}
	panic("unreachable!")
}

func (self *Lexer) scanLongString() string {
	startStr := reLongStringStart.FindString(self.chunk)
	if startStr == "" {
		panic(fmt.Sprintf("%s:%d: invalid long string delimiter",
			self.source, self.line))
	}

	endStr := strings.Replace(startStr, "[", "]", -1)
	endIdx := strings.Index(self.chunk, endStr)
	if endIdx < 0 {
		panic(fmt.Sprintf("%s:%d: unfinished long string or comment",
			self.source, self.line))
	}

	str := self.chunk[len(startStr):endIdx]
	self.chunk = self.chunk[endIdx+len(endStr):]
	self.line += strings.Count(str, "\n")
	return str
}

func (self *Lexer) scanShortString() (string, bool) {
	literal := reShortStr.FindString(self.chunk)
	if literal == "" {
		return "", false
	}

	self.chunk = self.chunk[len(literal):]
	literal = literal[1 : len(literal)-1]

	if strings.Index(literal, `\`) < 0 {
		return literal, true
	}

	str := self.escape(literal)
	return str, true
}

func (self *Lexer) escape(literal string) string {
	buf := make([]rune, 0, len(literal))

	for len(literal) > 0 {
		if literal[0] != '\\' {
			r, size := utf8.DecodeRuneInString(literal)
			literal = literal[size:]
			buf = append(buf, r)
		} else if len(literal) > 1 {
			switch literal[1] {
			case 'a':
				buf = append(buf, '\a')
				literal = literal[2:]
				continue
			case 'b':
				buf = append(buf, '\b')
				literal = literal[2:]
				continue
			case 'f':
				buf = append(buf, '\f')
				literal = literal[2:]
				continue
			case 'n':
				buf = append(buf, '\n')
				literal = literal[2:]
				continue
			case 'r':
				buf = append(buf, '\r')
				literal = literal[2:]
				continue
			case 't':
				buf = append(buf, '\t')
				literal = literal[2:]
				continue
			case 'v':
				buf = append(buf, '\v')
				literal = literal[2:]
				continue
			case '"':
				buf = append(buf, '"')
				literal = literal[2:]
				continue
			case '\'':
				buf = append(buf, '\'')
				literal = literal[2:]
				continue
			case '\\':
				buf = append(buf, '\\')
				literal = literal[2:]
				continue
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // \ddd
				if found := reEscapeDecimalDigits.FindString(literal); found != "" {
					if s, err := strconv.ParseInt(found[1:], 10, 32); err == nil {
						buf = append(buf, rune(s))
						literal = literal[len(found):]
						continue
					}
				}
			case 'x': // \xXX
				if found := reEscapeHexDigits.FindString(literal); found != "" {
					if s, err := strconv.ParseInt(found[2:], 16, 32); err == nil {
						buf = append(buf, rune(s))
						literal = literal[len(found):]
						continue
					}
				}
			case 'u': // \u{XXX}
				if found := reEscapeUnicode.FindString(literal); found != "" {
					if s, err := strconv.ParseInt(found[3:len(found)-1], 16, 32); err == nil {
						buf = append(buf, rune(s))
						literal = literal[len(found):]
						continue
					}
				}
			case 'z':
				literal = literal[2:]
				for len(literal) > 0 && isSpace(literal[0]) {
					if literal[0] == '\n' {
						self.line += 1
					}
					literal = literal[1:]
				}
				continue
			}
			panic(fmt.Sprintf("%s:%d: invalid escape sequence near '\\%c'!",
				self.source, self.line, literal[1]))
		} else {
			panic("unreachable!")
		}
	}

	return string(buf)
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
