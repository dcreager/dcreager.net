// This file was originally extracted from https://git.sr.ht/~adnano/gmnitohtml/
// under the following copyright:
//
// Copyright (c) 2021 Adnan Maolood
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// Later edits to the file are
//
// Copyright (c) 2023 Douglas Creager

package main

import (
	"fmt"
	"html"
	"io"
	"net/url"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
)

type HTMLWriter struct {
	isRoot    bool
	out       io.Writer
	Title     string
	haveTitle bool
	started   bool
	pre       bool
	list      bool
	br        bool
	sections  [3]bool
}

func (h *HTMLWriter) openSection(level int) {
	for i := 1; i <= level; i++ {
		if !h.sections[i-1] {
			fmt.Fprintf(h.out, "<section class=level%d>\n", i)
			h.sections[i-1] = true
		}
	}
}

func (h *HTMLWriter) closeSection(level int) {
	if h.sections[level-1] {
		fmt.Fprintf(h.out, "</section> <!-- level %d -->\n", level)
		h.sections[level-1] = false
	}
}

func (h *HTMLWriter) Handle(line gemini.Line) {
	if !h.started {
		fmt.Fprint(h.out, "<article class=gemtext>\n")
		h.started = true
	}
	if _, ok := line.(gemini.LineListItem); ok {
		if !h.list {
			h.list = true
			fmt.Fprint(h.out, "<ul>\n")
		}
	} else if h.list {
		h.list = false
		fmt.Fprint(h.out, "</ul>\n")
	}
	var blank bool
	switch line := line.(type) {
	case gemini.LineLink:
		href := line.URL
		parsed, err := url.Parse(href)
		if err != nil {
			panic(err)
		}
		if strings.HasSuffix(href, ".gmi") {
			if !parsed.IsAbs() {
				href = translateGmiPath(parsed.String())
			}
		}
		if !h.isRoot && !parsed.IsAbs() {
			href = "../" + href
		}
		href = html.EscapeString(href)
		name := html.EscapeString(line.Name)
		if name == "" {
			name = line.URL
		}
		fmt.Fprintf(h.out, "<p><a href='%s'>%s</a></p>\n", href, name)
	case gemini.LinePreformattingToggle:
		h.pre = !h.pre
		if h.pre {
			alt := strings.TrimSpace(string(line))
			if alt != "" {
				alt = html.EscapeString(alt)
				fmt.Fprintf(h.out, "<pre aria-label='%s'>\n", alt)
			} else {
				fmt.Fprint(h.out, "<pre>\n")
			}
		} else {
			fmt.Fprint(h.out, "</pre>\n")
		}
	case gemini.LinePreformattedText:
		fmt.Fprintf(h.out, "%s\n", html.EscapeString(string(line)))
	case gemini.LineHeading1:
		if !h.haveTitle {
			h.Title = string(line)
			h.haveTitle = true
		}
		h.closeSection(3)
		h.closeSection(2)
		h.closeSection(1)
		h.openSection(1)
		fmt.Fprintf(h.out, "<h1>%s</h1>\n", html.EscapeString(string(line)))
	case gemini.LineHeading2:
		h.closeSection(3)
		h.closeSection(2)
		h.openSection(2)
		fmt.Fprintf(h.out, "<h2>%s</h2>\n", html.EscapeString(string(line)))
	case gemini.LineHeading3:
		h.closeSection(3)
		h.openSection(3)
		fmt.Fprintf(h.out, "<h3>%s</h3>\n", html.EscapeString(string(line)))
	case gemini.LineListItem:
		fmt.Fprintf(h.out, "<li>%s</li>\n", html.EscapeString(string(line)))
	case gemini.LineQuote:
		fmt.Fprintf(h.out, "<blockquote>%s</blockquote>\n", html.EscapeString(string(line)))
	case gemini.LineText:
		if line == "" {
			blank = true
			if h.br {
				fmt.Fprint(h.out, "<br/>\n")
			} else {
				h.br = true
			}
		} else {
			fmt.Fprintf(h.out, "<p>%s</p>\n", html.EscapeString(string(line)))
		}
	}
	if h.br && !blank {
		h.br = false
	}
}

func (h *HTMLWriter) Finish() {
	if h.pre {
		fmt.Fprint(h.out, "</pre>\n")
	}
	if h.list {
		fmt.Fprint(h.out, "</ul>\n")
	}
	h.closeSection(3)
	h.closeSection(2)
	h.closeSection(1)
	fmt.Fprint(h.out, "</article> <!-- gemtext -->\n")
}
