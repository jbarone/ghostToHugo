package ghosttohugo

import (
	"bytes"
	"fmt"
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
	mkdir(c.path, "content")
	mkdir(c.path, "archetypes")
	mkdir(c.path, "static")
	mkdir(c.path, "data")
	mkdir(c.path, "themes")

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
