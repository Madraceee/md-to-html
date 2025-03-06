package main

import "fmt"

type AstPrinter struct {
}

func (a *AstPrinter) VisitStringPara(s *String) (string, error) {
	return s.Content.Lexeme, nil
}

func (a *AstPrinter) VisitBoldPara(b *Bold) (string, error) {
	s := ""
	for _, c := range b.Content {
		content, _ := c.Visit(a)
		s += content
	}

	return fmt.Sprintf("<Bold %s>", s), nil
}

func (a *AstPrinter) VisitItalicsPara(i *Italics) (string, error) {
	s := ""
	for _, c := range i.Content {
		content, _ := c.Visit(a)
		s += content
	}

	return fmt.Sprintf("<Italics %s>", s), nil
}

func (a *AstPrinter) VisitWhitespacePara(w *Whitespace) (string, error) {
	return fmt.Sprintf("<Whitespace %s>", getTokenTypeString(w.Whitespace.TokenType)), nil
}

func (a *AstPrinter) VisitHeadingChunk(h *Heading) (string, error) {
	s := ""
	for _, c := range h.Content {
		content, _ := c.Visit(a)
		s += content
	}

	return fmt.Sprintf("<Header %s %s>", getTokenTypeString(h.Header.TokenType), s), nil
}

func (a *AstPrinter) VisitParagraphChunk(p *Paragraph) (string, error) {
	s := ""
	for _, c := range p.Content {
		content, _ := c.Visit(a)
		s += content

	}

	return fmt.Sprintf("<Para %s>", s), nil
}

func (a *AstPrinter) VisitLineChunk(*Line) (string, error) {
	return fmt.Sprintf("\n<Line >\n"), nil
}

func (a *AstPrinter) VisitLineBreakChunk(*LineBreak) (string, error) {
	return fmt.Sprintf("\n<Line Break>\n"), nil
}

func (a *AstPrinter) VisitCodeChunk(c *Code) (string, error) {
	return fmt.Sprintf("<Code\n %s\n>", c.Code.Lexeme), nil
}

func (a *AstPrinter) VisitListChunk(l *List) (string, error) {
	s := ""
	for i := range l.Content {

		content := ""
		for _, p := range l.Content[i] {
			c, _ := p.Visit(a)
			content += c
		}
		s += fmt.Sprintf("<List %s Level->%d Content->%s>", getTokenTypeString(l.ListType[i].TokenType), l.Level[i], content)
		s += "\n"
	}
	return s, nil
}

func (a *AstPrinter) VisitHTMLLinkPara(h *HTMLLink) (string, error) {
	title := ""
	for _, t := range h.Title {
		s, _ := t.Visit(a)
		title += s
	}

	link := ""
	for _, l := range h.Link {
		s, _ := l.Visit(a)
		link += s
	}

	return fmt.Sprintf("<HTMLLink Title->%s Link->%s>", title, link), nil
}
