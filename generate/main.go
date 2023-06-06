package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	count := 0
	start := time.Now()
	sourceDir := path.Join(home, "notes")
	outputDir := ".html"
	sourceLen := len(sourceDir) + 1
	err = filepath.WalkDir(sourceDir, func(p string, d fs.DirEntry, err error) error {
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
		return err
	}

	duration := time.Since(start)
	fmt.Printf("Processed %d files in %.2f seconds.\n", count, duration.Seconds())
	return nil
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
