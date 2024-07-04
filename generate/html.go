package generate

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"path"
)

// OutputHTML renders a templated HTML output file.
func (s *Site) OutputHTML(outputPath string, mode fs.FileMode, title string, content string) error {
	out, err := s.OpenOutputFile(outputPath, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	customHead, err := s.LoadMetadataFile(outputPath, "custom-head")
	if err != nil {
		return err
	}

	data := htmlTemplateData{
		Title:      title,
		Content:    template.HTML(content),
		CustomHead: template.HTML(customHead),
	}
	t, err := s.FindTemplate(outputPath)
	if err != nil {
		return err
	}
	err = t.Execute(out, data)
	if err != nil {
		return err
	}

	return nil
}

type htmlTemplateData struct {
	Title      string
	Content    template.HTML
	CustomHead template.HTML
}

func (s *Site) LoadTemplate(outputPath string) (*template.Template, error) {
	if t, ok := s.templates[outputPath]; ok {
		return t, nil
	}

	templatePath := path.Join(s.BaseDir, "templates", outputPath)
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	templates[outputPath] = t
	return t, nil
}

func (s *Site) FindTemplate(outputPath string) (*template.Template, error) {
	t, err := s.LoadTemplate(outputPath)
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

	return nil, fmt.Errorf("Cannot find suitable template for %s", outputPath)
}
