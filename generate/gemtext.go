package generate

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
)

// Gemfile is an OutputFile that is generated from the contents of a gemtext
// file.
type Gemfile struct {
	sourcePath string
	outputPath string
	mode       fs.FileMode
	Source     gemini.Text
	content    string
	Title      string
	Published  string
	Updated    string
}

// GemfileVisitor lets you easily consume the contents of a gemtext source file.
type GemfileVisitor interface {
	Start() error
	Before(line gemini.Line) error
	After(line gemini.Line) error
	Finish() error

	Heading1(line gemini.LineHeading1) error
	Heading2(line gemini.LineHeading2) error
	Heading3(line gemini.LineHeading3) error
	Link(line gemini.LineLink) error
	ListItem(line gemini.LineListItem) error
	PreformattedText(line gemini.LinePreformattedText) error
	PreformattingToggle(line gemini.LinePreformattingToggle) error
	Quote(line gemini.LineQuote) error
	Text(line gemini.LineText) error
}

func urlForGempath(gempath string) string {
	if trimmed := strings.TrimSuffix(gempath, "index.gmi"); trimmed != gempath {
		return trimmed
	} else {
		return strings.TrimSuffix(gempath, ".gmi") + "/"
	}
}

func NewGemfile(site *Site, sourcePath string, sourceDir string, mode fs.FileMode) (*Gemfile, error) {
	outputPath := path.Join(urlForGempath(sourcePath), "index.html")

	in, err := os.Open(path.Join(sourceDir, sourcePath))
	if err != nil {
		return nil, err
	}
	defer in.Close()

	source, err := gemini.ParseText(in)
	if err != nil {
		return nil, err
	}

	f := &Gemfile{
		sourcePath: sourcePath,
		outputPath: outputPath,
		mode:       mode,
		Source:     source,
	}
	hw := GemtextHTML{Site: site, Gemfile: f}
	err = f.Visit(&hw)
	if err != nil {
		return nil, err
	}
	f.content = hw.content.String()

	return f, nil
}

func (f *Gemfile) SourceIsRootGemfile() bool {
	return path.Base(f.sourcePath) == "index.gmi"
}

func (f *Gemfile) Visit(visitor GemfileVisitor) error {
	err := visitor.Start()
	if err != nil {
		return err
	}

	for _, line := range f.Source {
		err := visitor.Before(line)
		if err != nil {
			return err
		}

		switch line := line.(type) {
		case gemini.LineHeading1:
			err := visitor.Heading1(line)
			if err != nil {
				return err
			}
		case gemini.LineHeading2:
			err := visitor.Heading2(line)
			if err != nil {
				return err
			}
		case gemini.LineHeading3:
			err := visitor.Heading3(line)
			if err != nil {
				return err
			}
		case gemini.LineLink:
			err := visitor.Link(line)
			if err != nil {
				return err
			}
		case gemini.LineListItem:
			err := visitor.ListItem(line)
			if err != nil {
				return err
			}
		case gemini.LinePreformattedText:
			err := visitor.PreformattedText(line)
			if err != nil {
				return err
			}
		case gemini.LinePreformattingToggle:
			err := visitor.PreformattingToggle(line)
			if err != nil {
				return err
			}
		case gemini.LineQuote:
			err := visitor.Quote(line)
			if err != nil {
				return err
			}
		case gemini.LineText:
			err := visitor.Text(line)
			if err != nil {
				return err
			}
		}

		err = visitor.After(line)
		if err != nil {
			return err
		}
	}

	return visitor.Finish()
}

func (f *Gemfile) OutputPath() string {
	return f.outputPath
}

func (f *Gemfile) Generate(site *Site) error {
	return site.OutputHTML(f.outputPath, f.mode, f.Title, f.content)
}

