package main

import "fmt"

type AstPrinter struct {
}

func (a *AstPrinter) VisitStringPara(s *String) (interface{}, error) {
	fmt.Printf("%s", s.Content.Lexeme)
	return nil, nil
}

func (a *AstPrinter) VisitBoldPara(b *Bold) (interface{}, error) {
	fmt.Print("BOLD\n")
	for _, c := range b.Content {
		c.Visit(a)
	}
	fmt.Print("BOLD OVER\n")

	return nil, nil
}

func (a *AstPrinter) VisitItalicsPara(i *Italics) (interface{}, error) {
	fmt.Print("ITALICS\n")
	for _, c := range i.Content {
		c.Visit(a)
	}
	fmt.Print("ITALICS OVER\n")

	return nil, nil
}

func (a *AstPrinter) VisitWhitespacePara(w *Whitespace) (interface{}, error) {
	fmt.Printf("Whitespace %s\n", getTokenTypeString(w.Whitespace.TokenType))
	return nil, nil
}

func (a *AstPrinter) VisitHeadingChunk(h *Heading) (interface{}, error) {
	fmt.Printf("Header %s - ", getTokenTypeString(h.Header.TokenType))
	for _, c := range h.Content {
		c.Visit(a)
	}

	fmt.Printf("Header Over\n")

	return nil, nil
}

func (a *AstPrinter) VisitParagraphChunk(p *Paragraph) (interface{}, error) {
	for _, c := range p.Content {
		c.Visit(a)
	}
	fmt.Println("")

	return nil, nil
}

func (a *AstPrinter) VisitLineChunk(*Line) (interface{}, error) {
	fmt.Println("\n---LINE---")
	return nil, nil
}

func (a *AstPrinter) VisitLineBreakChunk(*LineBreak) (interface{}, error) {
	fmt.Println("Line Break")
	return nil, nil
}

func (a *AstPrinter) VisitCodeChunk(c *Code) (interface{}, error) {
	fmt.Printf("Code\n%s\n\n", c.Code.Lexeme)
	return nil, nil
}

func (a *AstPrinter) VisitListChunk(l *List) (interface{}, error) {
	fmt.Print("Starting list\n")
	for i := range l.Content {
		for range l.Level[i] {
			fmt.Printf("\t")
		}
		fmt.Printf("%s %s ", getTokenTypeString(l.ListType[i].TokenType), l.ListType[i].Lexeme)
		for _, p := range l.Content[i] {
			p.Visit(a)
		}
		fmt.Println("")
	}
	fmt.Print("End list\n")
	return nil, nil
}
