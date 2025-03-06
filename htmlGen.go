package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

// templates
var (
	templates = map[string]string{
		"p":             "<p>{{.}}</p>",
		"bold":          "<strong>{{.}}</strong>",
		"italics":       "<em>{{.}}</em>",
		"h1":            "<h1>{{.}}</h1>",
		"h2":            "<h2>{{.}}</h2>",
		"h3":            "<h3>{{.}}</h3>",
		"h4":            "<h4>{{.}}</h4>",
		"h5":            "<h5>{{.}}</h5>",
		"h6":            "<h6>{{.}}</h6>",
		"line":          "\n<hr>\n",
		"break":         "\n\n",
		"orderedList":   "<ol>\n{{.}}</ol>",
		"unorderedList": "<ul>\n{{.}}</ul>",
		"listItem":      "<li>{{.}}</li>",
		"link":          "<a href={{.Link}} target='_blank'>{{.Title}}</a>",
		"code-multiline": `<pre><code>
{{.}}
</code></pre>`,
		"code-go": `<pre><code class="language-go">
{{.}}
</code></pre>`,
		"code-singleLine": `<code>{{.}}</code>`,
	}
)

type data struct {
	Title string
	Date  string
	Body  string
}

type HTMLGenerator struct {
	data     data
	template map[string]*template.Template
}

func NewHTMLGenerator() *HTMLGenerator {
	genTemplate := make(map[string]*template.Template)
	for k, v := range templates {
		t, err := template.New(k).Parse(v)
		if err != nil {
			log.Fatalln("Error while parsing template ", err)
		}

		genTemplate[k] = t
	}
	return &HTMLGenerator{
		data:     data{},
		template: genTemplate,
	}
}

//go:embed template/html.txt
var baseHTML string

