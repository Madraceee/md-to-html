package main

import (
	"fmt"
	"slices"
)

type Parser struct {
	Tokens   []Token
	Current  int
	Stack    []*Token
	StackTop int
}

const STACK_MAX = 8

func NewParser(tokens []Token) *Parser {
	return &Parser{
		Tokens:   tokens,
		Current:  0,
		Stack:    make([]*Token, STACK_MAX),
		StackTop: 0,
	}
}

func (p *Parser) Parse() []Chunk {
	chunks := make([]Chunk, 0)
	for !p.isAtEnd() {
		chunk := p.chunk()
		chunks = append(chunks, chunk)
	}

	return chunks
}

func (p *Parser) chunk() Chunk {
	DPrintf("In Chunk\n")
	if p.match(NEWLINE) {
		return NewLineBreak()
	}
	if p.isHeading() {
		return p.heading()
	} else if p.peek().TokenType == TRIPLE_DASH {
		p.advance()
		if p.match(SPACE, TAB, NEWLINE) {
			for !p.isAtEnd() {
				if p.match(NEWLINE) {
					break
				}
				p.advance()
			}
		} else {
			p.retreat()
			return NewParagraph(p.paragraph())
		}
		return NewLine()
	} else if p.match(CODE) {
		// Check whether code is single line or multiline
		noOfBackticks := 0
		code := p.previos().Lexeme

		for _, c := range code {
			if c == '`' {
				noOfBackticks++
				continue
			}
			break
		}

		defer p.advance()
		return NewCode(p.previos(), noOfBackticks)
	} else if p.match(LIST_NUMBER, DASH, STAR, PLUS) {
		DPrintf("In List %s\n", getTokenTypeString(p.previos().TokenType))
		p.retreat()
		listType := make([]Token, 0)
		levels := make([]int, 0)
		contents := make([][]Para, 0)
		for {
			level := 0
			for p.match(TAB) {
				level++
			}
			if p.match(LIST_NUMBER, DASH, STAR, PLUS) {
				levels = append(levels, level)
				listType = append(listType, p.previos())
			} else {
				p.retreat()
				break
			}

			contents = append(contents, p.paragraph())
		}
		// Consume last line break '\n'
		p.match(NEWLINE)
		return NewList(contents, levels, listType)
	} else {
		return NewParagraph(p.paragraph())
	}
}

func (p *Parser) heading() Chunk {
	DPrintf("Inside Heading\n")
	if p.match(H1, H2, H3, H4, H5, H6) {
		header := p.previos()
		// Get Space
		space := make([]Token, 0)
		for p.match(SPACE, TAB) {
			space = append(space, p.previos())
		}

		// Store Heading content
		content := p.paragraph()
		return NewHeading(header, content)
	} else {
		return NewParagraph(p.paragraph())
	}
}

func (p *Parser) bold() Para {
	DPrintf("Inside Bold\n")

	t := p.previos()
	contents := make([]Para, 0)
	for !p.match(DOUBLE_UNDERSCORE, DOUBLE_STAR) {
		content := p.paragraph()
		contents = append(contents, content...)

		if p.previos().TokenType == NEWLINE {
			fmt.Printf("Warning: Unclosed BOLD(%s) on line %d\n", getTokenTypeString(t.TokenType), t.Line)
			p.stackPop()
			break
		}

		if p.isAtEnd() {
			if p.stackTop().TokenType == DOUBLE_UNDERSCORE || p.stackTop().TokenType == DOUBLE_STAR {
				fmt.Printf("Warning: Unclosed BOLD(%s) on line %d\n", getTokenTypeString(t.TokenType), t.Line)
				p.stackPop()
			}
			break
		}
	}

	DPrintf("After getting Bold content %v\n", contents)

	return NewBold(contents)
}

func (p *Parser) italics() Para {
	DPrintf("Inside Italics\n")

	t := p.previos()
	contents := make([]Para, 0)
	for !p.match(UNDERSCORE, STAR) {
		content := p.paragraph()
		contents = append(contents, content...)

		if p.previos().TokenType == NEWLINE {
			fmt.Printf("Warning: Unclosed ITALICS(%s) on line %d\n", getTokenTypeString(t.TokenType), t.Line)
			p.stackPop()
			break
		}

		if p.isAtEnd() {
			if p.stackTop().TokenType == UNDERSCORE || p.stackTop().TokenType == STAR {
				fmt.Printf("Warning: Unclosed ITALICS(%s) on line %d\n", getTokenTypeString(t.TokenType), t.Line)
				p.stackPop()
			}
			break
		}
	}

	DPrintf("After getting Italics content %v\n", contents)

	return NewItalics(contents)
}

