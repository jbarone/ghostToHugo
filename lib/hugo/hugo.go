package hugo

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jbarone/ghostToHugo/lib/ghost"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/parser"
	"github.com/spf13/viper"
)

// Config represents the configuration for the Hugo project exporting to
type Config struct {
	Path string
}

func loadDefaultSettings(config Config) {
	viper.SetDefault("MetaDataFormat", "toml")
	viper.SetDefault("ContentDir", filepath.Join(config.Path, "content"))
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("StaticDir", "static")
	viper.SetDefault("ArchetypeDir", "archetypes")
	viper.SetDefault("PublishDir", "public")
	viper.SetDefault("DataDir", "data")
	viper.SetDefault("ThemesDir", "themes")
	viper.SetDefault("DefaultLayout", "post")
	viper.SetDefault("Verbose", false)
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
}

// Init ...
func Init(config Config) error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("hugo")
	viper.AddConfigPath(config.Path)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return err
		}
		return fmt.Errorf("Unable to locate Config file. Perhaps you need to create a new site. (%s)\n", err)
	}

	loadDefaultSettings(config)

	return nil
}

// ExportGhost will migrate entries from a Ghost export to a Hugo project
func ExportGhost(export *ghost.ExportData) error {
	var wg sync.WaitGroup
	for _, post := range export.Posts {
		wg.Add(1)
		go func(p ghost.Post) {
			defer wg.Done()
			writePost(p, export)
		}(post)
	}

	wg.Wait()
	return nil
}

func stripContentFolder(originalString string) string {
	return strings.Replace(originalString, "/content/", "/", -1)
}

func writePost(post ghost.Post, export *ghost.ExportData) error {
	var name = getPath(post)
	page, err := hugolib.NewPage(name)
	if err != nil {
		return err
	}

	if err := page.SetSourceMetaData(
		getMetadata(post, export),
		parser.FormatToLeadRune(viper.GetString("MetaDataFormat"))); err != nil {
		return err
	}

	page.SetSourceContent([]byte(post.Content))

	fmt.Println("saving file", name)
	return page.SafeSaveSourceAs(name)
}

func getMetadata(post ghost.Post, export *ghost.ExportData) map[string]interface{} {
	metadata := make(map[string]interface{})

	switch post.IsDraft() {
	case true:
		metadata["date"] = post.Created()
	case false:
		metadata["date"] = post.Published()
	}
	metadata["title"] = post.Title
	metadata["draft"] = post.IsDraft()
	metadata["slug"] = post.Slug
	metadata["description"] = post.MetaDescription
	if post.Image != "" {
		metadata["image"] = stripContentFolder(post.Image)
	}
	tags := post.Tags(export)
	if len(tags) > 0 {
		metadata["tags"] = tagsToSlice(tags)
		metadata["categories"] = tagsToSlice(tags)
	}
	author := post.Author(export)
	if author != nil {
		metadata["author"] = author.Name
	}

	return metadata
}

func tagsToSlice(tags []ghost.Tag) []interface{} {
	if len(tags) == 0 {
		return nil
	}

	t := make([]interface{}, len(tags))
	for i, tag := range tags {
		t[i] = tag.Name
	}

	return t
}

func getPath(post ghost.Post) string {
	if post.IsPage() {
		return helpers.AbsPathify(
			path.Join(viper.GetString("contentDir"), post.Slug+".md"))
	}

	return helpers.AbsPathify(
		path.Join(viper.GetString("contentDir"), "post", post.Slug+".md"))
}
