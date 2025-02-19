package main

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
