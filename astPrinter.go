package main

import "fmt"

type AstPrinter struct {
}

func (a *AstPrinter) VisitStringPara(s *String) (string, error) {
	fmt.Printf("%s", s.Content.Lexeme)
	return s.Content.Lexeme, nil
}

func (a *AstPrinter) VisitBoldPara(b *Bold) (string, error) {
	fmt.Print("BOLD\n")
	for _, c := range b.Content {
		c.Visit(a)
	}
	fmt.Print("BOLD OVER\n")

	return "", nil
}

func (a *AstPrinter) VisitItalicsPara(i *Italics) (string, error) {
	fmt.Print("ITALICS\n")
	for _, c := range i.Content {
		c.Visit(a)
	}
	fmt.Print("ITALICS OVER\n")

	return "", nil
}

func (a *AstPrinter) VisitWhitespacePara(w *Whitespace) (string, error) {
	fmt.Printf("Whitespace %s\n", getTokenTypeString(w.Whitespace.TokenType))
	return "", nil
}

func (a *AstPrinter) VisitHeadingChunk(h *Heading) (string, error) {
	fmt.Printf("Header %s - ", getTokenTypeString(h.Header.TokenType))
	for _, c := range h.Content {
		c.Visit(a)
	}

	fmt.Printf("Header Over\n")

	return "", nil
}

func (a *AstPrinter) VisitParagraphChunk(p *Paragraph) (string, error) {
	for _, c := range p.Content {
		c.Visit(a)
	}
	fmt.Println("")

	return "", nil
}

func (a *AstPrinter) VisitLineChunk(*Line) (string, error) {
	fmt.Println("\n---LINE---")
	return "", nil
}

func (a *AstPrinter) VisitLineBreakChunk(*LineBreak) (string, error) {
	fmt.Println("Line Break")
	return "", nil
}

func (a *AstPrinter) VisitCodeChunk(c *Code) (string, error) {
	fmt.Printf("Code\n%s\n\n", c.Code.Lexeme)
	return "", nil
}

func (a *AstPrinter) VisitListChunk(l *List) (string, error) {
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
	return "", nil
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

	fmt.Printf("\nLink-> Title:%s  Link:%s\n", title, link)

	return "", nil
}
