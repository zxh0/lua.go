package lexer

import "testing"
import "assert"

func TestNextToken(t *testing.T) {
	lexer := NewLexer("str", `;,()[]{}+-*^%%&|#`)
	assertNextTokenKind(t, lexer, TOKEN_SEP_SEMI)
	assertNextTokenKind(t, lexer, TOKEN_SEP_COMMA)
	assertNextTokenKind(t, lexer, TOKEN_SEP_LPAREN)
	assertNextTokenKind(t, lexer, TOKEN_SEP_RPAREN)
	assertNextTokenKind(t, lexer, TOKEN_SEP_LBRACK)
	assertNextTokenKind(t, lexer, TOKEN_SEP_RBRACK)
	assertNextTokenKind(t, lexer, TOKEN_SEP_LCURLY)
	assertNextTokenKind(t, lexer, TOKEN_SEP_RCURLY)
	assertNextTokenKind(t, lexer, TOKEN_OP_ADD)
	assertNextTokenKind(t, lexer, TOKEN_MINUS)
	assertNextTokenKind(t, lexer, TOKEN_OP_MUL)
	assertNextTokenKind(t, lexer, TOKEN_OP_POW)
	assertNextTokenKind(t, lexer, TOKEN_OP_MOD)
	assertNextTokenKind(t, lexer, TOKEN_OP_MOD)
	assertNextTokenKind(t, lexer, TOKEN_OP_BAND)
	assertNextTokenKind(t, lexer, TOKEN_OP_BOR)
	assertNextTokenKind(t, lexer, TOKEN_OP_LEN)
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken2(t *testing.T) {
	lexer := NewLexer("str", `... .. . :: : // / ~= ~ == = << <= < >> >= >`)
	assertNextTokenKind(t, lexer, TOKEN_VARARG)
	assertNextTokenKind(t, lexer, TOKEN_OP_CONCAT)
	assertNextTokenKind(t, lexer, TOKEN_SEP_DOT)
	assertNextTokenKind(t, lexer, TOKEN_SEP_LABEL)
	assertNextTokenKind(t, lexer, TOKEN_SEP_COLON)
	assertNextTokenKind(t, lexer, TOKEN_OP_IDIV)
	assertNextTokenKind(t, lexer, TOKEN_OP_DIV)
	assertNextTokenKind(t, lexer, TOKEN_OP_NE)
	assertNextTokenKind(t, lexer, TOKEN_WAVE)
	assertNextTokenKind(t, lexer, TOKEN_OP_EQ)
	assertNextTokenKind(t, lexer, TOKEN_ASSIGN)
	assertNextTokenKind(t, lexer, TOKEN_OP_SHL)
	assertNextTokenKind(t, lexer, TOKEN_OP_LE)
	assertNextTokenKind(t, lexer, TOKEN_OP_LT)
	assertNextTokenKind(t, lexer, TOKEN_OP_SHR)
	assertNextTokenKind(t, lexer, TOKEN_OP_GE)
	assertNextTokenKind(t, lexer, TOKEN_OP_GT)
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken_keywords(t *testing.T) {
	keywords := `
	and       break     do        else      elseif    end
	false     for       function  goto      if        in
	local     nil       not       or        repeat    return
	then      true      until     while
    `
	lexer := NewLexer("str", keywords)
	assertNextTokenKind(t, lexer, TOKEN_OP_AND)
	assertNextTokenKind(t, lexer, TOKEN_KW_BREAK)
	assertNextTokenKind(t, lexer, TOKEN_KW_DO)
	assertNextTokenKind(t, lexer, TOKEN_KW_ELSE)
	assertNextTokenKind(t, lexer, TOKEN_KW_ELSEIF)
	assertNextTokenKind(t, lexer, TOKEN_KW_END)
	assertNextTokenKind(t, lexer, TOKEN_KW_FALSE)
	assertNextTokenKind(t, lexer, TOKEN_KW_FOR)
	assertNextTokenKind(t, lexer, TOKEN_KW_FUNCTION)
	assertNextTokenKind(t, lexer, TOKEN_KW_GOTO)
	assertNextTokenKind(t, lexer, TOKEN_KW_IF)
	assertNextTokenKind(t, lexer, TOKEN_KW_IN)
	assertNextTokenKind(t, lexer, TOKEN_KW_LOCAL)
	assertNextTokenKind(t, lexer, TOKEN_KW_NIL)
	assertNextTokenKind(t, lexer, TOKEN_OP_NOT)
	assertNextTokenKind(t, lexer, TOKEN_OP_OR)
	assertNextTokenKind(t, lexer, TOKEN_KW_REPEAT)
	assertNextTokenKind(t, lexer, TOKEN_KW_RETURN)
	assertNextTokenKind(t, lexer, TOKEN_KW_THEN)
	assertNextTokenKind(t, lexer, TOKEN_KW_TRUE)
	assertNextTokenKind(t, lexer, TOKEN_KW_UNTIL)
	assertNextTokenKind(t, lexer, TOKEN_KW_WHILE)
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken_identifiers(t *testing.T) {
	identifiers := `_ __ ___ a _HW_ hello_world HelloWorld HELLO_WORLD`
	lexer := NewLexer("str", identifiers)
	assertNextIdentifier(t, lexer, "_")
	assertNextIdentifier(t, lexer, "__")
	assertNextIdentifier(t, lexer, "___")
	assertNextIdentifier(t, lexer, "a")
	assertNextIdentifier(t, lexer, "_HW_")
	assertNextIdentifier(t, lexer, "hello_world")
	assertNextIdentifier(t, lexer, "HelloWorld")
	assertNextIdentifier(t, lexer, "HELLO_WORLD")
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken_numbers(t *testing.T) {
	numbers := `
	3   345   0xff   0xBEBADA
	
	3.0     3.1416     314.16e-2     0.31416E1     34e1
	0x0.1E  0xA23p-4   0X1.921FB54442D18P+1
	`
	lexer := NewLexer("str", numbers)
	assertNextNumber(t, lexer, "3")
	assertNextNumber(t, lexer, "345")
	assertNextNumber(t, lexer, "0xff")
	assertNextNumber(t, lexer, "0xBEBADA")
	assertNextNumber(t, lexer, "3.0")
	assertNextNumber(t, lexer, "3.1416")
	assertNextNumber(t, lexer, "314.16e-2")
	assertNextNumber(t, lexer, "0.31416E1")
	assertNextNumber(t, lexer, "34e1")
	assertNextNumber(t, lexer, "0x0.1E")
	assertNextNumber(t, lexer, "0xA23p-4")
	assertNextNumber(t, lexer, "0X1.921FB54442D18P+1")
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken_comments(t *testing.T) {
	lexer := NewLexer("str", `
	--
	--[[]]
	a -- short comment
	+ --[[ long comment ]] b --[===[ long
	comment
	]===] - c
	--`)
	assertNextIdentifier(t, lexer, "a")
	assertNextTokenKind(t, lexer, TOKEN_OP_ADD)
	assertNextIdentifier(t, lexer, "b")
	assertNextTokenKind(t, lexer, TOKEN_MINUS)
	assertNextIdentifier(t, lexer, "c")
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken_strings(t *testing.T) {
	strs := `
	[[]] [[ long string ]]
	[===[long
	string]===]
	'' '"' 'short string'
	"" "'" "short string"
	'\a\b\f\n\r\t\v\\\"\''
	'\8 \64 \122 \x08 \x7a \x7A \u{6211} zzz'
	'foo \z  
	
	bar'
	`
	lexer := NewLexer("str", strs)
	assertNextString(t, lexer, "")
	assertNextString(t, lexer, " long string ")
	assertNextString(t, lexer, "long\n\tstring")
	assertNextString(t, lexer, "")
	assertNextString(t, lexer, "\"")
	assertNextString(t, lexer, "short string")
	assertNextString(t, lexer, "")
	assertNextString(t, lexer, "'")
	assertNextString(t, lexer, "short string")
	assertNextString(t, lexer, "\a\b\f\n\r\t\v\\\"'")
	assertNextString(t, lexer, "\b @ z \b z z 我 zzz")
	assertNextString(t, lexer, "foo bar")
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken_hw(t *testing.T) {
	src := `print("Hello, World!")`
	lexer := NewLexer("str", src)

	assertNextIdentifier(t, lexer, "print")
	assertNextTokenKind(t, lexer, TOKEN_SEP_LPAREN)
	assertNextString(t, lexer, "Hello, World!")
	assertNextTokenKind(t, lexer, TOKEN_SEP_RPAREN)
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestLookAhead(t *testing.T) {
	src := `print("Hello, World!")`
	lexer := NewLexer("str", src)

	assert.IntEqual(t, lexer.LookAhead(1), TOKEN_IDENTIFIER)
	lexer.NextToken()
	assert.IntEqual(t, lexer.LookAhead(1), TOKEN_SEP_LPAREN)
	lexer.NextToken()
	assert.IntEqual(t, lexer.LookAhead(1), TOKEN_STRING)
	lexer.NextToken()
	assert.IntEqual(t, lexer.LookAhead(1), TOKEN_SEP_RPAREN)
	lexer.NextToken()
	assert.IntEqual(t, lexer.LookAhead(1), TOKEN_EOF)
}

func assertNextTokenKind(t *testing.T, lexer *Lexer, expectedKind int) {
	_, kind, _ := lexer.NextToken()
	assert.IntEqual(t, kind, expectedKind)
}

func assertNextIdentifier(t *testing.T, lexer *Lexer, expectedToken string) {
	_, kind, token := lexer.NextToken()
	assert.IntEqual(t, kind, TOKEN_IDENTIFIER)
	assert.StringEqual(t, token, expectedToken)
}

func assertNextNumber(t *testing.T, lexer *Lexer, expectedToken string) {
	_, kind, token := lexer.NextToken()
	assert.IntEqual(t, kind, TOKEN_NUMBER)
	assert.StringEqual(t, token, expectedToken)
}

func assertNextString(t *testing.T, lexer *Lexer, expectedToken string) {
	_, kind, token := lexer.NextToken()
	assert.IntEqual(t, kind, TOKEN_STRING)
	assert.StringEqual(t, token, expectedToken)
}
