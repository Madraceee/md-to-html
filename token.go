package main

import (
	"strconv"
)

type TokenType int

const (
	// Single Character Tokens
	SPACE TokenType = iota //
	NEWLINE
	TAB           // \t
	UNDERSCORE    // _
	GREATER       // Blockquote
	LEFT_BRACKET  // [
	RIGHT_BRACKET // ]
	LEFT_PARAN    // (
	RIGHT_PARAN   // )
	EXCLAMATION   // !
	FORWARD_SLASH // \
	EQUAL         // =
	CODE          // `

	// List
	LIST_NUMBER // Ordered List - Number followed by a .
	DASH        // Unordered List
	STAR        // Unordered List. Also used for Bold, Italics
	PLUS        // Unordered List

	// One or more characters
	H1                // #
	H2                // ##
	H3                // ###
	H4                // ####
	H5                // #####
	H6                // ######
	DOUBLE_STAR       // **
	DOUBLE_UNDERSCORE // __
	TRIPLE_STAR       // ***
	TRIPLE_UNDERSCORE // ___
	TRIPLE_DASH       // ---

	CONTENT // Contents

	EOF //
)

type Token struct {
	TokenType TokenType
	Lexeme    string
	Line      int
}

func NewToken(tokenType TokenType, lexeme string, line int) Token {
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Line:      line,
	}
}

func getTokenTypeString(tokenType TokenType) string {
	switch tokenType {
	case SPACE:
		return "SPACE"
	case NEWLINE:
		return "NEWLINE"
	case TAB:
		return "TAB"
	case UNDERSCORE:
		return "UNDERSCORE"
	case GREATER:
		return "GREATER"
	case LEFT_BRACKET:
		return "LEFT_BRACKET"
	case RIGHT_BRACKET:
		return "RIGHT_BRACKET"
	case LEFT_PARAN:
		return "LEFT_PARAN"
	case RIGHT_PARAN:
		return "RIGHT_PARAN"
	case EXCLAMATION:
		return "EXCLAMATION"
	case FORWARD_SLASH:
		return "FORWARD_SLASH"
	case EQUAL:
		return "EQUAL"
	case CODE:
		return "CODE"
	case LIST_NUMBER:
		return "LIST_NUMBER"
	case DASH:
		return "DASH"
	case STAR:
		return "STAR"
	case PLUS:
		return "PLUS"
	case H1:
		return "H1"
	case H2:
		return "H2"
	case H3:
		return "H3"
	case H4:
		return "H4"
	case H5:
		return "H5"
	case H6:
		return "H6"
	case DOUBLE_STAR:
		return "DOUBLE_STAR"
	case DOUBLE_UNDERSCORE:
		return "DOUBLE_UNDERSCORE"
	case TRIPLE_STAR:
		return "TRIPLE_STAR"
	case TRIPLE_UNDERSCORE:
		return "TRIPLE_UNDERSCORE"
	case TRIPLE_DASH:
		return "TRIPLE_DASH"
	case CONTENT:
		return "CONTENT"
	case EOF:
		return "EOF"
	}

	return ""
}

func GetTokenString(token *Token) string {
	return getTokenTypeString(token.TokenType) + "-" + token.Lexeme + " Line:" + strconv.Itoa(token.Line)
}
