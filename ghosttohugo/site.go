package ghosttohugo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/parser"
)

func (c *Converter) createSite() error {
	s, err := hugolib.NewSiteDefaultLang()
	if err != nil {
		return err
	}

	fs := s.Fs.Source
	if exists, _ := helpers.Exists(c.path, fs); exists {
		if isDir, _ := helpers.IsDir(c.path, fs); !isDir {
			return fmt.Errorf(
				"target path %q exists but is not a directory",
				c.path,
			)
		}

		isEmpty, _ := helpers.IsEmpty(c.path, fs)

		if !isEmpty && !c.force {
			return fmt.Errorf(
				"target path %q exists and is not empty",
				c.path,
			)
		}
	}

	mkdir(c.path, "layouts")
	mkdir(c.path, filepath.Clean("layouts/shortcodes"))
	mkdir(c.path, "content")
	mkdir(c.path, "archetypes")
	mkdir(c.path, "static")
	mkdir(c.path, "data")
	mkdir(c.path, "themes")

	ioutil.WriteFile(
		filepath.Join(c.path, "layouts/shortcodes/bookmark.html"),
		bookmarkData,
		0644,
	)
	ioutil.WriteFile(
		filepath.Join(c.path, "layouts/shortcodes/gallery.html"),
		galleryData,
		0644,
	)
	ioutil.WriteFile(
		filepath.Join(c.path, "layouts/shortcodes/galleryImg.html"),
		galleryImgData,
		0644,
	)

	c.site = s

	c.createConfig()

	return nil
}

func (c Converter) createConfig() error {
	title := "My New Hugo Site"
	baseURL := "http://example.org/"

	for key, value := range c.info.settings {
		switch strings.ToLower(key) {
		case "title":
			title = value
		}
	}

	in := map[string]interface{}{
		"baseURL":            baseURL,
		"title":              title,
		"languageCode":       "en-us",
		"disablePathToLower": true,
		"markup": map[string]interface{}{
			"goldmark": map[string]interface{}{
				"renderer": map[string]interface{}{"unsafe": true},
			},
		},
	}

	var buf bytes.Buffer
	if err := parser.InterfaceToConfig(in, c.kind, &buf); err != nil {
		return err
	}

	return helpers.WriteToDisk(
		filepath.Join(c.path, "config."+string(c.kind)),
		&buf,
		c.site.Fs.Source,
	)
}

var bookmarkData = []byte(`<figure class="kg-card kg-bookmark-card">
  <a href="{{ .Get "url" }}" class="kg-bookmark-container">
    <div class="kg-bookmark-content">
      <div class="kg-bookmark-title">{{ .Get "title" }}</div>
      <div class="kg-bookmark-description">{{ .Get "description" }}</div>
      <div class="kg-bookmark-metadata">
        {{ with .Get "icon" }}<img src="{{ . }}" class="kg-bookmark-icon">{{ end }}
        {{ with .Get "author" }}<span class="kg-bookmark-author">{{ . }}</span>{{ end }}
        {{ with .Get "publisher" }}<span class="kg-bookmark-publisher">{{ . }}</span>{{ end }}
      </div>
    </div>
    {{ with .Get "thumbnail" }}
    <div class="kg-bookmark-thumbnail">
      <img src="{{ . }}">
    </div>
    {{ end }}
  </a>
  {{ if .Get "caption" }}
  <figcaption>{{ . }}</figcaption>
  {{ end }}
</figure>`)

var galleryData = []byte(`<figure class="kg-gallery-card kg-width-wide">
  <div class="kg-gallery-container">
    <div class="kg-gallery-row">
    {{ .Inner }}
    </div>
  </div>
  {{ with .Get "caption" }}
  <figcaption>{{ . }}</figcaption>
  {{ end }}
</figure>`)

var galleryImgData = []byte(`
  <div class="kg-gallery-image">
    <img src="{{ .Get "src" }}" width="{{ .Get "width" }}" height="{{ .Get "height" }}">
  </div>
{{ if mod .Ordinal 3 | eq 2 }}
</div>
<div class="kg-gallery-row">
{{ end }}`)
