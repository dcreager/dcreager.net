package main

import (
	"io/fs"
	"os"
	"path"

	"git.sr.ht/~adnano/go-gemini"
)

func translateGmiPath(p string) string {
	if path.Base(p) == "index.gmi" {
		return p[:len(p)-len("index.gmi")] + "/"
	} else {
		return p[:len(p)-len(".gmi")] + "/"
	}
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

	hw := HTMLWriter{out: out}
	gemini.ParseLines(in, hw.Handle)
	hw.Finish()
	return nil
}