func (hg *HTMLGenerator) GenerateHTML(filename string, chunks []Chunk) {
	t, err := template.New("webpage").Parse(baseHTML)
	if err != nil {
		log.Fatal(err)
	}

	hg.data.Title = getTitle(filename)
	outputFilename := strings.Split(filename, ".")[0] + ".html"
	output, err := os.Create("./" + outputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	items := make([]string, 0)
	for _, c := range chunks {
		s, err := c.Visit(hg)
		if err != nil {
			log.Println("Error while generating HTML ", err)
			continue
		}

		items = append(items, s)
	}

	body := ""
	for _, item := range items {
		body += item + "\n"
	}
	hg.data.Body = body

	year, month, date := time.Now().Date()
	dateString := fmt.Sprintf("%d-%s-%d", date, month.String(), year)
	hg.data.Date = dateString

	var buffer bytes.Buffer
	err = t.Execute(&buffer, hg.data)
	if err != nil {
		log.Fatal(err)
	}

	_, err = output.Write(buffer.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

// Functions to support AST operations
func (hg *HTMLGenerator) VisitStringPara(s *String) (string, error) {
	return s.Content.Lexeme, nil
}

func (hg *HTMLGenerator) VisitBoldPara(b *Bold) (string, error) {
	s := ""
	for _, c := range b.Content {
		output, err := c.Visit(hg)
		if err != nil {
			return "", err
		}
		s += output
	}

	var buff bytes.Buffer
	err := hg.template["bold"].Execute(&buff, s)
	return buff.String(), err
}

func (hg *HTMLGenerator) VisitItalicsPara(i *Italics) (string, error) {
	s := ""
	for _, c := range i.Content {
		output, err := c.Visit(hg)
		if err != nil {
			return "", err
		}
		s += output
	}

	var buff bytes.Buffer
	err := hg.template["italics"].Execute(&buff, s)
	return buff.String(), err
}

func (hg *HTMLGenerator) VisitWhitespacePara(w *Whitespace) (string, error) {
	if w.Whitespace.TokenType == SPACE {
		return " ", nil
	} else if w.Whitespace.TokenType == TAB {
		return "\t", nil
	}

	return "", nil
}

func (hg *HTMLGenerator) VisitHeadingChunk(h *Heading) (string, error) {
	s := ""
	for _, c := range h.Content {
		output, err := c.Visit(hg)
		if err != nil {
			return "", err
		}
		s += output
	}

	var err error
	var buff bytes.Buffer
	switch h.Header.TokenType {
	case H1:
		err = hg.template["h1"].Execute(&buff, s)
	case H2:
		err = hg.template["h2"].Execute(&buff, s)
	case H3:
		err = hg.template["h3"].Execute(&buff, s)
	case H4:
		err = hg.template["h4"].Execute(&buff, s)
	case H5:
		err = hg.template["h5"].Execute(&buff, s)
	case H6:
		err = hg.template["h6"].Execute(&buff, s)
	default:
		err = errors.New("Header not found")
	}

	return buff.String(), err
}

func (hg *HTMLGenerator) VisitParagraphChunk(p *Paragraph) (string, error) {
	s := ""
	for _, c := range p.Content {
		output, err := c.Visit(hg)
		if err != nil {
			return "", err
		}
		s += output
	}

	var buff bytes.Buffer
	err := hg.template["p"].Execute(&buff, s)
	return buff.String(), err
}

func (hg *HTMLGenerator) VisitLineChunk(*Line) (string, error) {
	return templates["line"], nil
}

func (hg *HTMLGenerator) VisitLineBreakChunk(*LineBreak) (string, error) {
	return templates["break"], nil
}

func (hg *HTMLGenerator) VisitCodeChunk(c *Code) (string, error) {
	code := c.Code.Lexeme

	key := "code-singleLine"
	if c.NoOfLines == 3 {
		key = "code-multiline"
		lines := strings.Split(code[3:len(code)-3], "\n")
		if len(lines) > 1 {
			language := strings.Trim(lines[0], " ")
			if language != "" {
				key = "code-" + language
			}
		}

		code = strings.Join(lines[1:], "\n")
	} else {
		code = code[c.NoOfLines : len(code)-c.NoOfLines]
	}

	code = strings.TrimRight(code, " \n\t")

	var buff bytes.Buffer
	if _, ok := hg.template[key]; !ok {
		return "", errors.New("Lanugage not present")
	}

	err := hg.template[key].Execute(&buff, code)
	return buff.String(), err
}

func (hg *HTMLGenerator) VisitListChunk(l *List) (string, error) {
	listItems := make([]string, 0)
	for i := range l.Content {
		s := ""
		for _, p := range l.Content[i] {
			output, err := p.Visit(hg)
			if err != nil {
				return "", err
			}
			out := strings.Trim(output, " ")
			s += out
		}

		var buff bytes.Buffer
		err := hg.template["listItem"].Execute(&buff, s)
		if err != nil {
			return "", nil
		}
		listItems = append(listItems, buff.String())
	}

	listMap := make(map[int]string)
	levelTypeMap := make(map[int]TokenType)

	level := l.Level[0]
	levelTypeMap[level] = l.ListType[0].TokenType

	for i := range l.Level {
		if l.Level[i] == level {
			listMap[l.Level[i]] += listItems[i] + "\n"
		}

		if l.Level[i] > level {
			levelTypeMap[l.Level[i]] = l.ListType[i].TokenType
			listMap[l.Level[i]] += listItems[i] + "\n"
			level = l.Level[i]
		}

		if l.Level[i] < level {
			for j := level; j > l.Level[i]; j-- {
				s := listMap[j]
				var err error
				var buff bytes.Buffer
				switch levelTypeMap[j] {
				case LIST_NUMBER:
					err = hg.template["orderedList"].Execute(&buff, s)
				case DASH, STAR, PLUS:
					err = hg.template["unorderedList"].Execute(&buff, s)
				}
				if err != nil {
					return "", err
				}

				listMap[j-1] += buff.String() + "\n"
				delete(listMap, j)
			}
			level = l.Level[i]
			listMap[level] += listItems[i] + "\n"
		}
	}

	for i := level; i >= 0; i-- {
		s := listMap[i]
		var err error
		var buff bytes.Buffer
		switch levelTypeMap[i] {
		case LIST_NUMBER:
			err = hg.template["orderedList"].Execute(&buff, s)
		case DASH, STAR, PLUS:
			err = hg.template["unorderedList"].Execute(&buff, s)
		}
		if err != nil {
			return "", err
		}

		if i == 0 {
			listMap[i] = buff.String() + "\n"
			break
		}

		listMap[i-1] += buff.String() + "\n"
	}

	return listMap[0], nil
}

func (hg *HTMLGenerator) VisitHTMLLinkPara(h *HTMLLink) (string, error) {
	title := ""
	for _, t := range h.Title {
		s, _ := t.Visit(hg)
		title += s
	}

	link := ""
	for _, l := range h.Link {
		s, _ := l.Visit(hg)
		link += s
	}

	var buff bytes.Buffer
	data := struct {
		Title string
		Link  string
	}{Title: title, Link: link}
	err := hg.template["link"].Execute(&buff, data)

	return buff.String(), err
}

// Helper functions
func getTitle(filename string) string {
	filenameWithoutExt := strings.Split(filename, ".")[0]
	wordArr := strings.FieldsFunc(filenameWithoutExt, func(r rune) bool {
		return r == ' ' || r == '-'
	})

	for i, word := range wordArr {
		wordArr[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
	}

	return strings.Join(wordArr, " ")
}
