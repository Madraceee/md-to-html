package main

type scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *scanner {
	return &scanner{
		source:  source,
		tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanTokens()
	}

	s.tokens = append(s.tokens, NewToken(EOF, "", s.line))
	return s.tokens
}

func (s *scanner) scanTokens() {
	c := s.advance()
	switch c {
	case ' ':
		s.addToken(SPACE)
	case '\n':
		s.addToken(NEWLINE)
		s.line++
	case '\t':
		s.addToken(TAB)
	case '>':
		s.addToken(GREATER)
	case '[':
		s.addToken(LEFT_BRACKET)
	case ']':
		s.addToken(RIGHT_BRACKET)
	case '(':
		s.addToken(LEFT_PARAN)
	case ')':
		s.addToken(RIGHT_PARAN)
	case '!':
		s.addToken(EXCLAMATION)
	case '\\':
		s.addToken(FORWARD_SLASH)
	case '=':
		s.addToken(EQUAL)
	// Multi Character
	case '`':
		s.code()
	case '-':
		{
			if s.match('-') && s.match('-') {
				s.addToken(TRIPLE_DASH)
			} else {
				s.addToken(DASH)
			}
		}
	case '_':
		{
			if s.match('_') {
				if s.match('_') {
					s.addToken(TRIPLE_UNDERSCORE)
				} else {
					s.addToken(DOUBLE_UNDERSCORE)
				}
			} else {
				s.addToken(UNDERSCORE)
			}
		}
	case '*':
		{
			if s.match('*') {
				if s.match('*') {
					s.addToken(TRIPLE_STAR)
				} else {
					s.addToken(DOUBLE_STAR)
				}
			} else {
				s.addToken(STAR)
			}
		}
	case '+':
		s.addToken(PLUS)
	case '#':
		s.heading()
	default:
		if isDigit(c) {
			for !s.isAtEnd() && isDigit(s.peek()) {
				s.advance()
			}

			// Number followed by period and space
			// Ex
			// 1. Hello
			if s.peek() == '.' && s.peekNext() == ' ' {
				s.advance()
				s.tokens = append(s.tokens, NewToken(LIST_NUMBER, s.source[s.start:s.current-2], s.line))
			} else {
				s.content()
			}
		} else {
			s.content()
		}
	}
}

// isAtEnd Check whether the end of the file has been reached.
func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// peek Get current rune
func (s *scanner) peek() rune {
	return rune(s.source[s.current])
}

// peekNext Get next rune if present
func (s *scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\n'
	}
	return rune(s.source[s.current+1])
}

func (s *scanner) match(expected rune) bool {
	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current += 1
	return true
}

// advance Returns the current character and moves to next character
func (s *scanner) advance() rune {
	defer func() { s.current += 1 }()
	return rune(s.source[s.current])
}

// addToken Adds the token to the list given a tokentype
func (s *scanner) addToken(tokenType TokenType) {
	s.tokens = append(s.tokens, NewToken(tokenType, "", s.line))
}

// heading # to ##### are considered heading else they are considered as normal string. They are only headings only if they are present at the start of the line. This will be checked by the parser
func (s *scanner) heading() {
	for !s.isAtEnd() && s.peek() == '#' {
		s.advance()
	}

	if s.match(' ') {
		length := s.current - s.start - 1
		switch length {
		case 1:
			s.addToken(H1)
		case 2:
			s.addToken(H2)
		case 3:
			s.addToken(H3)
		case 4:
			s.addToken(H4)
		case 5:
			s.addToken(H5)
		case 6:
			s.addToken(H6)
		default:
			s.content()
		}
	} else {
		s.content()
		return
	}
}

// content Contents of the markdown.
func (s *scanner) content() {
	for {
		if !s.isAtEnd() && isContentChar(s.peek()) {
			s.advance()
		} else if s.peek() == 92 && !isContentChar(s.peekNext()) { // To escapre characters. Ex "\\ \* \[ \]"
			s.advance()
			s.advance()
		} else {
			break
		}
	}

	s.tokens = append(s.tokens, NewToken(CONTENT, s.source[s.start:s.current], s.line))
}

// code Store code in a single token
func (s *scanner) code() {
	noOfBackticks := 1
	line := s.line
	for s.peek() == '`' {
		noOfBackticks++
		if noOfBackticks > 3 {
			break
		}
		s.advance()
	}

	for !s.isAtEnd() {
		if s.peek() == '`' {
			count := 0
			for !s.isAtEnd() && s.peek() == '`' {
				count++
				s.advance()
			}

			if count == noOfBackticks {
				s.tokens = append(s.tokens, NewToken(CODE, s.source[s.start:s.current], line))
				break
			}
		}

		if s.peek() == '\n' {
			s.line++
			if noOfBackticks != 3 {
				s.tokens = append(s.tokens, NewToken(CODE, s.source[s.start:s.current], line))
				break
			}
		}
		s.advance()
	}
}

func isContentChar(c rune) bool {
	// (,),*, [,\,]
	if (c >= 40 && c <= 42) || (c >= 91 && c <= 93) || (c >= 95 && c <= 96) || c == '\n' {
		return false
	}
	return true
}

func isSpecialCharacter(c rune) bool {
	return c == '*' || c == '_'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}
