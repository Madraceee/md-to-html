package main

import (
	"slices"
)

type Parser struct {
	Tokens  []Token
	Current int
}

const STACK_MAX = 64

func NewParser(tokens []Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
	}
}

func (p *Parser) Parse() {
	// chunks = make([]chunks,0)
	chunks := make([]Para, 0)
	for !p.isAtEnd() {
		chunk := p.paragraph(false)
		// chunk , err := p.chunk();
		// chunks = append(chunks, chunk)
		chunks = append(chunks, chunk...)
	}

	astPrinter := AstPrinter{}

	for _, c := range chunks {
		c.Visit(&astPrinter)
	}
}

func (p *Parser) heading() {
	// For heading
	if p.match(H1, H2, H3, H4, H5, H6) {
		// Heading followed by WhiteSpace and string
	} else {
		// Paragraph
		// p.paragraph()
	}
}

func (p *Parser) bold() Para {
	DPrintf("Inside Bold\n")
	p.match(DOUBLE_UNDERSCORE, DOUBLE_STAR)
	content := p.paragraph(true)

	DPrintf("After getting Bold content %v\n", content)
	p.match(DOUBLE_UNDERSCORE, DOUBLE_STAR)

	return NewBold(content)
}

func (p *Parser) italics() Para {
	DPrintf("Inside Italics\n")
	p.match(UNDERSCORE, STAR)
	content := p.paragraph(true)

	DPrintf("After getting Italics content %v\n", content)
	p.match(UNDERSCORE, STAR)

	return NewItalics(content)
}

func (p *Parser) paragraph(isInsideFormat bool) []Para {
	DPrintf("Inside Para\n")
	content := make([]Para, 0)
	for {
		if p.match(CONTENT, SPACE, TAB) {
			switch p.previos().TokenType {
			case CONTENT:
				content = append(content, NewString(p.previos()))
			case SPACE, TAB:
				content = append(content, NewWhitespace(p.previos()))
			}
		}

		if p.isFormat() && !isInsideFormat {
			content = append(content, p.format())
		} else if isInsideFormat {
			break
		}

		if p.match(NEWLINE) {
			break
		}
	}

	return content
}

func (p *Parser) format() Para {
	if p.match(STAR, UNDERSCORE) {
		return p.italics()
	} else if p.match(DOUBLE_STAR, DOUBLE_UNDERSCORE) {
		return p.bold()
	}

	DPrintf("BOLD AND ITALICS\n")
	p.match(TRIPLE_STAR, TRIPLE_UNDERSCORE)
	content := p.paragraph(true)
	p.match(TRIPLE_STAR, TRIPLE_UNDERSCORE)
	return NewBold([]Para{NewItalics(content)})
}

func (p *Parser) isFormat() bool {
	return p.peek().TokenType == STAR || p.peek().TokenType == DOUBLE_STAR || p.peek().TokenType == UNDERSCORE || p.peek().TokenType == DOUBLE_UNDERSCORE || p.peek().TokenType == TRIPLE_STAR || p.peek().TokenType == TRIPLE_UNDERSCORE
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
