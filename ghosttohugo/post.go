package ghosttohugo

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/parser"
	"github.com/jbarone/mobiledoc"
	jww "github.com/spf13/jwalterweatherman"
)

type post struct {
	ID              json.RawMessage `json:"id"`
	Title           string          `json:"title"`
	Slug            string          `json:"slug"`
	Content         string          `json:"markdown"`
	Plain           string          `json:"plaintext"`
	MobileDoc       string          `json:"mobiledoc,omitempty"`
	Image           string          `json:"image"`
	FeaturedImage   string          `json:"feature_image,omitempty"`
	Type            string          `json:"type"`
	Status          string          `json:"status"`
	MetaDescription string          `json:"meta_description"`
	AuthorID        json.RawMessage `json:"author_id"`
	PublishedAt     json.RawMessage `json:"published_at"`
	CreatedAt       json.RawMessage `json:"created_at"`
	Summary         string          `json:"custom_excerpt"`

	Published time.Time
	Created   time.Time
	Author    string
	Tags      []string
}

func (p post) isDraft() bool {
	return strings.ToLower(p.Status) == "draft"
}

func (p post) isPage() bool {
	return strings.ToLower(p.Type) == "page"
}

func (p post) frontMatter() map[string]interface{} {
	metadata := make(map[string]interface{})

	switch p.isDraft() {
	case true:
		metadata["date"] = p.Created
	case false:
		metadata["date"] = p.Published
	}
	metadata["title"] = p.Title
	metadata["draft"] = p.isDraft()
	metadata["slug"] = p.Slug
	metadata["description"] = p.MetaDescription
	if p.Image != "" {
		metadata["image"] = stripContentFolder(p.Image)
	} else if p.FeaturedImage != "" {
		metadata["image"] = stripContentFolder(p.FeaturedImage)
	}
	if len(p.Tags) > 0 {
		metadata["tags"] = p.Tags
		metadata["categories"] = p.Tags
	}
	if p.Author != "" {
		metadata["author"] = p.Author
	}
	if p.Summary != "" {
		metadata["summary"] = p.Summary
	}

	return metadata
}

func (c *Converter) writePost(p post) error {
	jww.DEBUG.Printf("converting: %s", p.Title)
	path := filepath.Join(c.path, "content")
	switch p.isPage() {
	case true:
		path = filepath.Join(path, p.Slug+".md")
	case false:
		path = filepath.Join(path, "post", p.Slug+".md")
	}

	buf := bytes.NewBuffer(nil)
	if err := parser.InterfaceToFrontMatter(p.frontMatter(), c.kind, buf); err != nil {
		return err
	}
	if _, err := buf.Write([]byte("\n\n")); err != nil {
		return err
	}

	switch {
	case p.Content != "":
		if _, err := buf.Write([]byte(p.Content)); err != nil {
			return err
		}
	case p.MobileDoc != "":
		if _, err := buf.Write([]byte(p.mobiledocMarkdown())); err != nil {
			return err
		}
	default:
		if _, err := buf.Write([]byte(p.Plain)); err != nil {
			return err
		}
	}

	return helpers.WriteToDisk(path, bytes.NewReader(buf.Bytes()), c.site.Fs.Source)
}

func (p post) mobiledocMarkdown() string {
	if p.MobileDoc == "" {
		return ""
	}

	r := strings.NewReader(p.MobileDoc)
	var buf bytes.Buffer

	md := mobiledoc.NewMobiledoc(r).
		WithAtom("soft-break", atomSoftReturn).
		WithAtom("soft-return", atomSoftReturn).
		WithCard("card-markdown", cardMarkdown).
		WithCard("markdown", cardMarkdown).
		WithCard("hr", cardHR).
		WithCard("image", cardImage).
		WithCard("code", cardCode).
		WithCard("embed", cardEmbed).
		WithCard("gallery", cardGallery).
		WithCard("html", cardHTML).
		WithCard("bookmark", cardBookmark)

	err := md.Render(&buf)
	if err != nil {
		jww.ERROR.Printf("error rendering post %s (%v)\n", p.ID, err)
		return ""
	}

	return buf.String()
}
