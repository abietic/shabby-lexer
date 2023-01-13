// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package token defines constants representing the lexical tokens of the Go
// programming language and basic operations on tokens (printing, predicates).
package token

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

// Token is the set of lexical tokens of the Go programming language.
type Token int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	literal_beg // 字面量开始
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT // main
	INT   // 12345
	FLOAT // 123.45
	// IMAG   // 123.45i
	// CHAR   // 'a' // 只用字符串类吧,单引双引都代表单引
	STRING      // "abc"
	literal_end // 字面量结束

	operator_beg // 操作符开始
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND // &
	OR  // |
	XOR // ^
	SHL // <<
	SHR // >>
	// AND_NOT // &^ // 不使用这个改用正常c的按位非
	NOT // ~ // 新加入的c中的按位非

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN // &=
	OR_ASSIGN  // |=
	XOR_ASSIGN // ^=
	SHL_ASSIGN // <<=
	SHR_ASSIGN // >>=
	// AND_NOT_ASSIGN // &^= // 不使用这个改用正常c的按位非

	LAND // &&
	LOR  // ||
	// ARROW // <- // 不用channel了所以不用
	INC // ++
	DEC // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	LNOT   // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN       // )
	RBRACK       // ]
	RBRACE       // }
	SEMICOLON    // ;
	COLON        // :
	operator_end // 操作符结束

	keyword_beg // 关键词开始
	// Keywords
	BREAK
	CASE
	// CHAN // 不用channel管道
	CONST
	CONTINUE

	DEFAULT
	DEFER // 对于没有exception的语言来说还是应该加上
	ELSE
	FALLTHROUGH
	FOR
	WHILE // 添加while

	FUNC
	// GO // 不使用go的协程了
	// GOTO // 不加goto
	IF
	IMPORT

	INTERFACE
	MAP
	PACKAGE
	RANGE
	RETURN

	// SELECT // 因为不用管道所以select也不用
	STRUCT
	SWITCH
	TYPE
	VAR
	keyword_end // 关键词结束

	additional_beg
	// additional tokens, handled in an ad-hoc manner
	// TILDE // 不用用于泛型的~符号把这个符号让给c的按位非
	additional_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT: "IDENT",
	INT:   "INT",
	FLOAT: "FLOAT",
	// IMAG:   "IMAG",
	// CHAR:   "CHAR",
	STRING: "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND: "&",
	OR:  "|",
	XOR: "^",
	SHL: "<<",
	SHR: ">>",
	// AND_NOT: "&^", //
	NOT: "~", // 使用与c类似的按位非操作

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN: "&=",
	OR_ASSIGN:  "|=",
	XOR_ASSIGN: "^=",
	SHL_ASSIGN: "<<=",
	SHR_ASSIGN: ">>=",
	// AND_NOT_ASSIGN: "&^=",

	LAND: "&&",
	LOR:  "||",
	// ARROW: "<-", // 不使用channel所以放弃左箭头表达式
	INC: "++",
	DEC: "--",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	LNOT:   "!",

	NEQ:      "!=",
	LEQ:      "<=",
	GEQ:      ">=",
	DEFINE:   ":=",
	ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	BREAK: "break",
	CASE:  "case",
	// CHAN:     "chan", // 放弃管道支持
	CONST:    "const",
	CONTINUE: "continue",

	DEFAULT:     "default",
	DEFER:       "defer",
	ELSE:        "else",
	FALLTHROUGH: "fallthrough",
	FOR:         "for",
	WHILE:       "while",

	FUNC: "func",
	// GO:     "go",
	// GOTO:   "goto",
	IF:     "if",
	IMPORT: "import",

	INTERFACE: "interface",
	MAP:       "map",
	PACKAGE:   "package",
	RANGE:     "range",
	RETURN:    "return",

	// SELECT: "select", // 放弃了channel所以也放弃select
	STRUCT: "struct",
	SWITCH: "switch",
	TYPE:   "type",
	VAR:    "var",

	// TILDE: "~",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// 表达式的优先级,没有操作符的表达式为最低优先级,其次是前有单目操作符的表达式
//
// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence, followed by operators
// starting with precedence 1 up to unary operators. The highest
// precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
func (op Token) Precedence() int {
	switch op {
	case LOR: // 逻辑或优先级最低,使与或表达式不用加括号
		return 1
	case LAND: // 逻辑与优先级高一点,保证表达式的逻辑运算已经完成
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ: // 所有的比较运算符,比逻辑运算符优先级高,比计算运算符低,保证比较值前完成数值运算
		return 3
	case ADD, SUB, OR, XOR: // 加减法与按位或和按位异或在同一优先级
		return 4
	case MUL, QUO, REM, SHL, SHR, AND /*, AND_NOT*/ : // 乘除法,求余和按位移动,按位与是双目运算符的最高优先级
		return 5
	}
	return LowestPrec
}

// 私有全局变量keyword,将字符串转化成相应关键字Token的映射
var keywords map[string]Token

// 初始化keyword这个私有全局变量的方法
func init() {
	keywords = make(map[string]Token, keyword_end-(keyword_beg+1))
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// 判断一个字符串是标识符还是关键字
// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

// 判断谓词
// Predicates

// 检查token是否是一个字面量
// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// 检查token是否是一个操作符
// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
func (tok Token) IsOperator() bool {
	return (operator_beg < tok && tok < operator_end) /*|| tok == TILDE*/
}

// 检查token是否是关键字
// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

// 检查名字是否是一个大写字母开头的名字
// IsExported reports whether name starts with an upper-case letter.
func IsExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

// 检查一个字符串是否是一个关键字
// IsKeyword reports whether name is a Go keyword, such as "func" or "return".
func IsKeyword(name string) bool {
	// TODO: opt: use a perfect hash function instead of a global map.
	_, ok := keywords[name]
	return ok
}

// 检查字符串是否是一个标识符
// 一个标识符是非空的有数字,字母和下划线组成的,首个字符不是数字的字符串
// 关键字不是标识符
// IsIdentifier reports whether name is a Go identifier, that is, a non-empty
// string made up of letters, digits, and underscores, where the first character
// is not a digit. Keywords are not identifiers.
func IsIdentifier(name string) bool {
	// 如果是空串或关键字那不是标识符
	if name == "" || IsKeyword(name) {
		return false
	}
	for i, c := range name {
		// 由字母,下划线,数字组成,首个字符不是数字
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return true
}
