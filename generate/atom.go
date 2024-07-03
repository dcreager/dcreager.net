package generate

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"golang.org/x/tools/blog/atom"
)

func GenerateAtomFeed(domain, sourceDir, outputDir, author string) error {
	siteURL := fmt.Sprintf("https://%s/", domain)
	feedURL := fmt.Sprintf("https://%s/feed.atom", domain)
	siteLink := atom.Link{
		Rel:  "alternate",
		Type: "text/html",
		Href: siteURL,
	}
	feedLink := atom.Link{
		Rel:  "self",
		Type: "application/atom+xml",
		Href: feedURL,
	}

	archivePath := path.Join(sourceDir, "index.gmi")
	archive, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	var entries atomEntries
	err = entries.parseArchive(domain, sourceDir, archive)
	if err != nil {
		return err
	}

	feed := &atom.Feed{
		ID:      siteURL,
		Title:   domain,
		Link:    []atom.Link{siteLink, feedLink},
		Updated: atom.Time(entries.latestUpdate),
		Author: &atom.Person{
			Name: author,
			URI:  siteURL,
		},
		Entry: entries.entries,
	}

	encoded, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	output := path.Join(outputDir, "feed.atom")
	err = os.WriteFile(output, encoded, 0644)
	if err != nil {
		return err
	}

	return nil
}

type atomEntries struct {
	entries      []*atom.Entry
	latestUpdate time.Time
}

func (e *atomEntries) parseArchive(domain, sourceDir string, in io.Reader) error {
	return gemini.ParseLines(in, func(line gemini.Line) {
		switch line := line.(type) {
		case gemini.LineLink:
			if !dateRegexp.MatchString(line.Name) {
				return
			}

			publishedDate := line.Name[0:10]
			linkTitle := line.Name[11:]

			_, linkURL := gemlink(true, line.URL)
			id := fmt.Sprintf("tag:%s,%s:%s", domain, publishedDate, linkURL)

			targetPath := path.Join(sourceDir, line.URL)
			target, err := os.Open(targetPath)
			if err != nil {
				panic(err)
			}
			defer target.Close()

			var content strings.Builder
			hw := HTMLWriter{
				domain:      domain,
				Path:        "index.gmi",
				out:         &content,
				useAbsLinks: true,
			}
			gemini.ParseLines(target, hw.Handle)
			hw.Finish()

			if hw.Published != "" {
				publishedDate = hw.Published
			}

			updatedDate := publishedDate
			if hw.Updated != "" {
				updatedDate = hw.Updated
			}

			published, err := time.Parse("2006-01-02", publishedDate)
			if err != nil {
				panic(err)
			}
			updated, err := time.Parse("2006-01-02", updatedDate)
			if err != nil {
				panic(err)
			}

			if updated.After(e.latestUpdate) {
				e.latestUpdate = updated
			}

			contentText := &atom.Text{
				Type: "html",
				Body: content.String(),
			}

			entry := &atom.Entry{
				ID:        id,
				Title:     linkTitle,
				Published: atom.Time(published),
				Updated:   atom.Time(updated),
				Content:   contentText,
			}
			e.entries = append(e.entries, entry)
		}
	})
}