func (p *Parser) link() []Para {
	title := make([]Para, 0)
	// [ already consumed... consuming Title till ]
	for !p.isAtEnd() && !p.match(RIGHT_BRACKET) {
		title = append(title, NewString(p.peek()))
		p.advance()
	}

	// After ],  ( should be present for a link
	link := make([]Para, 0)
	if p.match(LEFT_PARAN) {
		for !p.isAtEnd() && !p.match(RIGHT_PARAN) {
			link = append(link, NewString(p.peek()))
			p.advance()
		}
		return []Para{NewHTMLLink(title, link)}
	}

	para := []Para{NewString(NewToken(CONTENT, "[", p.previos().Line))}
	para = append(para, title...)
	para = append(para, NewString(NewToken(CONTENT, "]", p.previos().Line)))
	return para

}

func (p *Parser) paragraph() []Para {
	DPrintf("Inside Para %s %s\n", getTokenTypeString(p.peek().TokenType), p.peek().Lexeme)
	content := make([]Para, 0)
	for {
		if p.match(CONTENT, SPACE, TAB, TRIPLE_DASH, LEFT_PARAN, RIGHT_PARAN, LEFT_BRACKET) {
			switch p.previos().TokenType {
			case CONTENT:
				content = append(content, NewString(p.previos()))
			case SPACE, TAB:
				content = append(content, NewWhitespace(p.previos()))
			case TRIPLE_DASH, LEFT_PARAN, RIGHT_PARAN:
				content = append(content, NewString(NewToken(CONTENT, p.previos().Lexeme, p.previos().Line)))
			case LEFT_BRACKET:
				content = append(content, p.link()...)
			}
		}

		if p.isFormat() {
			if p.stackTop() != nil && p.stackPop().TokenType == p.peek().TokenType {
				break
			}
			content = append(content, p.format())
		}

		if p.match(NEWLINE) || p.previos().TokenType == NEWLINE {
			break
		}

		if p.peek().TokenType == CODE {
			break
		}

		if p.isAtEnd() {
			break
		}
	}

	return content
}

func (p *Parser) format() Para {
	token := p.peek()
	p.stackPush(&token)
	defer func() {
		if p.stackTop() != nil && p.stackTop().TokenType == token.TokenType {
			p.stackPop()
		}
	}()

	if p.match(STAR, UNDERSCORE) {
		return p.italics()
	} else if p.match(DOUBLE_STAR, DOUBLE_UNDERSCORE) {
		return p.bold()
	}

	DPrintf("BOLD AND ITALICS\n")
	p.match(TRIPLE_STAR, TRIPLE_UNDERSCORE)
	content := p.paragraph()
	p.match(TRIPLE_STAR, TRIPLE_UNDERSCORE)
	return NewBold([]Para{NewItalics(content)})
}

func (p *Parser) isFormat() bool {
	return p.peek().TokenType == STAR || p.peek().TokenType == DOUBLE_STAR || p.peek().TokenType == UNDERSCORE || p.peek().TokenType == DOUBLE_UNDERSCORE || p.peek().TokenType == TRIPLE_STAR || p.peek().TokenType == TRIPLE_UNDERSCORE
}

func (p *Parser) isHeading() bool {
	return p.peek().TokenType == H1 || p.peek().TokenType == H2 || p.peek().TokenType == H3 || p.peek().TokenType == H4 || p.peek().TokenType == H5 || p.peek().TokenType == H6
}

// Parser helper functions
func (p *Parser) isAtEnd() bool {
	return p.Tokens[p.Current].TokenType == EOF
}

func (p *Parser) match(_type ...TokenType) bool {
	res := slices.ContainsFunc(_type, func(t TokenType) bool {
		return t == p.peek().TokenType
	})

	if res {
		p.advance()
	}
	return res
}

func (p *Parser) peek() Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previos() Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.Current++
	}

	return p.previos()
}

// TODO:
// Bad structuring...
// Refactor code so main function reads the first then branches out into child
// Refactor bold() and italics()
func (p *Parser) retreat() Token {
	if p.Current != 0 {
		p.Current--
	}

	return p.peek()
}

// Stack Methods
func (p *Parser) stackPush(t *Token) {
	// Should return false or error
	if p.StackTop == STACK_MAX {
		return
	}
	p.Stack[p.StackTop] = t
	p.StackTop++
}

func (p *Parser) stackTop() *Token {
	if p.StackTop > 0 {
		return p.Stack[p.StackTop-1]
	}

	return nil
}

func (p *Parser) stackPop() *Token {
	if p.StackTop > 0 {
		p.StackTop--
		return p.Stack[p.StackTop]
	}

	return nil
}