func ProcessSourceDir(domain, sourceDir, outputDir string) (int, error) {
	count := 0
	sourceLen := len(sourceDir) + 1
	err := filepath.WalkDir(sourceDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		count++
		base := p[sourceLen:]
		ext := path.Ext(base)
		output := path.Join(outputDir, base)

		if ext == ".gmi" {
			return translateGemtext(domain, base, p, output, d)
		}

		return copyFile(base, p, output, d)
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func copyFile(base, from, to string, d fs.DirEntry) error {
	err := os.MkdirAll(path.Dir(to), 0750)
	if err != nil {
		return err
	}

	in, err := os.Open(from)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := d.Info()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

var templates = map[string]*template.Template{}

func loadTemplate(base string) (*template.Template, error) {
	if t, ok := templates[base]; ok {
		return t, nil
	}

	templatePath := path.Join("templates", base)
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	templates[base] = t
	return t, nil
}

func findTemplate(base string) (*template.Template, error) {
	t, err := loadTemplate(base)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	t, err = loadTemplate("default.html")
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	return nil, fmt.Errorf("Cannot find suitable template for %s", base)
}

func loadMetadataFile(base string, meta string) (string, error) {
	for {
		metadataPath := path.Join("metadata", base, meta)
		metadata, err := os.ReadFile(metadataPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				next := path.Dir(base)
				if next == base {
					return "", nil
				}
				base = next
				continue
			}
			return "", err
		}
		return string(metadata), nil
	}
}

func translateGmiPath(p string) string {
	if path.Base(p) == "index.gmi" {
		return p[:len(p)-len("index.gmi")] + "/"
	} else {
		return p[:len(p)-len(".gmi")] + "/"
	}
}

type geminiTemplateData struct {
	Title      string
	Content    template.HTML
	CustomHead template.HTML
}

func translateGemtext(domain, base, from, to string, d fs.DirEntry) error {
	to = translateGmiPath(to) + "index.html"
	err := os.MkdirAll(path.Dir(to), 0750)
	if err != nil {
		return err
	}

	in, err := os.Open(from)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := d.Info()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	var content strings.Builder
	hw := HTMLWriter{
		domain: domain,
		Path:   base,
		out:    &content,
	}
	gemini.ParseLines(in, hw.Handle)
	hw.Finish()

	customHead, err := loadMetadataFile(base, "custom-head")
	if err != nil {
		return err
	}

	data := geminiTemplateData{
		Title:      hw.Title,
		Content:    template.HTML(content.String()),
		CustomHead: template.HTML(customHead),
	}
	t, err := findTemplate(base)
	if err != nil {
		return err
	}
	err = t.Execute(out, data)
	if err != nil {
		return err
	}

	return nil
}

var anchorRegexp = regexp.MustCompile("\\w+")
var dateRegexp = regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}")
var updatedRegexp = regexp.MustCompile("last updated ([0-9]{4}-[0-9]{2}-[0-9]{2})")

type GemtextHTML struct {
	*Site
	*Gemfile
	useAbsLinks bool

	content    strings.Builder
	blockquote bool
	pre        bool
	list       bool
	blanks     int
	sections   [3]bool
	anchors    map[string]struct{}
}

func (h *GemtextHTML) openSection(level int) {
	for i := 1; i <= level; i++ {
		if !h.sections[i-1] {
			fmt.Fprintf(&h.content, "<section class=level%d>\n", i)
			h.sections[i-1] = true
		}
	}
}

func (h *GemtextHTML) closeSection(level int) {
	if h.sections[level-1] {
		fmt.Fprintf(&h.content, "</section> <!-- level %d -->\n", level)
		h.sections[level-1] = false
	}
}

func (h *GemtextHTML) generateAnchor(text string) string {
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

func (h *GemtextHTML) spacingClass(classes ...string) string {
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

func markupRegexp(left, right string) *regexp.Regexp {
	return regexp.MustCompile(left + `([A-Za-z0-9 \\~./:_&#;()-]+?)` + right)
}

var smartquoteTT = markupRegexp("‘", "’")
var backtickTT = markupRegexp("`", "`")
var underlineItalic = markupRegexp("_", "_")
var doublestarBold = markupRegexp(`\*\*`, `\*\*`)
var starBold = markupRegexp(`\*`, `\*`)

func renderLine(line string) string {
	line = html.EscapeString(line)
	line = smartquoteTT.ReplaceAllString(line, `<tt>$1</tt>`)
	line = backtickTT.ReplaceAllString(line, `<tt>$1</tt>`)
	line = underlineItalic.ReplaceAllString(line, `<em>$1</em>`)
	line = doublestarBold.ReplaceAllString(line, `<strong>$1</strong>`)
	line = starBold.ReplaceAllString(line, `<strong>$1</strong>`)
	return line
}

func (h *GemtextHTML) Start() error {
	fmt.Fprint(&h.content, "<article class=gemtext>\n")
	return nil
}

func (h *GemtextHTML) Before(line gemini.Line) error {
	if _, ok := line.(gemini.LineListItem); ok {
		if !h.list {
			h.list = true
			fmt.Fprint(&h.content, "<ul>\n")
		}
	} else if h.list {
		h.list = false
		fmt.Fprint(&h.content, "</ul>\n")
	}
	return nil
}

func (h *GemtextHTML) Heading1(line gemini.LineHeading1) error {
	if h.Title == "" {
		h.Title = string(line)
	}
	h.closeSection(3)
	h.closeSection(2)
	h.closeSection(1)
	h.openSection(1)
	heading := string(line)
	anchor := h.generateAnchor(heading)
	fmt.Fprintf(&h.content, "<h1 id=%s><a class=header-link href=\"#%s\">¶</a><span>%s</span></h1>\n", anchor, anchor, html.EscapeString(heading))
	return nil
}

func (h *GemtextHTML) Heading2(line gemini.LineHeading2) error {
	h.closeSection(3)
	h.closeSection(2)
	h.openSection(2)
	heading := string(line)
	anchor := h.generateAnchor(heading)
	fmt.Fprintf(&h.content, "<h2 id=%s><a class=header-link href=\"#%s\">¶</a><span>%s</span></h2>\n", anchor, anchor, html.EscapeString(heading))
	return nil
}

func (h *GemtextHTML) Heading3(line gemini.LineHeading3) error {
	h.closeSection(3)
	h.openSection(3)
	heading := string(line)
	anchor := h.generateAnchor(heading)
	fmt.Fprintf(&h.content, "<h3 id=%s><a class=header-link href=\"#%s\">¶</a><span>%s</span></h3>\n", anchor, anchor, html.EscapeString(heading))
	return nil
}

func isImage(p string) bool {
	return strings.HasSuffix(p, ".png") || strings.HasSuffix(p, ".gif") || strings.HasSuffix(p, ".jpg")
}

func (h *GemtextHTML) Link(line gemini.LineLink) error {
	parsed, err := url.Parse(line.URL)
	if err != nil {
		return nil
	}
	linkIsRemote := parsed.Host != ""

	href := line.URL
	if strings.HasSuffix(href, ".gmi") && !parsed.IsAbs() {
		href = urlForGempath(href)
	}
	if !linkIsRemote && !parsed.IsAbs() && !path.IsAbs(href) && !h.SourceIsRootGemfile() {
		href = "../" + href
	}
	href = html.EscapeString(href)
	name := line.Name
	if name == "" {
		name = line.URL
	}
	name = html.EscapeString(name)

	if !linkIsRemote && h.useAbsLinks {
		href = fmt.Sprintf("https://%s%s", h.Domain, href)
	}

	if isImage(href) {
		fmt.Fprintf(&h.content, "<p%s><img src=\"%s\" alt=\"%s\"></p>\n", h.spacingClass(), href, name)
	} else {
		linkClass := "local"
		if linkIsRemote {
			linkClass = "remote"
		}
		fmt.Fprintf(&h.content, "<p%s><a class=%s href=\"%s\">%s</a></p>\n", h.spacingClass(), linkClass, href, name)
	}

	return nil
}

func (h *GemtextHTML) ListItem(line gemini.LineListItem) error {
	fmt.Fprintf(&h.content, "<li%s>%s</li>\n", h.spacingClass(), renderLine(string(line)))
	return nil
}

func (h *GemtextHTML) PreformattedText(line gemini.LinePreformattedText) error {
	fmt.Fprintf(&h.content, "%s\n", html.EscapeString(string(line)))
	return nil
}

func (h *GemtextHTML) PreformattingToggle(line gemini.LinePreformattingToggle) error {
	h.pre = !h.pre
	if h.pre {
		alt := strings.TrimSpace(string(line))
		if alt != "" {
			alt = html.EscapeString(alt)
			langClass := "lang-" + alt
			fmt.Fprintf(&h.content, "<pre%s aria-label='%s'>\n", h.spacingClass(langClass), alt)
		} else {
			fmt.Fprintf(&h.content, "<pre%s>\n", h.spacingClass())
		}
	} else {
		fmt.Fprint(&h.content, "</pre>\n")
	}
	return nil
}

func (h *GemtextHTML) Quote(line gemini.LineQuote) error {
	if line == "" {
		return nil
	}

	if dateRegexp.MatchString(string(line)) {
		h.Published = string(line[0:10])
		if updated := updatedRegexp.FindStringSubmatch(string(line)); updated != nil {
			h.Updated = updated[1]
		}
		fmt.Fprintf(&h.content, "<p%s>%s</p>\n", h.spacingClass("date"), html.EscapeString(string(line)))
		return nil
	}

	class := h.spacingClass()
	if h.blockquote {
		class = " class=nogap"
	}
	fmt.Fprintf(&h.content, "<blockquote%s>%s</blockquote>\n", class, html.EscapeString(string(line)))
	return nil
}

func (h *GemtextHTML) Text(line gemini.LineText) error {
	if line == "" {
		return nil
	}
	fmt.Fprintf(&h.content, "<p%s>%s</p>\n", h.spacingClass(), renderLine(string(line)))
	return nil
}

func (h *GemtextHTML) After(line gemini.Line) error {
	_, h.blockquote = line.(gemini.LineQuote)
	if line.String() == "" {
		h.blanks++
	} else {
		h.blanks = 0
	}
	return nil
}

func (h *GemtextHTML) Finish() error {
	if h.pre {
		fmt.Fprint(&h.content, "</pre>\n")
	}
	if h.list {
		fmt.Fprint(&h.content, "</ul>\n")
	}
	h.closeSection(3)
	h.closeSection(2)
	h.closeSection(1)
	fmt.Fprint(&h.content, "</article> <!-- gemtext -->\n")
	gemURL := strings.TrimSuffix(fmt.Sprintf("gemini://%s/%s", h.Domain, h.sourcePath), "index.gmi")
	fmt.Fprintf(&h.content, "<p class=gemlink>This page is also available via <a href=\"https://gemini.circumlunar.space/\">Gemini</a> at <a rel=alternate type=text/gemini href=\"%s\">%s</a></p>", gemURL, gemURL)
	return nil
}

type EmptyGemfileVisitor struct{}

func (_ EmptyGemfileVisitor) Start() error {
	return nil
}

func (_ EmptyGemfileVisitor) Before(line gemini.Line) error {
	return nil
}

func (_ EmptyGemfileVisitor) After(line gemini.Line) error {
	return nil
}

func (_ EmptyGemfileVisitor) Finish() error {
	return nil
}

func (_ EmptyGemfileVisitor) Heading1(line gemini.LineHeading1) error {
	return nil
}

func (_ EmptyGemfileVisitor) Heading2(line gemini.LineHeading2) error {
	return nil
}

func (_ EmptyGemfileVisitor) Heading3(line gemini.LineHeading3) error {
	return nil
}

func (_ EmptyGemfileVisitor) Link(line gemini.LineLink) error {
	return nil
}

func (_ EmptyGemfileVisitor) ListItem(line gemini.LineListItem) error {
	return nil
}

func (_ EmptyGemfileVisitor) PreformattedText(line gemini.LinePreformattedText) error {
	return nil
}

func (_ EmptyGemfileVisitor) PreformattingToggle(line gemini.LinePreformattingToggle) error {
	return nil
}

func (_ EmptyGemfileVisitor) Quote(line gemini.LineQuote) error {
	return nil
}

func (_ EmptyGemfileVisitor) Text(line gemini.LineText) error {
	return nil
}
