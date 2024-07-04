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

type AtomFeed struct{}

func (f AtomFeed) OutputPath() string {
	return "feed.atom"
}

func (f AtomFeed) Generate(site *Site) error {
	siteURL := fmt.Sprintf("https://%s/", site.Domain)
	feedURL := fmt.Sprintf("https://%s/feed.atom", site.Domain)
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

	archive := site.OutputFiles["index.html"]
	var entries atomEntries
	err := entries.parseArchive(site, archive.Source)
	if err != nil {
		return err
	}

	feed := &atom.Feed{
		ID:      siteURL,
		Title:   site.Domain,
		Link:    []atom.Link{siteLink, feedLink},
		Updated: atom.Time(entries.latestUpdate),
		Author: &atom.Person{
			Name: site.Author,
			URI:  siteURL,
		},
		Entry: entries.entries,
	}

	encoded, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	out, err := site.OpenOutputFile("feed.atom", 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	err = out.Write(encoded)
	if err != nil {
		return err
	}

	return nil
}

type atomEntries struct {
	entries      []*atom.Entry
	latestUpdate time.Time
}

func (e *atomEntries) parseArchive(site *Site, lines gemini.Text) error {
	for _, line := range lines {
		switch line := line.(type) {
		case gemini.LineLink:
			if !dateRegexp.MatchString(line.Name) {
				continue
			}

			publishedDate := line.Name[0:10]
			linkTitle := line.Name[11:]

			_, linkURL := gemlink(true, line.URL)
			id := fmt.Sprintf("tag:%s,%s:%s", domain, publishedDate, linkURL)

			targetPath := path.Join(sourceDir, line.URL)
			target, err := os.Open(targetPath)
			if err != nil {
				return err
			}
			defer target.Close()

			if hw.Published != "" {
				publishedDate = hw.Published
			}

			updatedDate := publishedDate
			if hw.Updated != "" {
				updatedDate = hw.Updated
			}

			published, err := time.Parse("2006-01-02", publishedDate)
			if err != nil {
				return err
			}
			updated, err := time.Parse("2006-01-02", updatedDate)
			if err != nil {
				return err
			}

			if updated.After(e.latestUpdate) {
				e.latestUpdate = updated
			}

			entryLink := atom.Link{
				Rel:  "alternate",
				Type: "text/html",
				Href: fmt.Sprintf("https://%s%s", domain, linkURL),
			}

			summary := &atom.Text{
				Type: "html",
				Body: content.String(),
			}

			entry := &atom.Entry{
				ID:        id,
				Title:     linkTitle,
				Published: atom.Time(published),
				Updated:   atom.Time(updated),
				Summary:   summary,
				Link:      []atom.Link{entryLink},
			}
			e.entries = append(e.entries, entry)
		}
	}
}
