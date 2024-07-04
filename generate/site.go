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
)

// Site contains information about the site being generated.
type Site struct {
	// Domain is the domain name of the site
	Domain string

	// BaseDir is the local filesystem directory where the site is located.
	BaseDir string

	// OutputFiles contains all of the output files that will be created for
	// this site.
	OutputFiles map[string]OutputFile

	templates map[string]*template.Template
}

func NewSite(domain string, baseDir string) Site {
	return Site{
		Domain:      domain,
		BaseDir:     baseDir,
		OutputFiles: map[string]OutputFile{},
		templates:   map[string]*template.Template{},
	}
}

func (s *Site) OutputDir() string {
	return path.Join(s.BaseDir, ".html-new")
}

func (s *Site) AddOutputFile(file OutputFile) error {
	path := file.OutputPath()
	if existing, ok := s.OutputFiles[path]; ok {
		return fmt.Errorf("output file %s already exists: %s and %s", path, existing, file)
	}
	s.OutputFiles[path] = file
	return nil
}

func (s *Site) OpenOutputFile(outputPath string, mode fs.FileMode) (*os.File, error) {
	outputPath = path.Join(s.OutputDir(), outputPath)
	err := os.MkdirAll(path.Dir(outputPath), 0750)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
}

// URLForPath returns the site-relative URL for an output file (identified by
// its relative path in the local filesystem).
func (s *Site) URLForPath(outputPath string) string {
	if path.Base(outputPath) == "index.html" {
		dir := path.Dir(outputPath)
		if dir == "." {
			return "/"
		} else {
			return dir + "/"
		}
	}
	return outputPath
}

func (s *Site) AddSourceDir(sourceDir string) error {
	sourceLen := len(sourceDir) + 1
	return filepath.WalkDir(sourceDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		sourcePath := p[sourceLen:]
		ext := path.Ext(sourcePath)
		if ext == ".gmi" {
			file, err := NewGemfile(s, sourcePath, sourceDir, info.Mode())
			if err != nil {
				return err
			}
			err = s.AddOutputFile(file)
			if err != nil {
				return err
			}
		} else {
			file := NewCopyFile(sourcePath, sourceDir, info.Mode())
			err := s.AddOutputFile(file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Site) Generate() error {
	for _, file := range s.OutputFiles {
		err := file.Generate(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadMetadataFile loads a particular metadata file for an output path.  If
// there is no metadata file for the specific output path, we try looking for
// metadata files for its parent directories.
func (s *Site) LoadMetadataFile(outputPath string, meta string) (string, error) {
	for {
		metadataPath := path.Join(s.BaseDir, "metadata", outputPath, meta)
		metadata, err := os.ReadFile(metadataPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				next := path.Dir(outputPath)
				if next == outputPath {
					return "", nil
				}
				outputPath = next
				continue
			}
			return "", err
		}
		return string(metadata), nil
	}
}

// OutputFile represents one file that will be generated in the output
// directory.
type OutputFile interface {
	// OutputPath returns the local filesystem path (relative to the output
	// directory) of this output file.
	OutputPath() string

	// Generate creates the content of this output file and writes it to the
	// output directory.
	Generate(site *Site) error
}

// CopyFile is an OutputFile that directly copies content from a source
// directory into the output directory.
type CopyFile struct {
	path      string
	sourceDir string
	mode      fs.FileMode
}

func NewCopyFile(path string, sourceDir string, mode fs.FileMode) CopyFile {
	return CopyFile{path, sourceDir, mode}
}

func (f CopyFile) OutputPath() string {
	return f.path
}

func (f CopyFile) Generate(site *Site) error {
	out, err := site.OpenOutputFile(f.path, f.mode)
	if err != nil {
		return err
	}
	defer out.Close()

	sourcePath := path.Join(f.sourceDir, f.path)
	in, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	return err
}
