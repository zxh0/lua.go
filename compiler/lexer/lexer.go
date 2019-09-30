package lexer

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//var reSpaces = regexp.MustCompile(`^\s+`)
var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")
var reIdentifier = regexp.MustCompile(`^[_\d\w]+`)
var reNumber = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?|^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)
var reShortStr = regexp.MustCompile(`(?s)(^'(\\\\|\\'|\\\n|\\z\s*|[^'\n])*')|(^"(\\\\|\\"|\\\n|\\z\s*|[^"\n])*")`)
var reOpeningLongBracket = regexp.MustCompile(`^\[=*\[`)

var reDecEscapeSeq = regexp.MustCompile(`^\\[0-9]{1,3}`)
var reHexEscapeSeq = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)
var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

type Lexer struct {
	chunk         string // source code
	chunkName     string // source name
	line          int    // current line number
	nextToken     string
	nextTokenKind int
	nextTokenLine int
}

func NewLexer(chunk, chunkName string) *Lexer {
	return &Lexer{chunk, chunkName, 1, "", 0, 0}
}

func (lexer *Lexer) Line() int {
	return lexer.line
}

func (lexer *Lexer) LookAhead() int {
	if lexer.nextTokenLine > 0 {
		return lexer.nextTokenKind
	}
	currentLine := lexer.line
	line, kind, token := lexer.NextToken()
	lexer.line = currentLine
	lexer.nextTokenLine = line
	lexer.nextTokenKind = kind
	lexer.nextToken = token
	return kind
}

func (lexer *Lexer) NextIdentifier() (line int, token string) {
	return lexer.NextTokenOfKind(TOKEN_IDENTIFIER)
}

func (lexer *Lexer) NextTokenOfKind(kind int) (line int, token string) {
	line, _kind, token := lexer.NextToken()
	if kind != _kind {
		lexer.error("syntax error near '%s'", token)
	}
	return line, token
}

func (lexer *Lexer) NextToken() (line, kind int, token string) {
	if lexer.nextTokenLine > 0 {
		line = lexer.nextTokenLine
		kind = lexer.nextTokenKind
		token = lexer.nextToken
		lexer.line = lexer.nextTokenLine
		lexer.nextTokenLine = 0
		return
	}

	lexer.skipWhiteSpaces()
	if len(lexer.chunk) == 0 {
		return lexer.line, TOKEN_EOF, "EOF"
	}

	switch lexer.chunk[0] {
	case ';':
		lexer.next(1)
		return lexer.line, TOKEN_SEP_SEMI, ";"
	case ',':
		lexer.next(1)
		return lexer.line, TOKEN_SEP_COMMA, ","
	case '(':
		lexer.next(1)
		return lexer.line, TOKEN_SEP_LPAREN, "("
	case ')':
		lexer.next(1)
		return lexer.line, TOKEN_SEP_RPAREN, ")"
	case ']':
		lexer.next(1)
		return lexer.line, TOKEN_SEP_RBRACK, "]"
	case '{':
		lexer.next(1)
		return lexer.line, TOKEN_SEP_LCURLY, "{"
	case '}':
		lexer.next(1)
		return lexer.line, TOKEN_SEP_RCURLY, "}"
	case '+':
		lexer.next(1)
		return lexer.line, TOKEN_OP_ADD, "+"
	case '-':
		lexer.next(1)
		return lexer.line, TOKEN_OP_MINUS, "-"
	case '*':
		lexer.next(1)
		return lexer.line, TOKEN_OP_MUL, "*"
	case '^':
		lexer.next(1)
		return lexer.line, TOKEN_OP_POW, "^"
	case '%':
		lexer.next(1)
		return lexer.line, TOKEN_OP_MOD, "%"
	case '&':
		lexer.next(1)
		return lexer.line, TOKEN_OP_BAND, "&"
	case '|':
		lexer.next(1)
		return lexer.line, TOKEN_OP_BOR, "|"
	case '#':
		lexer.next(1)
		return lexer.line, TOKEN_OP_LEN, "#"
	case ':':
		if lexer.test("::") {
			lexer.next(2)
			return lexer.line, TOKEN_SEP_LABEL, "::"
		} else {
			lexer.next(1)
			return lexer.line, TOKEN_SEP_COLON, ":"
		}
	case '/':
		if lexer.test("//") {
			lexer.next(2)
			return lexer.line, TOKEN_OP_IDIV, "//"
		} else {
			lexer.next(1)
			return lexer.line, TOKEN_OP_DIV, "/"
		}
	case '~':
		if lexer.test("~=") {
			lexer.next(2)
			return lexer.line, TOKEN_OP_NE, "~="
		} else {
			lexer.next(1)
			return lexer.line, TOKEN_OP_WAVE, "~"
		}
	case '=':
		if lexer.test("==") {
			lexer.next(2)
			return lexer.line, TOKEN_OP_EQ, "=="
		} else {
			lexer.next(1)
			return lexer.line, TOKEN_OP_ASSIGN, "="
		}
	case '<':
		if lexer.test("<<") {
			lexer.next(2)
			return lexer.line, TOKEN_OP_SHL, "<<"
		} else if lexer.test("<=") {
			lexer.next(2)
			return lexer.line, TOKEN_OP_LE, "<="
		} else {
			lexer.next(1)
			return lexer.line, TOKEN_OP_LT, "<"
		}
	case '>':
		if lexer.test(">>") {
			lexer.next(2)
			return lexer.line, TOKEN_OP_SHR, ">>"
		} else if lexer.test(">=") {
			lexer.next(2)
			return lexer.line, TOKEN_OP_GE, ">="
		} else {
			lexer.next(1)
			return lexer.line, TOKEN_OP_GT, ">"
		}
	case '.':
		if lexer.test("...") {
			lexer.next(3)
			return lexer.line, TOKEN_VARARG, "..."
		} else if lexer.test("..") {
			lexer.next(2)
			return lexer.line, TOKEN_OP_CONCAT, ".."
		} else if len(lexer.chunk) == 1 || !isDigit(lexer.chunk[1]) {
			lexer.next(1)
			return lexer.line, TOKEN_SEP_DOT, "."
		}
	case '[':
		if lexer.test("[[") || lexer.test("[=") {
			return lexer.line, TOKEN_STRING, lexer.scanLongString()
		} else {
			lexer.next(1)
			return lexer.line, TOKEN_SEP_LBRACK, "["
		}
	case '\'', '"':
		return lexer.line, TOKEN_STRING, lexer.scanShortString()
	}

	c := lexer.chunk[0]
	if c == '.' || isDigit(c) {
		token := lexer.scanNumber()
		return lexer.line, TOKEN_NUMBER, token
	}
	if c == '_' || isLetter(c) {
		token := lexer.scanIdentifier()
		if kind, found := keywords[token]; found {
			return lexer.line, kind, token // keyword
		} else {
			return lexer.line, TOKEN_IDENTIFIER, token
		}
	}

	lexer.error("unexpected symbol near %q", c)
	return
}

