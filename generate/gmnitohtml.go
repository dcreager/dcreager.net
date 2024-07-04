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

package generate

import (
	"fmt"
	"html"
	"io"
	"net/url"
	"path"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
)

type HTMLWriter struct {
	domain      string
	out         io.Writer
	Title       string
	Path        string
	haveTitle   bool
	started     bool
	blockquote  bool
	pre         bool
	list        bool
	blanks      int
	sections    [3]bool
	anchors     map[string]struct{}
	useAbsLinks bool

	Published string
	Updated   string
}

func (h *HTMLWriter) isRoot() bool {
	return path.Base(h.Path) == "index.gmi"
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

func (h *HTMLWriter) generateAnchor(text string) string {
	var anchor string
	words := anchorRegexp.FindAllString(text, 4)
	if len(words) == 0 {
		anchor = "heading"
	} else {
		anchor = strings.Join(words, "-")
	}

	if h.anchors == nil {
		h.anchors = map[string]struct{}{}
	}
	if _, ok := h.anchors[anchor]; !ok {
		h.anchors[anchor] = struct{}{}
		return anchor
	}

	suffix := 1
	for {
		suffixedAnchor := fmt.Sprintf("%s-%d", anchor, suffix)
		if _, ok := h.anchors[suffixedAnchor]; !ok {
			h.anchors[suffixedAnchor] = struct{}{}
			return suffixedAnchor
		}
	}
}

func (h *HTMLWriter) spacingClass(classes ...string) string {
	if h.blanks == 0 {
		classes = append(classes, "nogap")
	} else if h.blanks > 1 {
		classes = append(classes, "doublegap")
	}
	if len(classes) == 0 {
		return ""
	}
	var css strings.Builder
	css.WriteString(" class=\"")
	for i, class := range classes {
		if i > 0 {
			css.WriteString(" ")
		}
		css.WriteString(class)
	}
	css.WriteString("\"")
	return css.String()
}

func gemlink(isRoot bool, href string) (*url.URL, string) {
	parsed, err := url.Parse(href)
	if err != nil {
		panic(err)
	}
	if strings.HasSuffix(href, ".gmi") {
		if !parsed.IsAbs() {
			href = translateGmiPath(parsed.String())
		}
	}
	if !isRoot && !parsed.IsAbs() && !path.IsAbs(href) {
		href = "../" + href
	}
	return parsed, href
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
	var blockquote bool
	switch line := line.(type) {
	case gemini.LineLink:
		parsed, href := gemlink(h.isRoot(), line.URL)
		href = html.EscapeString(href)
		name := html.EscapeString(line.Name)
		if name == "" {
			name = line.URL
		}

		linkIsRemote := parsed.IsAbs() || parsed.Host != ""
		if !linkIsRemote && h.useAbsLinks {
			href = fmt.Sprintf("https://%s%s", h.domain, href)
		}

		if isImage(href) {
			fmt.Fprintf(h.out, "<p%s><img src=\"%s\" alt=\"%s\"></p>\n", h.spacingClass(), href, name)
		} else {
			linkClass := "local"
			if linkIsRemote {
				linkClass = "remote"
			}
			fmt.Fprintf(h.out, "<p%s><a class=%s href=\"%s\">%s</a></p>\n", h.spacingClass(), linkClass, href, name)
		}
	case gemini.LinePreformattingToggle:
		h.pre = !h.pre
		if h.pre {
			alt := strings.TrimSpace(string(line))
			if alt != "" {
				alt = html.EscapeString(alt)
				langClass := "lang-" + alt
				fmt.Fprintf(h.out, "<pre%s aria-label='%s'>\n", h.spacingClass(langClass), alt)
			} else {
				fmt.Fprintf(h.out, "<pre%s>\n", h.spacingClass())
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
		heading := string(line)
		anchor := h.generateAnchor(heading)
		fmt.Fprintf(h.out, "<h1 id=%s><a class=header-link href=\"#%s\">¶</a><span>%s</span></h1>\n", anchor, anchor, html.EscapeString(heading))
	case gemini.LineHeading2:
		h.closeSection(3)
		h.closeSection(2)
		h.openSection(2)
		heading := string(line)
		anchor := h.generateAnchor(heading)
		fmt.Fprintf(h.out, "<h2 id=%s><a class=header-link href=\"#%s\">¶</a><span>%s</span></h2>\n", anchor, anchor, html.EscapeString(heading))
	case gemini.LineHeading3:
		h.closeSection(3)
		h.openSection(3)
		heading := string(line)
		anchor := h.generateAnchor(heading)
		fmt.Fprintf(h.out, "<h3 id=%s><a class=header-link href=\"#%s\">¶</a><span>%s</span></h3>\n", anchor, anchor, html.EscapeString(heading))
	case gemini.LineListItem:
		fmt.Fprintf(h.out, "<li%s>%s</li>\n", h.spacingClass(), renderLine(string(line)))
	case gemini.LineQuote:
		blockquote = true
		if line == "" {
			blank = true
		} else if dateRegexp.MatchString(string(line)) {
			h.Published = string(line[0:10])
			if updated := updatedRegexp.FindStringSubmatch(string(line)); updated != nil {
				h.Updated = updated[1]
			}
			fmt.Fprintf(h.out, "<p%s>%s</p>\n", h.spacingClass("date"), html.EscapeString(string(line)))
		} else {
			class := h.spacingClass()
			if h.blockquote {
				class = " class=nogap"
			}
			fmt.Fprintf(h.out, "<blockquote%s>%s</blockquote>\n", class, html.EscapeString(string(line)))
		}
	case gemini.LineText:
		if line == "" {
			blank = true
		} else {
			fmt.Fprintf(h.out, "<p%s>%s</p>\n", h.spacingClass(), renderLine(string(line)))
		}
	}
	h.blockquote = blockquote
	if blank {
		h.blanks++
	} else {
		h.blanks = 0
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
	gemURL := strings.TrimSuffix(fmt.Sprintf("gemini://%s/%s", h.domain, h.Path), "index.gmi")
	fmt.Fprintf(h.out, "<p class=gemlink>This page is also available via <a href=\"https://gemini.circumlunar.space/\">Gemini</a> at <a rel=alternate type=text/gemini href=\"%s\">%s</a></p>", gemURL, gemURL)
}
