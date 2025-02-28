package main

type Chunk interface {
	Visit(VisitChunk) (interface{}, error)
}

type VisitChunk interface {
	VisitHeadingChunk(*Heading) (interface{}, error)
	VisitParagraphChunk(*Paragraph) (interface{}, error)
	VisitLineChunk(*Line) (interface{}, error)
	VisitLineBreakChunk(*LineBreak) (interface{}, error)
}

type Heading struct {
	Header  Token
	Content []Para
}

func NewHeading(header Token, content []Para) Chunk {
	return &Heading{
		Header:  header,
		Content: content,
	}
}

func (expr *Heading) Visit(visitor VisitChunk) (interface{}, error) {
	return visitor.VisitHeadingChunk(expr)
}

type Paragraph struct {
	Content []Para
}

func NewParagraph(content []Para) Chunk {
	return &Paragraph{
		Content: content,
	}
}

func (expr *Paragraph) Visit(visitor VisitChunk) (interface{}, error) {
	return visitor.VisitParagraphChunk(expr)
}

type Line struct {
}

func NewLine() Chunk {
	return &Line{}
}

func (expr *Line) Visit(visitor VisitChunk) (interface{}, error) {
	return visitor.VisitLineChunk(expr)
}

type LineBreak struct {
}

func NewLineBreak() Chunk {
	return &LineBreak{}
}

func (expr *LineBreak) Visit(visitor VisitChunk) (interface{}, error) {
	return visitor.VisitLineBreakChunk(expr)
}