func (lexer *Lexer) next(n int) {
	lexer.chunk = lexer.chunk[n:]
}

func (lexer *Lexer) test(s string) bool {
	return strings.HasPrefix(lexer.chunk, s)
}

func (lexer *Lexer) error(f string, a ...interface{}) {
	err := fmt.Sprintf(f, a...)
	err = fmt.Sprintf("%s:%d: %s", lexer.chunkName, lexer.line, err)
	panic(err)
}

func (lexer *Lexer) skipWhiteSpaces() {
	for len(lexer.chunk) > 0 {
		if lexer.test("--") {
			lexer.skipComment()
		} else if lexer.test("\r\n") || lexer.test("\n\r") {
			lexer.next(2)
			lexer.line += 1
		} else if isNewLine(lexer.chunk[0]) {
			lexer.next(1)
			lexer.line += 1
		} else if isWhiteSpace(lexer.chunk[0]) {
			lexer.next(1)
		} else {
			break
		}
	}
}

func (lexer *Lexer) skipComment() {
	lexer.next(2) // skip --

	// long comment ?
	if lexer.test("[") {
		if reOpeningLongBracket.FindString(lexer.chunk) != "" {
			lexer.scanLongString()
			return
		}
	}

	// short comment
	for len(lexer.chunk) > 0 && !isNewLine(lexer.chunk[0]) {
		lexer.next(1)
	}
}

func (lexer *Lexer) scanIdentifier() string {
	return lexer.scan(reIdentifier)
}

func (lexer *Lexer) scanNumber() string {
	return lexer.scan(reNumber)
}

func (lexer *Lexer) scan(re *regexp.Regexp) string {
	if token := re.FindString(lexer.chunk); token != "" {
		lexer.next(len(token))
		return token
	}
	panic("unreachable!")
}

func (lexer *Lexer) scanLongString() string {
	openingLongBracket := reOpeningLongBracket.FindString(lexer.chunk)
	if openingLongBracket == "" {
		lexer.error("invalid long string delimiter near '%s'",
			lexer.chunk[0:2])
	}

	closingLongBracket := strings.Replace(openingLongBracket, "[", "]", -1)
	closingLongBracketIdx := strings.Index(lexer.chunk, closingLongBracket)
	if closingLongBracketIdx < 0 {
		lexer.error("unfinished long string or comment")
	}

	str := lexer.chunk[len(openingLongBracket):closingLongBracketIdx]
	lexer.next(closingLongBracketIdx + len(closingLongBracket))

	str = reNewLine.ReplaceAllString(str, "\n")
	lexer.line += strings.Count(str, "\n")
	if len(str) > 0 && str[0] == '\n' {
		str = str[1:]
	}

	return str
}

func (lexer *Lexer) scanShortString() string {
	if str := reShortStr.FindString(lexer.chunk); str != "" {
		lexer.next(len(str))
		str = str[1 : len(str)-1]
		if strings.Index(str, `\`) >= 0 {
			lexer.line += len(reNewLine.FindAllString(str, -1))
			str = lexer.escape(str)
		}
		return str
	}
	lexer.error("unfinished string")
	return ""
}

func (lexer *Lexer) escape(str string) string {
	var buf bytes.Buffer

	for len(str) > 0 {
		if str[0] != '\\' {
			buf.WriteByte(str[0])
			str = str[1:]
			continue
		}

		if len(str) == 1 {
			lexer.error("unfinished string")
		}

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
		case 'n', '\n':
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
				lexer.error("decimal escape too large near '%s'", found)
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
				d, err := strconv.ParseInt(found[3:len(found)-1], 16, 32)
				if err == nil && d <= 0x10FFFF {
					buf.WriteRune(rune(d))
					str = str[len(found):]
					continue
				}
				lexer.error("UTF-8 value too large near '%s'", found)
			}
		case 'z':
			str = str[2:]
			for len(str) > 0 && isWhiteSpace(str[0]) { // todo
				str = str[1:]
			}
			continue
		}
		lexer.error("invalid escape sequence near '\\%c'", str[1])
	}

	return buf.String()
}

func isWhiteSpace(c byte) bool {
	switch c {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}

func isNewLine(c byte) bool {
	return c == '\r' || c == '\n'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}
