package generate

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"git.sr.ht/~adnano/go-gemini"
)

func ProcessSourceDir(sourceDir, outputDir string) (int, error) {
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
			return translateGemtext(base, p, output, d)
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
	t, err := loadTemplate("default.html")
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	t, err = loadTemplate(base)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	return nil, fmt.Errorf("Cannot find suitable template for %s", base)
}

func translateGmiPath(p string) string {
	if path.Base(p) == "index.gmi" {
		return p[:len(p)-len("index.gmi")] + "/"
	} else {
		return p[:len(p)-len(".gmi")] + "/"
	}
}

type geminiTemplateData struct {
	Title           string
	RenderedGemtext template.HTML
}

func translateGemtext(base, from, to string, d fs.DirEntry) error {
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
		Path:   base,
		out:    &content,
		isRoot: path.Base(base) == "index.gmi",
	}
	gemini.ParseLines(in, hw.Handle)
	hw.Finish()

	data := geminiTemplateData{
		Title:           hw.Title,
		RenderedGemtext: template.HTML(content.String()),
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
