package main

type Para interface {
	Visit(VisitPara) (string, error)
}

type VisitPara interface {
	VisitStringPara(*String) (string, error)
	VisitBoldPara(*Bold) (string, error)
	VisitItalicsPara(*Italics) (string, error)
	VisitWhitespacePara(*Whitespace) (string, error)
	VisitHTMLLinkPara(*HTMLLink) (string, error)
}

type String struct {
	Content Token
}

func NewString(content Token) Para {
	return &String{
		Content: content,
	}
}

func (expr *String) Visit(visitor VisitPara) (string, error) {
	return visitor.VisitStringPara(expr)
}

type Bold struct {
	Content []Para
}

func NewBold(content []Para) Para {
	return &Bold{
		Content: content,
	}
}

func (expr *Bold) Visit(visitor VisitPara) (string, error) {
	return visitor.VisitBoldPara(expr)
}

type Italics struct {
	Content []Para
}

func NewItalics(content []Para) Para {
	return &Italics{
		Content: content,
	}
}

func (expr *Italics) Visit(visitor VisitPara) (string, error) {
	return visitor.VisitItalicsPara(expr)
}

type Whitespace struct {
	Whitespace Token
}

func NewWhitespace(whitespace Token) Para {
	return &Whitespace{
		Whitespace: whitespace,
	}
}

func (expr *Whitespace) Visit(visitor VisitPara) (string, error) {
	return visitor.VisitWhitespacePara(expr)
}

type HTMLLink struct {
	Title []Para
	Link  []Para
}

func NewHTMLLink(Title []Para, Link []Para) Para {
	return &HTMLLink{
		Title: Title,
		Link:  Link,
	}
}

func (expr *HTMLLink) Visit(visitor VisitPara) (string, error) {
	return visitor.VisitHTMLLinkPara(expr)
}
