package main

type Para interface {
	Visit(VisitPara) (interface{}, error)
}

type VisitPara interface {
	VisitStringPara(*String) (interface{}, error)
	VisitBoldPara(*Bold) (interface{}, error)
	VisitItalicsPara(*Italics) (interface{}, error)
	VisitWhitespacePara(*Whitespace) (interface{}, error)
}

type String struct {
	Content Token
}

func NewString(content Token) Para {
	return &String{
		Content: content,
	}
}

func (expr *String) Visit(visitor VisitPara) (interface{}, error) {
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

func (expr *Bold) Visit(visitor VisitPara) (interface{}, error) {
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

func (expr *Italics) Visit(visitor VisitPara) (interface{}, error) {
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

func (expr *Whitespace) Visit(visitor VisitPara) (interface{}, error) {
	return visitor.VisitWhitespacePara(expr)
}
