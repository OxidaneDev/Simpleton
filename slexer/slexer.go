/*
	SLEXER.go
	SIMPLETON LANG
	v1.0
*/

package slexer

/*
	IMPORTS
*/

import (
	"strconv"
	"strings"
)

/*
	CONSTANTS
*/

const DIGITS = "0123456789"

/*
	ERRORS
*/

type Err struct {
	error_name string
	details    string
}

// Error implements error.
func (e *Err) Error() string {
	panic("unimplemented")
}

func NewErrz(error_name string, details string) *Err {
	return &Err{
		error_name: error_name,
		details:    details,
	}
}

func ErrAsString(setSelf Err) string {
	var result = "" + setSelf.error_name + ": " + setSelf.details + ""
	return result
}

/*
	ERROR TYPES
*/

func illegal_char_error(details string) *Err {
	return NewErrz("UNSUPPORTED CHARACTER", details)
}

/*
	TOKENS
*/

const TT_INT = "INT"
const TT_FLOAT = "FLOAT"
const TT_PLUS = "PLUS"
const TT_MINUS = "MINUS"
const TT_MUL = "MULTIPLY"
const TT_DIV = "DIVIDE"
const TT_LPAREN = "LPAREN"
const TT_RPAREN = "RPAREN"

type __Token__ struct {
	_type_  string
	_value_ any
}

func initToken(type_ string, value any) *__Token__ {
	return &__Token__{
		_type_:  type_,
		_value_: value,
	}
}

func stringifyToken(setToken __Token__) string {
	if setToken._value_ != nil {
		value, ok := setToken._value_.(int)
		if ok {
			return "" + setToken._type_ + ":" + strconv.Itoa(value)
		}
	}

	return "" + setToken._type_ + ""
}

/*
	MAIN LEXER
*/

type __lexer__ struct {
	text         string
	pos          int
	current_char string
}

func StartLexer(text string) *__lexer__ {
	return &__lexer__{
		text:         text,
		pos:          -1,
		current_char: "",
	}
}

func LexerAdvance(lexData __lexer__) {
	lexData.pos += 1
	if lexData.pos < len(lexData.text) {
		lexData.current_char = string(lexData.text[lexData.pos])
	} else {
		lexData.current_char = ""
	}
}

func LexerMakeNumber(lexData __lexer__) (__Token__, error) {
	var num_str string = ""
	var dot_count uint8 = 0

	for lexData.current_char != "" && (strings.Contains(DIGITS, lexData.current_char) || lexData.current_char == ".") {
		if lexData.current_char == "." {
			if dot_count >= 1 {
				break
			}
			dot_count += 1
			num_str += "."
		} else {
			num_str += lexData.current_char
		}
		LexerAdvance(lexData)
	}

	if dot_count == 0 {
		num, errz := strconv.Atoi(num_str)
		if errz != nil {
			panic("ERROR WHILE LEXING NUMBER")
		}
		return *initToken(TT_INT, num), nil
	} else {
		num, errz := strconv.ParseFloat(num_str, 64)
		if errz != nil {
			panic("ERROR WHILE LEXING FLOAT")
		}
		return *initToken(TT_FLOAT, num), nil
	}
}

func LexerTOKENS(lexData __lexer__) ([]any, error) {
	var retVar []any

	for lexData.current_char != "" {
		if lexData.current_char == "\t" {
			LexerAdvance(lexData)
		} else if strings.Contains(DIGITS, lexData.current_char) {
			token, err := LexerMakeNumber(lexData)

			if err != nil {
				panic("SOMETHING WRONG OCCURED WHILE PARSING NUMERICAL VALUES")
			}

			retVar = append(retVar, token)
		} else if lexData.current_char == "+" {
			retVar = append(retVar, initToken(TT_PLUS, nil))
			LexerAdvance(lexData)
		} else if lexData.current_char == "-" {
			retVar = append(retVar, initToken(TT_MINUS, nil))
			LexerAdvance(lexData)
		} else if lexData.current_char == "*" {
			retVar = append(retVar, initToken(TT_MUL, nil))
			LexerAdvance(lexData)
		} else if lexData.current_char == "/" {
			retVar = append(retVar, initToken(TT_DIV, nil))
			LexerAdvance(lexData)
		} else if lexData.current_char == "(" {
			retVar = append(retVar, initToken(TT_LPAREN, nil))
			LexerAdvance(lexData)
		} else if lexData.current_char == ")" {
			retVar = append(retVar, initToken(TT_RPAREN, nil))
			LexerAdvance(lexData)
		} else {
			var __chartemp__ = lexData.current_char
			LexerAdvance(lexData)
			return nil, illegal_char_error("'" + __chartemp__ + "'")
		}
	}

	return retVar, nil
}

func run(text string) ([]any, error) {
	Lexer := StartLexer(text)
	var Tokens, Errz = LexerTOKENS(*Lexer)

	var TokenStrings []any

	for _, token := range Tokens {
		tokenTyped, ok := token.(*__Token__)
		if !ok {
			panic("invalid token type")
		}
		TokenStrings = append(TokenStrings, stringifyToken(*tokenTyped))
	}

	return TokenStrings, Errz
}
